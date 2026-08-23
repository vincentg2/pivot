package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")
var ErrInvitationInvalid = errors.New("invitation invalid")
var ErrEmailTaken = errors.New("email taken")

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func scanUser(row pgx.Row) (User, error) {
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Nickname, &user.AvatarSeed, &user.Theme, &user.Role, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return user, err
}

func (r *PostgresRepository) FindUserByEmail(ctx context.Context, email string) (User, error) {
	return scanUser(r.pool.QueryRow(ctx, `SELECT id,email,password_hash,nickname,avatar_seed,theme,role,created_at FROM users WHERE email=$1`, email))
}

func (r *PostgresRepository) FindUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	return scanUser(r.pool.QueryRow(ctx, `SELECT id,email,password_hash,nickname,avatar_seed,theme,role,created_at FROM users WHERE id=$1`, id))
}

func (r *PostgresRepository) RegisterWithInvitation(ctx context.Context, input NewUser, codeHash []byte, now time.Time) (User, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return User{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var invitationID uuid.UUID
	err = tx.QueryRow(ctx, `SELECT id FROM invitations WHERE code_hash=$1 AND revoked_at IS NULL AND (expires_at IS NULL OR expires_at>$2) AND uses<max_uses FOR UPDATE`, codeHash, now).Scan(&invitationID)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrInvitationInvalid
	}
	if err != nil {
		return User{}, err
	}
	input.InvitationID = invitationID
	user, err := scanUser(tx.QueryRow(ctx, `INSERT INTO users(id,email,password_hash,nickname,avatar_seed) VALUES($1,$2,$3,$4,$5) RETURNING id,email,password_hash,nickname,avatar_seed,theme,role,created_at`, input.ID, input.Email, input.PasswordHash, input.Nickname, input.AvatarSeed))
	if err != nil {
		if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, ErrEmailTaken
		}
		return User{}, err
	}
	if _, err = tx.Exec(ctx, `UPDATE invitations SET uses=uses+1 WHERE id=$1`, invitationID); err != nil {
		return User{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO invitation_uses(invitation_id,user_id) VALUES($1,$2)`, invitationID, user.ID); err != nil {
		return User{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *PostgresRepository) CreateSession(ctx context.Context, userID uuid.UUID, hash []byte, expires time.Time, agent, ip string) error {
	var ipValue any
	if ip != "" {
		ipValue = ip
	}
	_, err := r.pool.Exec(ctx, `INSERT INTO sessions(id,token_hash,user_id,expires_at,user_agent,ip) VALUES($1,$2,$3,$4,$5,$6)`, uuid.New(), hash, userID, expires, agent, ipValue)
	return err
}

func (r *PostgresRepository) FindUserBySession(ctx context.Context, hash []byte) (User, error) {
	return scanUser(r.pool.QueryRow(ctx, `SELECT u.id,u.email,u.password_hash,u.nickname,u.avatar_seed,u.theme,u.role,u.created_at FROM sessions s JOIN users u ON u.id=s.user_id WHERE s.token_hash=$1 AND s.expires_at>now()`, hash))
}

func (r *PostgresRepository) DeleteSession(ctx context.Context, hash []byte) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE token_hash=$1`, hash)
	return err
}
