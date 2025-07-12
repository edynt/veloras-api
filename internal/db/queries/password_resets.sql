-- name: CreatePasswordReset :one
INSERT INTO password_resets (user_id, reset_token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPasswordReset :one
SELECT * FROM password_resets WHERE user_id = $1 AND reset_token = $2;
