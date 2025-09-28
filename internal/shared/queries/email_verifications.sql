-- name: CreateEmailVerification :one
INSERT INTO email_verifications (user_id, code, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetEmailVerification :one
SELECT * FROM email_verifications WHERE user_id = $1 AND code = $2;

-- name: DeleteExpiredEmailVerifications :exec
DELETE FROM email_verifications WHERE expires_at < $1;

-- name: CountExpiredEmailVerifications :one
SELECT COUNT(*) FROM email_verifications WHERE expires_at < $1;