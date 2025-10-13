-- name: AssignPermissionToRole :exec
INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2);

-- name: GetPermissionsByRole :many
SELECT p.* FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
WHERE rp.role_id = $1;

-- name: UserHasPermission :one
SELECT EXISTS(
    SELECT 1 FROM user_roles ur
    JOIN role_permissions rp ON rp.role_id = ur.role_id
    JOIN permissions p ON p.id = rp.permission_id
    WHERE ur.user_id = $1 AND p.name = $2
);
