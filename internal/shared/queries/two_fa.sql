-- name: Create2FACode :one
INSERT INTO two_fa_codes (user_id, code, expires_at, type)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: Get2FACode :one
SELECT * FROM two_fa_codes WHERE user_id = $1 AND code = $2 AND type = $3;

-- name: Delete2FACode :exec
DELETE FROM two_fa_codes WHERE user_id = $1 AND type = $2;

-- name: DeleteExpired2FACodes :exec
DELETE FROM two_fa_codes WHERE expires_at < $1;

-- name: GetUser2FACodes :many
SELECT * FROM two_fa_codes WHERE user_id = $1 AND type = $2;

-- name: Update2FAEnabled :exec
UPDATE users SET two_fa_enabled = $1, two_fa_secret = $2 WHERE id = $3;
