-- name: CreateEmailVerification :one
INSERT INTO email_verifications (ev_user_id, ev_code, ev_expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetEmailVerification :one
SELECT * FROM email_verifications WHERE ev_user_id = $1 AND ev_code = $2;
