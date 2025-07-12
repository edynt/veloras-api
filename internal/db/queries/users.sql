-- name: CreateUser :one
INSERT INTO users (email, username, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: VerifyUser :exec
UPDATE users SET is_verified = TRUE WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
