-- name: CreateUser :one
INSERT INTO users (email, username, password, phone_number, first_name, last_name)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, email, status;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: VerifyUser :exec
UPDATE users SET is_verified = TRUE WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: GetUserEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: GetUsernameExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);

-- name: UpdateUserStatus :one
UPDATE users SET status = $1 WHERE id = $2 RETURNING id, email, status;

-- name: ActiveUser :one
UPDATE users SET is_verified = TRUE WHERE id = $1 RETURNING id, email, status;

-- name: DeleteVerificationCode :exec
DELETE FROM email_verifications WHERE user_id = $1;