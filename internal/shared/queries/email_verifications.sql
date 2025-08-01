-- name: CreateEmailVerification :one
INSERT INTO email_verifications (user_id, code, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetEmailVerification :one
SELECT * FROM email_verifications WHERE user_id = $1 AND code = $2;