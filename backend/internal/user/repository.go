package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}
func (r *PostgresRepository) UpdateProfile(ctx context.Context, id uuid.UUID, nickname, theme, locale string) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET nickname=$2,theme=$3,locale=$4,updated_at=now() WHERE id=$1`, id, nickname, theme, locale)
	return err
}
func (r *PostgresRepository) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}
