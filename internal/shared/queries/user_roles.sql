-- name: AssignRoleToUser :exec
INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2);

-- name: GetRolesByUser :many
SELECT r.* FROM roles r
JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = $1;

-- name: UserHasRole :one
SELECT EXISTS(
    SELECT 1 FROM user_roles ur
    JOIN roles r ON r.id = ur.role_id
    WHERE ur.user_id = $1 AND r.name = $2
);
