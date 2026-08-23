-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE user_role AS ENUM ('user', 'admin');
CREATE TYPE theme_preference AS ENUM ('system', 'light', 'dark');

CREATE TABLE users (
    id uuid PRIMARY KEY,
    email citext NOT NULL UNIQUE,
    password_hash text NOT NULL,
    nickname varchar(40) NOT NULL,
    avatar_seed varchar(64) NOT NULL,
    theme theme_preference NOT NULL DEFAULT 'system',
    role user_role NOT NULL DEFAULT 'user',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE invitations (
    id uuid PRIMARY KEY,
    code_hash bytea NOT NULL UNIQUE,
    label varchar(100) NOT NULL DEFAULT '',
    expires_at timestamptz,
    max_uses integer NOT NULL CHECK (max_uses > 0),
    uses integer NOT NULL DEFAULT 0 CHECK (uses >= 0),
    created_by uuid REFERENCES users(id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    revoked_at timestamptz,
    CHECK (uses <= max_uses)
);

CREATE TABLE invitation_uses (
    invitation_id uuid NOT NULL REFERENCES invitations(id) ON DELETE RESTRICT,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    used_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (invitation_id, user_id)
);

CREATE TABLE sessions (
    id uuid PRIMARY KEY,
    token_hash bytea NOT NULL UNIQUE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    last_seen_at timestamptz NOT NULL DEFAULT now(),
    user_agent varchar(300) NOT NULL DEFAULT '',
    ip inet
);
CREATE INDEX sessions_user_id_idx ON sessions(user_id);
CREATE INDEX sessions_expires_at_idx ON sessions(expires_at);

CREATE TABLE password_reset_links (
    id uuid PRIMARY KEY,
    token_hash bytea NOT NULL UNIQUE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_by uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL,
    used_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE password_reset_links;
DROP TABLE sessions;
DROP TABLE invitation_uses;
DROP TABLE invitations;
DROP TABLE users;
DROP TYPE theme_preference;
DROP TYPE user_role;
