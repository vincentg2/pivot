package invitation

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("invitation not found")

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, invite Invitation, hash []byte, createdBy uuid.UUID) (Invitation, error) {
	err := r.pool.QueryRow(ctx, `INSERT INTO invitations(id,code_hash,label,expires_at,max_uses,created_by) VALUES($1,$2,$3,$4,$5,$6) RETURNING created_at`, invite.ID, hash, invite.Label, invite.ExpiresAt, invite.MaxUses, createdBy).Scan(&invite.CreatedAt)
	return invite, err
}

func (r *PostgresRepository) List(ctx context.Context) ([]Invitation, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,label,expires_at,max_uses,uses,created_at,revoked_at FROM invitations ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Invitation, 0)
	for rows.Next() {
		var item Invitation
		if err := rows.Scan(&item.ID, &item.Label, &item.ExpiresAt, &item.MaxUses, &item.Uses, &item.CreatedAt, &item.RevokedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `UPDATE invitations SET revoked_at=now() WHERE id=$1 AND revoked_at IS NULL`, id)
	if err == nil && result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

var _ = pgx.ErrNoRows
