-- name: CreateAuditLog :one
INSERT INTO audit_logs (user_id, action, resource_type, resource_id, details, ip_address, user_agent)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAuditLogsByUser :many
SELECT * FROM audit_logs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2;

-- name: GetAuditLogsByResource :many
SELECT * FROM audit_logs WHERE resource_type = $1 AND resource_id = $2 ORDER BY created_at DESC;

-- name: GetAuditLogsByAction :many
SELECT * FROM audit_logs WHERE action = $1 ORDER BY created_at DESC LIMIT $2;

-- name: GetAuditLogsByDateRange :many
SELECT * FROM audit_logs WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC;

-- name: DeleteOldAuditLogs :exec
DELETE FROM audit_logs WHERE created_at < $1;
