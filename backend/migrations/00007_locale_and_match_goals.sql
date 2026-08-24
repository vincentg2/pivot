-- +goose Up
ALTER TABLE users ADD COLUMN locale text NOT NULL DEFAULT 'fr' CHECK (locale IN ('fr', 'en'));

CREATE TABLE match_goals (
    match_id uuid NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    position smallint NOT NULL,
    minute smallint NOT NULL,
    injury_time smallint,
    goal_type text NOT NULL,
    team_provider_id integer NOT NULL,
    scorer_name text NOT NULL,
    PRIMARY KEY (match_id, position)
);

-- +goose Down
DROP TABLE match_goals;
ALTER TABLE users DROP COLUMN locale;
