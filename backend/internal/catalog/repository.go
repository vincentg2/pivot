package catalog

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) ListCompetitions(ctx context.Context) ([]Competition, error) {
	rows, err := r.pool.Query(ctx, `SELECT c.id,c.provider_code,c.name,c.country,c.emblem_url,count(cc.club_id) FROM competitions c LEFT JOIN competition_clubs cc ON cc.competition_id=c.id GROUP BY c.id ORDER BY c.country,c.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Competition, 0)
	for rows.Next() {
		var item Competition
		if err := rows.Scan(&item.ID, &item.ProviderCode, &item.Name, &item.Country, &item.EmblemURL, &item.ClubCount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) ListClubs(ctx context.Context, competitionCode string) ([]Club, error) {
	rows, err := r.pool.Query(ctx, `SELECT DISTINCT cl.id,cl.provider_id,cl.name,cl.short_name,cl.tla,cl.crest_url,cl.website_url,cl.venue FROM clubs cl LEFT JOIN competition_clubs cc ON cc.club_id=cl.id LEFT JOIN competitions c ON c.id=cc.competition_id WHERE $1='' OR c.provider_code=$1 ORDER BY cl.name`, competitionCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Club, 0)
	for rows.Next() {
		item, err := scanClub(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) GetClub(ctx context.Context, id uuid.UUID) (Club, error) {
	club, err := scanClub(r.pool.QueryRow(ctx, `SELECT id,provider_id,name,short_name,tla,crest_url,website_url,venue FROM clubs WHERE id=$1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return Club{}, ErrClubNotFound
	}
	if err != nil {
		return Club{}, err
	}
	rows, err := r.pool.Query(ctx, `SELECT c.id,c.provider_code,c.name,c.country,c.emblem_url,0 FROM competitions c JOIN competition_clubs cc ON cc.competition_id=c.id WHERE cc.club_id=$1 ORDER BY c.name`, id)
	if err != nil {
		return Club{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var item Competition
		if err := rows.Scan(&item.ID, &item.ProviderCode, &item.Name, &item.Country, &item.EmblemURL, &item.ClubCount); err != nil {
			return Club{}, err
		}
		club.Competitions = append(club.Competitions, item)
	}
	return club, rows.Err()
}

func (r *PostgresRepository) ListFavorites(ctx context.Context, userID uuid.UUID) ([]Club, error) {
	rows, err := r.pool.Query(ctx, `SELECT cl.id,cl.provider_id,cl.name,cl.short_name,cl.tla,cl.crest_url,cl.website_url,cl.venue FROM favorite_clubs f JOIN clubs cl ON cl.id=f.club_id WHERE f.user_id=$1 ORDER BY f.position`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Club, 0)
	for rows.Next() {
		item, err := scanClub(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) ReplaceFavorites(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) ([]Club, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `DELETE FROM favorite_clubs WHERE user_id=$1`, userID); err != nil {
		return nil, err
	}
	for index, id := range ids {
		result, execErr := tx.Exec(ctx, `INSERT INTO favorite_clubs(user_id,club_id,position) SELECT $1,id,$3 FROM clubs WHERE id=$2`, userID, id, index+1)
		if execErr != nil {
			return nil, execErr
		}
		if result.RowsAffected() == 0 {
			return nil, ErrClubNotFound
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.ListFavorites(ctx, userID)
}

func (r *PostgresRepository) UpsertCompetition(ctx context.Context, competition Competition, clubs []Club) (int, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	var competitionID uuid.UUID
	err = tx.QueryRow(ctx, `INSERT INTO competitions(id,provider,provider_code,name,country,emblem_url) VALUES($1,'football-data.org',$2,$3,$4,$5) ON CONFLICT(provider,provider_code) DO UPDATE SET name=excluded.name,country=excluded.country,emblem_url=excluded.emblem_url,updated_at=now() RETURNING id`, uuid.New(), competition.ProviderCode, competition.Name, competition.Country, competition.EmblemURL).Scan(&competitionID)
	if err != nil {
		return 0, err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM competition_clubs WHERE competition_id=$1`, competitionID); err != nil {
		return 0, err
	}
	for _, club := range clubs {
		var clubID uuid.UUID
		err = tx.QueryRow(ctx, `INSERT INTO clubs(id,provider,provider_id,name,short_name,tla,crest_url,website_url,venue) VALUES($1,'football-data.org',$2,$3,$4,$5,$6,$7,$8) ON CONFLICT(provider,provider_id) DO UPDATE SET name=excluded.name,short_name=excluded.short_name,tla=excluded.tla,crest_url=excluded.crest_url,website_url=excluded.website_url,venue=excluded.venue,updated_at=now() RETURNING id`, uuid.New(), club.ProviderID, club.Name, club.ShortName, club.TLA, club.CrestURL, club.WebsiteURL, club.Venue).Scan(&clubID)
		if err != nil {
			return 0, err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO competition_clubs(competition_id,club_id) VALUES($1,$2)`, competitionID, clubID); err != nil {
			return 0, err
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return len(clubs), nil
}

func (r *PostgresRepository) StartCollection(ctx context.Context) (CollectionRun, error) {
	item := CollectionRun{ID: uuid.New(), Status: "running"}
	err := r.pool.QueryRow(ctx, `INSERT INTO collection_runs(id,connector,status) VALUES($1,'football-data.org','running') RETURNING started_at`, item.ID).Scan(&item.StartedAt)
	return item, err
}

func (r *PostgresRepository) FinishCollection(ctx context.Context, id uuid.UUID, status string, count int, message string) error {
	_, err := r.pool.Exec(ctx, `UPDATE collection_runs SET status=$2,records_count=$3,error_message=$4,finished_at=now() WHERE id=$1`, id, status, count, message)
	return err
}

func (r *PostgresRepository) LatestCollection(ctx context.Context) (*CollectionRun, error) {
	var item CollectionRun
	err := r.pool.QueryRow(ctx, `SELECT id,status,records_count,error_message,started_at,finished_at FROM collection_runs WHERE connector='football-data.org' ORDER BY started_at DESC LIMIT 1`).Scan(&item.ID, &item.Status, &item.RecordsCount, &item.ErrorMessage, &item.StartedAt, &item.FinishedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &item, err
}

type rowScanner interface{ Scan(...any) error }

func scanClub(row rowScanner) (Club, error) {
	var item Club
	err := row.Scan(&item.ID, &item.ProviderID, &item.Name, &item.ShortName, &item.TLA, &item.CrestURL, &item.WebsiteURL, &item.Venue)
	return item, err
}
