-- +goose Up
CREATE TABLE news_feeds (
    id uuid PRIMARY KEY,
    club_id uuid NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    feed_url text NOT NULL,
    source_name varchar(120) NOT NULL,
    enabled boolean NOT NULL DEFAULT true,
    created_by uuid REFERENCES users(id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (club_id, feed_url)
);

CREATE TABLE news_items (
    id uuid PRIMARY KEY,
    feed_id uuid NOT NULL REFERENCES news_feeds(id) ON DELETE CASCADE,
    club_id uuid NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    external_id text NOT NULL,
    title varchar(300) NOT NULL,
    source_name varchar(120) NOT NULL,
    published_at timestamptz NOT NULL,
    link_url text NOT NULL,
    imported_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (feed_id, external_id)
);

CREATE INDEX news_items_published_idx ON news_items (published_at DESC);
CREATE INDEX news_items_club_published_idx ON news_items (club_id, published_at DESC);

-- +goose Down
DROP TABLE IF EXISTS news_items;
DROP TABLE IF EXISTS news_feeds;
