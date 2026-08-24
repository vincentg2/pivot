-- +goose Up
CREATE TABLE seasons (
    id uuid PRIMARY KEY,
    competition_id uuid NOT NULL REFERENCES competitions(id) ON DELETE CASCADE,
    provider_id integer,
    start_date date NOT NULL,
    end_date date NOT NULL,
    current boolean NOT NULL DEFAULT false,
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (competition_id, start_date)
);

CREATE TABLE matches (
    id uuid PRIMARY KEY,
    provider varchar(40) NOT NULL,
    provider_id integer NOT NULL,
    competition_id uuid NOT NULL REFERENCES competitions(id) ON DELETE RESTRICT,
    season_id uuid NOT NULL REFERENCES seasons(id) ON DELETE RESTRICT,
    utc_date timestamptz NOT NULL,
    status varchar(30) NOT NULL,
    stage varchar(40) NOT NULL DEFAULT '',
    matchday integer,
    home_club_id uuid REFERENCES clubs(id) ON DELETE SET NULL,
    away_club_id uuid REFERENCES clubs(id) ON DELETE SET NULL,
    home_name varchar(140) NOT NULL,
    away_name varchar(140) NOT NULL,
    home_score integer,
    away_score integer,
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (provider, provider_id)
);
CREATE INDEX matches_utc_date_idx ON matches(utc_date);
CREATE INDEX matches_home_club_date_idx ON matches(home_club_id, utc_date);
CREATE INDEX matches_away_club_date_idx ON matches(away_club_id, utc_date);

CREATE TABLE standings (
    season_id uuid NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
    standing_type varchar(20) NOT NULL DEFAULT 'TOTAL',
    position integer NOT NULL,
    club_id uuid NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    played integer NOT NULL DEFAULT 0,
    won integer NOT NULL DEFAULT 0,
    drawn integer NOT NULL DEFAULT 0,
    lost integer NOT NULL DEFAULT 0,
    goals_for integer NOT NULL DEFAULT 0,
    goals_against integer NOT NULL DEFAULT 0,
    goal_difference integer NOT NULL DEFAULT 0,
    points integer NOT NULL DEFAULT 0,
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (season_id, standing_type, club_id),
    UNIQUE (season_id, standing_type, position)
);

ALTER TABLE collection_runs ADD COLUMN scope varchar(30) NOT NULL DEFAULT 'catalog';
DROP INDEX collection_runs_connector_started_idx;
CREATE INDEX collection_runs_connector_scope_started_idx ON collection_runs(connector, scope, started_at DESC);

-- +goose Down
DROP INDEX collection_runs_connector_scope_started_idx;
ALTER TABLE collection_runs DROP COLUMN scope;
CREATE INDEX collection_runs_connector_started_idx ON collection_runs(connector, started_at DESC);
DROP TABLE standings;
DROP TABLE matches;
DROP TABLE seasons;
