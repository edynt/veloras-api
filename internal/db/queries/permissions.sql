-- name: CreatePermission :one
INSERT INTO permissions (permission_name, permission_description) VALUES ($1, $2) RETURNING *;

-- name: GetPermissionByName :one
SELECT * FROM permissions WHERE permission_name = $1;
