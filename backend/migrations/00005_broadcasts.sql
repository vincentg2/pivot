-- +goose Up
CREATE TABLE tv_listings (
    id uuid PRIMARY KEY,
    source text NOT NULL CHECK (source IN ('footao')),
    source_key text NOT NULL,
    match_id uuid REFERENCES matches(id) ON DELETE SET NULL,
    starts_at timestamptz NOT NULL,
    home_name text NOT NULL DEFAULT '',
    away_name text NOT NULL DEFAULT '',
    label text NOT NULL,
    competition_name text NOT NULL DEFAULT '',
    broadcast_kind text NOT NULL DEFAULT 'live' CHECK (broadcast_kind IN ('live', 'delayed', 'replay')),
    source_url text,
    imported_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (source, source_key)
);

CREATE INDEX tv_listings_starts_at_idx ON tv_listings (starts_at);
CREATE INDEX tv_listings_match_id_idx ON tv_listings (match_id) WHERE match_id IS NOT NULL;

CREATE TABLE tv_listing_channels (
    listing_id uuid NOT NULL REFERENCES tv_listings(id) ON DELETE CASCADE,
    channel text NOT NULL,
    position smallint NOT NULL DEFAULT 0,
    PRIMARY KEY (listing_id, channel)
);

CREATE TABLE tv_listing_corrections (
    listing_id uuid PRIMARY KEY REFERENCES tv_listings(id) ON DELETE CASCADE,
    starts_at timestamptz NOT NULL,
    home_name text NOT NULL DEFAULT '',
    away_name text NOT NULL DEFAULT '',
    label text NOT NULL,
    competition_name text NOT NULL DEFAULT '',
    broadcast_kind text NOT NULL CHECK (broadcast_kind IN ('live', 'delayed', 'replay')),
    channels text[] NOT NULL DEFAULT '{}',
    hidden boolean NOT NULL DEFAULT false,
    note text NOT NULL DEFAULT '',
    updated_by uuid REFERENCES users(id) ON DELETE SET NULL,
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE tv_listing_audit (
    id uuid PRIMARY KEY,
    listing_id uuid NOT NULL,
    admin_id uuid REFERENCES users(id) ON DELETE SET NULL,
    action text NOT NULL CHECK (action IN ('corrected', 'cleared')),
    before_state jsonb,
    after_state jsonb,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX tv_listing_audit_listing_idx ON tv_listing_audit (listing_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS tv_listing_audit;
DROP TABLE IF EXISTS tv_listing_corrections;
DROP TABLE IF EXISTS tv_listing_channels;
DROP TABLE IF EXISTS tv_listings;
