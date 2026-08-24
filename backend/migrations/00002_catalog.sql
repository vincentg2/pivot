-- +goose Up
CREATE TABLE competitions (
    id uuid PRIMARY KEY,
    provider varchar(40) NOT NULL,
    provider_code varchar(20) NOT NULL,
    name varchar(120) NOT NULL,
    country varchar(80) NOT NULL DEFAULT '',
    emblem_url text,
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (provider, provider_code)
);

CREATE TABLE clubs (
    id uuid PRIMARY KEY,
    provider varchar(40) NOT NULL,
    provider_id integer NOT NULL,
    name varchar(140) NOT NULL,
    short_name varchar(80) NOT NULL DEFAULT '',
    tla varchar(5) NOT NULL DEFAULT '',
    crest_url text,
    website_url text,
    venue varchar(140) NOT NULL DEFAULT '',
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (provider, provider_id)
);

CREATE TABLE competition_clubs (
    competition_id uuid NOT NULL REFERENCES competitions(id) ON DELETE CASCADE,
    club_id uuid NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    PRIMARY KEY (competition_id, club_id)
);

CREATE TABLE favorite_clubs (
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    club_id uuid NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    position smallint NOT NULL CHECK (position BETWEEN 1 AND 5),
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, club_id),
    UNIQUE (user_id, position)
);

CREATE TABLE collection_runs (
    id uuid PRIMARY KEY,
    connector varchar(40) NOT NULL,
    status varchar(20) NOT NULL CHECK (status IN ('running', 'succeeded', 'failed', 'disabled')),
    records_count integer NOT NULL DEFAULT 0,
    error_message varchar(300) NOT NULL DEFAULT '',
    started_at timestamptz NOT NULL DEFAULT now(),
    finished_at timestamptz
);
CREATE INDEX collection_runs_connector_started_idx ON collection_runs(connector, started_at DESC);

-- +goose Down
DROP TABLE collection_runs;
DROP TABLE favorite_clubs;
DROP TABLE competition_clubs;
DROP TABLE clubs;
DROP TABLE competitions;
