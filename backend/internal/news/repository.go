package news

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) ListItems(ctx context.Context, userID uuid.UUID, clubID *uuid.UUID, limit int) ([]Item, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT n.id,n.club_id,c.name,c.tla,n.title,n.source_name,n.published_at,n.link_url
		FROM news_items n JOIN clubs c ON c.id=n.club_id
		WHERE ($2::uuid IS NOT NULL AND n.club_id=$2)
		   OR ($2::uuid IS NULL AND EXISTS(SELECT 1 FROM favorite_clubs f WHERE f.user_id=$1 AND f.club_id=n.club_id))
		ORDER BY n.published_at DESC,n.id LIMIT $3`, userID, clubID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Item, 0)
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.ClubID, &item.ClubName, &item.ClubTLA, &item.Title, &item.SourceName, &item.PublishedAt, &item.LinkURL); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) ListFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := r.pool.Query(ctx, `SELECT f.id,f.club_id,c.name,f.feed_url,f.source_name,f.enabled,f.updated_at FROM news_feeds f JOIN clubs c ON c.id=f.club_id ORDER BY c.name,f.source_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Feed, 0)
	for rows.Next() {
		var item Feed
		if err := rows.Scan(&item.ID, &item.ClubID, &item.ClubName, &item.URL, &item.SourceName, &item.Enabled, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
func (r *PostgresRepository) SaveFeed(ctx context.Context, adminID uuid.UUID, input FeedInput) (Feed, error) {
	item := Feed{ID: uuid.New(), ClubID: input.ClubID, URL: input.URL, SourceName: input.SourceName, Enabled: input.Enabled}
	err := r.pool.QueryRow(ctx, `INSERT INTO news_feeds(id,club_id,feed_url,source_name,enabled,created_by) SELECT $1,id,$3,$4,$5,$6 FROM clubs WHERE id=$2 ON CONFLICT(club_id,feed_url) DO UPDATE SET source_name=excluded.source_name,enabled=excluded.enabled,updated_at=now() RETURNING id,updated_at`, item.ID, input.ClubID, input.URL, input.SourceName, input.Enabled, adminID).Scan(&item.ID, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Feed{}, ErrClubMissing
	}
	if err != nil {
		return Feed{}, err
	}
	_ = r.pool.QueryRow(ctx, `SELECT name FROM clubs WHERE id=$1`, input.ClubID).Scan(&item.ClubName)
	return item, nil
}
func (r *PostgresRepository) DeleteFeed(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM news_feeds WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrFeedMissing
	}
	return nil
}
func (r *PostgresRepository) UpsertItems(ctx context.Context, feed Feed, items []ImportedItem) (int, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	count := 0
	for _, item := range items {
		result, execErr := tx.Exec(ctx, `INSERT INTO news_items(id,feed_id,club_id,external_id,title,source_name,published_at,link_url) VALUES($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT(feed_id,external_id) DO UPDATE SET title=excluded.title,source_name=excluded.source_name,published_at=excluded.published_at,link_url=excluded.link_url,imported_at=now()`, uuid.New(), feed.ID, feed.ClubID, item.ExternalID, item.Title, feed.SourceName, item.PublishedAt, item.LinkURL)
		if execErr != nil {
			return 0, execErr
		}
		count += int(result.RowsAffected())
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return count, nil
}
func (r *PostgresRepository) Cleanup(ctx context.Context, now time.Time) (int64, error) {
	result, err := r.pool.Exec(ctx, `DELETE FROM news_items WHERE published_at < $1::timestamptz-interval '30 days'`, now)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
func (r *PostgresRepository) StartCollection(ctx context.Context) (CollectionRun, error) {
	item := CollectionRun{ID: uuid.New(), Status: "running"}
	err := r.pool.QueryRow(ctx, `INSERT INTO collection_runs(id,connector,scope,status) VALUES($1,'official-rss','news','running') RETURNING started_at`, item.ID).Scan(&item.StartedAt)
	return item, err
}
func (r *PostgresRepository) FinishCollection(ctx context.Context, id uuid.UUID, status string, count int, message string) error {
	_, err := r.pool.Exec(ctx, `UPDATE collection_runs SET status=$2,records_count=$3,error_message=$4,finished_at=now() WHERE id=$1`, id, status, count, message)
	return err
}
func (r *PostgresRepository) LatestCollection(ctx context.Context) (*CollectionRun, error) {
	var item CollectionRun
	err := r.pool.QueryRow(ctx, `SELECT id,status,records_count,error_message,started_at,finished_at FROM collection_runs WHERE connector='official-rss' AND scope='news' ORDER BY started_at DESC LIMIT 1`).Scan(&item.ID, &item.Status, &item.RecordsCount, &item.ErrorMessage, &item.StartedAt, &item.FinishedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &item, err
}
