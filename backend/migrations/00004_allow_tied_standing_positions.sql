-- +goose Up
ALTER TABLE standings DROP CONSTRAINT standings_season_id_standing_type_position_key;

-- +goose Down
ALTER TABLE standings ADD CONSTRAINT standings_season_id_standing_type_position_key UNIQUE (season_id, standing_type, position);
