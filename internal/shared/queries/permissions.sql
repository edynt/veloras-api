-- name: CreatePermission :exec
INSERT INTO permissions (name, description) VALUES ($1, $2);

-- name: UpdatePermission :exec
UPDATE permissions SET name = $1, description = $2 WHERE id = $3;

-- name: DeletePermission :exec
DELETE FROM permissions WHERE id = $1;

-- name: GetPermissionById :one
SELECT * FROM permissions WHERE id = $1;

-- name: GetPermissionByName :one
SELECT * FROM permissions WHERE name = $1;

-- name: GetPermissions :many
SELECT * FROM permissions;
