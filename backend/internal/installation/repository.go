package installation

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}
func (r *PostgresRepository) Required(ctx context.Context) (bool, error) {
	var required bool
	err := r.pool.QueryRow(ctx, `SELECT NOT EXISTS(SELECT 1 FROM users)`).Scan(&required)
	return required, err
}
func (r *PostgresRepository) CreateFirstAdmin(ctx context.Context, input NewAdmin) (Admin, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Admin{}, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `SELECT pg_advisory_xact_lock(719840617)`); err != nil {
		return Admin{}, err
	}
	var exists bool
	if err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users)`).Scan(&exists); err != nil {
		return Admin{}, err
	}
	if exists {
		return Admin{}, ErrAlreadyInstalled
	}
	var item Admin
	err = tx.QueryRow(ctx, `INSERT INTO users(id,email,password_hash,nickname,avatar_seed,role) VALUES($1,$2,$3,$4,$5,'admin') RETURNING id,email,nickname,role`, input.ID, input.Email, input.PasswordHash, input.Nickname, input.AvatarSeed).Scan(&item.ID, &item.Email, &item.Nickname, &item.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		return Admin{}, ErrAlreadyInstalled
	}
	if err != nil {
		return Admin{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return Admin{}, err
	}
	return item, nil
}
