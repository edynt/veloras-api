-- name: CreatePasswordReset :one
INSERT INTO password_resets (ps_user_id, ps_reset_token, ps_expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPasswordReset :one
SELECT * FROM password_resets WHERE ps_user_id = $1 AND ps_reset_token = $2;
