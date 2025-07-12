-- name: AssignRoleToUser :exec
INSERT INTO user_roles (ur_user_id, ur_role_id) VALUES ($1, $2);

-- name: GetRolesByUser :many
SELECT r.* FROM roles r
JOIN user_roles ur ON ur.ur_role_id = r.role_id
WHERE ur.ur_user_id = $1;
