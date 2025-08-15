-- name: CreateSession :one
INSERT INTO sessions (user_id, refresh_token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions WHERE id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;

-- name: GetSessionByRefreshToken :one
SELECT * FROM sessions WHERE refresh_token = $1;

-- name: GetSessionsByUser :many
SELECT * FROM sessions WHERE user_id = $1;

-- name: DeleteSessionByRefreshToken :exec
DELETE FROM sessions WHERE refresh_token = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expires_at < $1;

-- name: DeleteAllUserSessions :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: UpdateSessionExpiry :exec
UPDATE sessions SET expires_at = $1 WHERE id = $2;