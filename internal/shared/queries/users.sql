-- name: CreateUser :one
INSERT INTO users (email, password, phone_number, first_name, last_name)
VALUES ($1, $2, $3, $4, $5)
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


-- name: UpdateUserStatus :one
UPDATE users SET status = $1 WHERE id = $2 RETURNING id, email, status;

-- name: ActiveUser :one
UPDATE users SET is_verified = TRUE WHERE id = $1 RETURNING id, email, status;

-- name: DeleteVerificationCode :exec
DELETE FROM email_verifications WHERE user_id = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserPassword :one
UPDATE users SET password = $2 WHERE id = $1 RETURNING id, email, username;

-- name: GetAllUsers :many
SELECT id, email, username, is_verified, phone_number, first_name, last_name, status, language, created_at, updated_at 
FROM users 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAllUsers :one
SELECT COUNT(*) FROM users;

-- name: UpdateUserProfile :one
UPDATE users 
SET 
    username = COALESCE(NULLIF($2, ''), username),
    phone_number = COALESCE(NULLIF($3, ''), phone_number),
    first_name = COALESCE(NULLIF($4, ''), first_name),
    last_name = COALESCE(NULLIF($5, ''), last_name),
    language = COALESCE(NULLIF($6, ''), language),
    updated_at = extract(epoch from now())
WHERE id = $1
RETURNING id, email, username, phone_number, first_name, last_name, language, updated_at;