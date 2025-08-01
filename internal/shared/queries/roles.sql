-- name: CreateRole :exec
INSERT INTO roles (name, description) VALUES ($1, $2);

-- name: UpdateRole :exec
UPDATE roles SET name = $1, description = $2 WHERE id = $3;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;

-- name: GetRoleById :one
SELECT * FROM roles WHERE id = $1;

-- name: GetRoles :many
SELECT * FROM roles;