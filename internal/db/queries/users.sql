-- name: CreateUser :one
INSERT INTO users (user_email, user_username, user_hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE user_email = $1;

-- name: VerifyUser :exec
UPDATE users SET user_is_verified = TRUE WHERE user_id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE user_id = $1;
