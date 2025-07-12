-- name: CreateSession :one
INSERT INTO sessions (session_user_id, session_refresh_token, session_user_agent, session_client_ip, session_expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions WHERE session_id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE session_id = $1;
