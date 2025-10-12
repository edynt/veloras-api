-- name: CreateAddress :one
INSERT INTO addresses (
    user_id, street, city, state, zip_code, country, is_default
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetAddress :one
SELECT * FROM addresses WHERE id = $1;

-- name: ListUserAddresses :many
SELECT * FROM addresses 
WHERE user_id = $1
ORDER BY is_default DESC, created_at DESC;

-- name: GetDefaultAddress :one
SELECT * FROM addresses 
WHERE user_id = $1 AND is_default = true
LIMIT 1;

-- name: UpdateAddress :one
UPDATE addresses 
SET street = $2, city = $3, state = $4, zip_code = $5, 
    country = $6, is_default = $7, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SetDefaultAddress :exec
UPDATE addresses 
SET is_default = CASE WHEN id = $2 THEN true ELSE false END,
    updated_at = NOW()
WHERE user_id = $1;

-- name: DeleteAddress :exec
DELETE FROM addresses WHERE id = $1;

-- name: GetAddressCount :one
SELECT COUNT(*) FROM addresses WHERE user_id = $1;