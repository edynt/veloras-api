-- name: CreateOrganization :one
INSERT INTO organizations (name, description, parent_id, created_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrganizationById :one
SELECT * FROM organizations WHERE id = $1;

-- name: GetOrganizationsByParent :many
SELECT * FROM organizations WHERE parent_id = $1;

-- name: UpdateOrganization :exec
UPDATE organizations SET name = $1, description = $2, updated_at = $3 WHERE id = $4;

-- name: DeleteOrganization :exec
DELETE FROM organizations WHERE id = $1;

-- name: GetUserOrganization :one
SELECT o.* FROM organizations o
JOIN users u ON u.organization_id = o.id
WHERE u.id = $1;

-- name: AssignUserToOrganization :exec
UPDATE users SET organization_id = $1 WHERE id = $2;

-- name: GetOrganizationUsers :many
SELECT u.* FROM users u WHERE u.organization_id = $1;
