-- name: AssignPermissionToRole :exec
INSERT INTO role_permissions (rp_role_id, rp_permission_id) VALUES ($1, $2);

-- name: GetPermissionsByRole :many
SELECT p.* FROM permissions p
JOIN role_permissions rp ON rp.rp_permission_id = p.permission_id
WHERE rp.rp_role_id = $1;
