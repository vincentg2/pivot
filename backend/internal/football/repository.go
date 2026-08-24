package football

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

func (r *PostgresRepository) ListMatches(ctx context.Context, userID uuid.UUID, from, to time.Time, competitionCode string, clubID *uuid.UUID) ([]Match, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT m.id,m.provider_id,m.utc_date,m.status,m.stage,m.matchday,m.home_score,m.away_score,
		       c.id,c.provider_code,c.name,s.provider_id,s.start_date,s.end_date,s.current,
		       hc.id,m.home_name,coalesce(hc.short_name,''),coalesce(hc.tla,''),hc.crest_url,
		       ac.id,m.away_name,coalesce(ac.short_name,''),coalesce(ac.tla,''),ac.crest_url,
		       EXISTS(SELECT 1 FROM favorite_clubs f WHERE f.user_id=$1 AND f.club_id IN (m.home_club_id,m.away_club_id))
		FROM matches m
		JOIN competitions c ON c.id=m.competition_id
		JOIN seasons s ON s.id=m.season_id
		LEFT JOIN clubs hc ON hc.id=m.home_club_id
		LEFT JOIN clubs ac ON ac.id=m.away_club_id
		WHERE (m.utc_date AT TIME ZONE 'Europe/Paris')::date BETWEEN $2::date AND $3::date
		  AND ($4='' OR c.provider_code=$4)
		  AND ($5::uuid IS NULL OR m.home_club_id=$5 OR m.away_club_id=$5)
		ORDER BY m.utc_date,c.name,m.provider_id`, userID, from, to, competitionCode, clubID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Match, 0)
	for rows.Next() {
		var item Match
		if err := rows.Scan(&item.ID, &item.ProviderID, &item.UTCDate, &item.Status, &item.Stage, &item.Matchday, &item.HomeScore, &item.AwayScore,
			&item.Competition.ID, &item.Competition.Code, &item.Competition.Name, &item.Season.ProviderID, &item.Season.StartDate, &item.Season.EndDate, &item.Season.Current,
			&item.Home.ID, &item.Home.Name, &item.Home.ShortName, &item.Home.TLA, &item.Home.CrestURL,
			&item.Away.ID, &item.Away.Name, &item.Away.ShortName, &item.Away.TLA, &item.Away.CrestURL, &item.Favorite); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) CurrentStanding(ctx context.Context, competitionCode string) (Standing, error) {
	var result Standing
	err := r.pool.QueryRow(ctx, `SELECT c.id,c.provider_code,c.name,s.provider_id,s.start_date,s.end_date,s.current FROM competitions c JOIN seasons s ON s.competition_id=c.id AND s.current WHERE c.provider_code=$1`, competitionCode).Scan(&result.Competition.ID, &result.Competition.Code, &result.Competition.Name, &result.Season.ProviderID, &result.Season.StartDate, &result.Season.EndDate, &result.Season.Current)
	if errors.Is(err, pgx.ErrNoRows) {
		return Standing{}, ErrCompetitionMissing
	}
	if err != nil {
		return Standing{}, err
	}
	rows, err := r.pool.Query(ctx, `SELECT st.position,cl.id,cl.name,cl.short_name,cl.tla,cl.crest_url,st.played,st.won,st.drawn,st.lost,st.goals_for,st.goals_against,st.goal_difference,st.points FROM standings st JOIN seasons s ON s.id=st.season_id JOIN clubs cl ON cl.id=st.club_id WHERE s.competition_id=$1 AND s.current AND st.standing_type='TOTAL' ORDER BY st.position,cl.name`, result.Competition.ID)
	if err != nil {
		return Standing{}, err
	}
	defer rows.Close()
	result.Rows = make([]StandingRow, 0)
	for rows.Next() {
		var row StandingRow
		if err := rows.Scan(&row.Position, &row.Club.ID, &row.Club.Name, &row.Club.ShortName, &row.Club.TLA, &row.Club.CrestURL, &row.Played, &row.Won, &row.Drawn, &row.Lost, &row.GoalsFor, &row.GoalsAgainst, &row.GoalDifference, &row.Points); err != nil {
			return Standing{}, err
		}
		result.Rows = append(result.Rows, row)
	}
	return result, rows.Err()
}

func (r *PostgresRepository) SaveCompetitionData(ctx context.Context, code string, season Season, matches []Match, standings []StandingRow) (int, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	var competitionID uuid.UUID
	if err = tx.QueryRow(ctx, `SELECT id FROM competitions WHERE provider='football-data.org' AND provider_code=$1`, code).Scan(&competitionID); errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrCompetitionMissing
	}
	if err != nil {
		return 0, err
	}
	var seasonID uuid.UUID
	err = tx.QueryRow(ctx, `INSERT INTO seasons(id,competition_id,provider_id,start_date,end_date,current) VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT(competition_id,start_date) DO UPDATE SET provider_id=excluded.provider_id,end_date=excluded.end_date,current=excluded.current,updated_at=now() RETURNING id`, uuid.New(), competitionID, season.ProviderID, season.StartDate, season.EndDate, season.Current).Scan(&seasonID)
	if err != nil {
		return 0, err
	}
	if season.Current {
		if _, err = tx.Exec(ctx, `UPDATE seasons SET current=false WHERE competition_id=$1 AND id<>$2`, competitionID, seasonID); err != nil {
			return 0, err
		}
	}
	for _, item := range matches {
		_, err = tx.Exec(ctx, `INSERT INTO matches(id,provider,provider_id,competition_id,season_id,utc_date,status,stage,matchday,home_club_id,away_club_id,home_name,away_name,home_score,away_score) VALUES($1,'football-data.org',$2,$3,$4,$5,$6,$7,$8,(SELECT id FROM clubs WHERE provider='football-data.org' AND provider_id=$9),(SELECT id FROM clubs WHERE provider='football-data.org' AND provider_id=$10),$11,$12,$13,$14) ON CONFLICT(provider,provider_id) DO UPDATE SET competition_id=excluded.competition_id,season_id=excluded.season_id,utc_date=excluded.utc_date,status=excluded.status,stage=excluded.stage,matchday=excluded.matchday,home_club_id=excluded.home_club_id,away_club_id=excluded.away_club_id,home_name=excluded.home_name,away_name=excluded.away_name,home_score=excluded.home_score,away_score=excluded.away_score,updated_at=now()`, uuid.New(), item.ProviderID, competitionID, seasonID, item.UTCDate, item.Status, item.Stage, item.Matchday, item.HomeProviderID, item.AwayProviderID, item.Home.Name, item.Away.Name, item.HomeScore, item.AwayScore)
		if err != nil {
			return 0, err
		}
	}
	if len(standings) > 0 {
		if _, err = tx.Exec(ctx, `DELETE FROM standings WHERE season_id=$1 AND standing_type='TOTAL'`, seasonID); err != nil {
			return 0, err
		}
		for _, row := range standings {
			result, insertErr := tx.Exec(ctx, `INSERT INTO standings(season_id,standing_type,position,club_id,played,won,drawn,lost,goals_for,goals_against,goal_difference,points) SELECT $1,'TOTAL',$2,id,$4,$5,$6,$7,$8,$9,$10,$11 FROM clubs WHERE provider='football-data.org' AND provider_id=$3`, seasonID, row.Position, row.ClubProviderID, row.Played, row.Won, row.Drawn, row.Lost, row.GoalsFor, row.GoalsAgainst, row.GoalDifference, row.Points)
			if insertErr != nil {
				return 0, insertErr
			}
			if result.RowsAffected() == 0 {
				return 0, ErrCompetitionMissing
			}
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return len(matches) + len(standings), nil
}

func (r *PostgresRepository) StartCollection(ctx context.Context) (CollectionRun, error) {
	item := CollectionRun{ID: uuid.New(), Status: "running"}
	err := r.pool.QueryRow(ctx, `INSERT INTO collection_runs(id,connector,scope,status) VALUES($1,'football-data.org','sport','running') RETURNING started_at`, item.ID).Scan(&item.StartedAt)
	return item, err
}
func (r *PostgresRepository) FinishCollection(ctx context.Context, id uuid.UUID, status string, count int, message string) error {
	_, err := r.pool.Exec(ctx, `UPDATE collection_runs SET status=$2,records_count=$3,error_message=$4,finished_at=now() WHERE id=$1`, id, status, count, message)
	return err
}
func (r *PostgresRepository) LatestCollection(ctx context.Context) (*CollectionRun, error) {
	var item CollectionRun
	err := r.pool.QueryRow(ctx, `SELECT id,status,records_count,error_message,started_at,finished_at FROM collection_runs WHERE connector='football-data.org' AND scope='sport' ORDER BY started_at DESC LIMIT 1`).Scan(&item.ID, &item.Status, &item.RecordsCount, &item.ErrorMessage, &item.StartedAt, &item.FinishedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &item, err
}
