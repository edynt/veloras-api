-- name: CreateRole :one
INSERT INTO roles (role_name, role_description) VALUES ($1, $2) RETURNING *;

-- name: GetRoleByName :one
SELECT * FROM roles WHERE role_name = $1;
