-- name: CreatePermission :exec
INSERT INTO permissions (name, description, resource_type, resource_action, organization_id)
VALUES ($1, $2, $3, $4, $5);

-- name: GetPermissionById :one
SELECT * FROM permissions WHERE id = $1;

-- name: GetPermissionByName :one
SELECT * FROM permissions WHERE name = $1;

-- name: GetPermissions :many
SELECT * FROM permissions;

-- name: UpdatePermission :exec
UPDATE permissions SET name = $1, description = $2, resource_type = $3, resource_action = $4 WHERE id = $5;

-- name: DeletePermission :exec
DELETE FROM permissions WHERE id = $1;

-- name: GetPermissionsByResourceType :many
SELECT * FROM permissions WHERE resource_type = $1;

-- name: GetPermissionsByOrganization :many
SELECT * FROM permissions WHERE organization_id = $1;

-- name: GetPermissionsByAction :many
SELECT * FROM permissions WHERE resource_action = $1;

-- name: GetUserPermissions :many
SELECT DISTINCT p.* FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
JOIN user_roles ur ON rp.role_id = ur.role_id
WHERE ur.user_id = $1;

-- name: CheckUserPermission :one
SELECT EXISTS(
    SELECT 1 FROM permissions p
    JOIN role_permissions rp ON p.id = rp.permission_id
    JOIN user_roles ur ON rp.role_id = ur.role_id
    WHERE ur.user_id = $1 AND p.resource_type = $2 AND p.resource_action = $3
);
