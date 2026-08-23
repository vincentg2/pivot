-- name: GetUserByEmail :one
SELECT id, email, password_hash, nickname, avatar_seed, theme, role, created_at, updated_at
FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, password_hash, nickname, avatar_seed, theme, role, created_at, updated_at
FROM users WHERE id = $1;

-- name: DeleteExpiredSessions :execrows
DELETE FROM sessions WHERE expires_at < now();
