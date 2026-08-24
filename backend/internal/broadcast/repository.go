package broadcast

import (
	"context"
	"encoding/json"
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

func (r *PostgresRepository) List(ctx context.Context, from, to time.Time, includeHidden bool) ([]Listing, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT l.id,l.match_id,coalesce(c.starts_at,l.starts_at),coalesce(c.home_name,l.home_name),
		       coalesce(c.away_name,l.away_name),coalesce(c.label,l.label),
		       coalesce(c.competition_name,l.competition_name),coalesce(c.broadcast_kind,l.broadcast_kind),
		       CASE WHEN c.listing_id IS NULL THEN coalesce(array_agg(ch.channel ORDER BY ch.position) FILTER (WHERE ch.channel IS NOT NULL),'{}') ELSE c.channels END,
		       l.source_url,l.match_id IS NULL,c.listing_id IS NOT NULL,coalesce(c.hidden,false)
		FROM tv_listings l
		LEFT JOIN tv_listing_channels ch ON ch.listing_id=l.id
		LEFT JOIN tv_listing_corrections c ON c.listing_id=l.id
		WHERE (coalesce(c.starts_at,l.starts_at) AT TIME ZONE 'Europe/Paris')::date BETWEEN $1::date AND $2::date
		  AND ($3 OR NOT coalesce(c.hidden,false))
		GROUP BY l.id,c.listing_id,c.starts_at,c.home_name,c.away_name,c.label,c.competition_name,c.broadcast_kind,c.channels,c.hidden
		ORDER BY coalesce(c.starts_at,l.starts_at),coalesce(c.label,l.label)`, from, to, includeHidden)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Listing, 0)
	for rows.Next() {
		var item Listing
		if err := rows.Scan(&item.ID, &item.MatchID, &item.StartsAt, &item.HomeName, &item.AwayName, &item.Label, &item.CompetitionName, &item.Kind, &item.Channels, &item.SourceURL, &item.External, &item.Corrected, &item.Hidden); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) FindMatchCandidates(ctx context.Context, at time.Time) ([]MatchCandidate, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,home_name,away_name FROM matches WHERE utc_date BETWEEN $1::timestamptz-interval '20 minutes' AND $1::timestamptz+interval '20 minutes'`, at)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]MatchCandidate, 0)
	for rows.Next() {
		var item MatchCandidate
		if err := rows.Scan(&item.ID, &item.HomeName, &item.AwayName); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) Upsert(ctx context.Context, item ImportedListing, matchID *uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var id uuid.UUID
	err = tx.QueryRow(ctx, `INSERT INTO tv_listings(id,source,source_key,match_id,starts_at,home_name,away_name,label,competition_name,broadcast_kind,source_url)
		VALUES($1,'footao',$2,$3,$4,$5,$6,$7,$8,$9,$10)
		ON CONFLICT(source,source_key) DO UPDATE SET match_id=excluded.match_id,starts_at=excluded.starts_at,home_name=excluded.home_name,away_name=excluded.away_name,label=excluded.label,competition_name=excluded.competition_name,broadcast_kind=excluded.broadcast_kind,source_url=excluded.source_url,imported_at=now(),updated_at=now()
		RETURNING id`, uuid.New(), item.SourceKey, matchID, item.StartsAt, item.HomeName, item.AwayName, item.Label, item.CompetitionName, item.Kind, item.SourceURL).Scan(&id)
	if err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM tv_listing_channels WHERE listing_id=$1`, id); err != nil {
		return err
	}
	for i, channel := range item.Channels {
		if _, err = tx.Exec(ctx, `INSERT INTO tv_listing_channels(listing_id,channel,position) VALUES($1,$2,$3)`, id, channel, i); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *PostgresRepository) CleanupPast(ctx context.Context, now time.Time) (int64, error) {
	result, err := r.pool.Exec(ctx, `DELETE FROM tv_listings WHERE (starts_at AT TIME ZONE 'Europe/Paris')::date < ($1::timestamptz AT TIME ZONE 'Europe/Paris')::date`, now)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (r *PostgresRepository) Correct(ctx context.Context, listingID, adminID uuid.UUID, input CorrectionInput) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	before, err := listingSnapshot(ctx, tx, listingID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrListingMissing
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO tv_listing_corrections(listing_id,starts_at,home_name,away_name,label,competition_name,broadcast_kind,channels,hidden,note,updated_by)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT(listing_id) DO UPDATE SET starts_at=excluded.starts_at,home_name=excluded.home_name,away_name=excluded.away_name,label=excluded.label,competition_name=excluded.competition_name,broadcast_kind=excluded.broadcast_kind,channels=excluded.channels,hidden=excluded.hidden,note=excluded.note,updated_by=excluded.updated_by,updated_at=now()`, listingID, input.StartsAt, input.HomeName, input.AwayName, input.Label, input.CompetitionName, input.Kind, input.Channels, input.Hidden, input.Note, adminID)
	if err != nil {
		return err
	}
	after := map[string]any{"startsAt": input.StartsAt, "homeName": input.HomeName, "awayName": input.AwayName, "label": input.Label, "competitionName": input.CompetitionName, "kind": input.Kind, "channels": input.Channels, "hidden": input.Hidden, "note": input.Note}
	return writeAudit(ctx, tx, listingID, adminID, "corrected", before, after)
}

func (r *PostgresRepository) ClearCorrection(ctx context.Context, listingID, adminID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	before, err := correctionSnapshot(ctx, tx, listingID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrListingMissing
	}
	if err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM tv_listing_corrections WHERE listing_id=$1`, listingID); err != nil {
		return err
	}
	return writeAudit(ctx, tx, listingID, adminID, "cleared", before, nil)
}

func listingSnapshot(ctx context.Context, tx pgx.Tx, id uuid.UUID) (map[string]any, error) {
	var starts time.Time
	var home, away, label, competition, kind, note string
	var channels []string
	var hidden bool
	err := tx.QueryRow(ctx, `SELECT coalesce(c.starts_at,l.starts_at),coalesce(c.home_name,l.home_name),coalesce(c.away_name,l.away_name),coalesce(c.label,l.label),coalesce(c.competition_name,l.competition_name),coalesce(c.broadcast_kind,l.broadcast_kind),CASE WHEN c.listing_id IS NULL THEN ARRAY(SELECT channel FROM tv_listing_channels WHERE listing_id=l.id ORDER BY position) ELSE c.channels END,coalesce(c.hidden,false),coalesce(c.note,'') FROM tv_listings l LEFT JOIN tv_listing_corrections c ON c.listing_id=l.id WHERE l.id=$1`, id).Scan(&starts, &home, &away, &label, &competition, &kind, &channels, &hidden, &note)
	return map[string]any{"startsAt": starts, "homeName": home, "awayName": away, "label": label, "competitionName": competition, "kind": kind, "channels": channels, "hidden": hidden, "note": note}, err
}
func correctionSnapshot(ctx context.Context, tx pgx.Tx, id uuid.UUID) (map[string]any, error) {
	var starts time.Time
	var home, away, label, competition, kind, note string
	var channels []string
	var hidden bool
	err := tx.QueryRow(ctx, `SELECT starts_at,home_name,away_name,label,competition_name,broadcast_kind,channels,hidden,note FROM tv_listing_corrections WHERE listing_id=$1`, id).Scan(&starts, &home, &away, &label, &competition, &kind, &channels, &hidden, &note)
	return map[string]any{"startsAt": starts, "homeName": home, "awayName": away, "label": label, "competitionName": competition, "kind": kind, "channels": channels, "hidden": hidden, "note": note}, err
}
func writeAudit(ctx context.Context, tx pgx.Tx, listingID, adminID uuid.UUID, action string, before, after map[string]any) error {
	b, _ := json.Marshal(before)
	a, _ := json.Marshal(after)
	if _, err := tx.Exec(ctx, `INSERT INTO tv_listing_audit(id,listing_id,admin_id,action,before_state,after_state) VALUES($1,$2,$3,$4,$5,$6)`, uuid.New(), listingID, adminID, action, b, a); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *PostgresRepository) Audit(ctx context.Context, limit int) ([]AuditEntry, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,listing_id,admin_id,action,before_state,after_state,created_at FROM tv_listing_audit ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]AuditEntry, 0)
	for rows.Next() {
		var item AuditEntry
		var before, after []byte
		if err := rows.Scan(&item.ID, &item.ListingID, &item.AdminID, &item.Action, &before, &after, &item.CreatedAt); err != nil {
			return nil, err
		}
		if len(before) > 0 {
			_ = json.Unmarshal(before, &item.Before)
		}
		if len(after) > 0 {
			_ = json.Unmarshal(after, &item.After)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
func (r *PostgresRepository) StartCollection(ctx context.Context) (CollectionRun, error) {
	item := CollectionRun{ID: uuid.New(), Status: "running"}
	err := r.pool.QueryRow(ctx, `INSERT INTO collection_runs(id,connector,scope,status) VALUES($1,'footao','broadcasts','running') RETURNING started_at`, item.ID).Scan(&item.StartedAt)
	return item, err
}
func (r *PostgresRepository) FinishCollection(ctx context.Context, id uuid.UUID, status string, count int, message string) error {
	_, err := r.pool.Exec(ctx, `UPDATE collection_runs SET status=$2,records_count=$3,error_message=$4,finished_at=now() WHERE id=$1`, id, status, count, message)
	return err
}
func (r *PostgresRepository) LatestCollection(ctx context.Context) (*CollectionRun, error) {
	var item CollectionRun
	err := r.pool.QueryRow(ctx, `SELECT id,status,records_count,error_message,started_at,finished_at FROM collection_runs WHERE connector='footao' AND scope='broadcasts' ORDER BY started_at DESC LIMIT 1`).Scan(&item.ID, &item.Status, &item.RecordsCount, &item.ErrorMessage, &item.StartedAt, &item.FinishedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &item, err
}
