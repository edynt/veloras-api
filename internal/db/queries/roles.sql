-- name: CreateRole :one
INSERT INTO roles (name, description) VALUES ($1, $2) RETURNING *;

-- name: GetRoleByName :one
SELECT * FROM roles WHERE name = $1;
