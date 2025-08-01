-- name: CreatePasswordReset :exec
INSERT INTO password_resets (user_id, reset_token, expires_at)
VALUES ($1, $2, $3);

-- name: GetPasswordReset :one
SELECT * FROM password_resets WHERE user_id = $1 AND reset_token = $2;

-- name: DeletePasswordReset :exec
DELETE FROM password_resets WHERE user_id = $1 AND reset_token = $2;
