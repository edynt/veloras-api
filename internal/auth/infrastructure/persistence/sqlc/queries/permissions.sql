-- name: CreatePermission :one
INSERT INTO permissions (name, description) VALUES ($1, $2) RETURNING *;

-- name: GetPermissionByName :one
SELECT * FROM permissions WHERE name = $1;
