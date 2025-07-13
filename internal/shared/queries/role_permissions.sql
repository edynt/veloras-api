-- name: AssignPermissionToRole :exec
INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2);

-- name: GetPermissionsByRole :many
SELECT p.* FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.permission_id
WHERE rp.role_id = $1;
