-- name: CreateStore :one
INSERT INTO stores (
    user_id, name, description, logo, banner, address, phone, email, is_verified
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetStore :one
SELECT s.*, u.first_name || ' ' || u.last_name as owner_name, u.email as owner_email
FROM stores s
LEFT JOIN users u ON s.user_id = u.id
WHERE s.id = $1;

-- name: GetStoreByUserID :one
SELECT s.*, u.first_name || ' ' || u.last_name as owner_name, u.email as owner_email
FROM stores s
LEFT JOIN users u ON s.user_id = u.id
WHERE s.user_id = $1;

-- name: ListStores :many
SELECT s.*, u.first_name || ' ' || u.last_name as owner_name, u.email as owner_email
FROM stores s
LEFT JOIN users u ON s.user_id = u.id
WHERE s.is_verified = true
ORDER BY s.rating DESC, s.review_count DESC
LIMIT $1 OFFSET $2;

-- name: ListVerifiedStores :many
SELECT s.*, u.first_name || ' ' || u.last_name as owner_name, u.email as owner_email
FROM stores s
LEFT JOIN users u ON s.user_id = u.id
WHERE s.is_verified = true
ORDER BY s.rating DESC, s.review_count DESC
LIMIT $1 OFFSET $2;

-- name: UpdateStore :one
UPDATE stores 
SET name = $2, description = $3, logo = $4, banner = $5, 
    address = $6, phone = $7, email = $8, is_verified = $9, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateStoreRating :exec
UPDATE stores 
SET rating = $2, review_count = $3, updated_at = NOW()
WHERE id = $1;

-- name: VerifyStore :one
UPDATE stores 
SET is_verified = true, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteStore :exec
DELETE FROM stores WHERE id = $1;

-- name: GetStoreStats :one
SELECT 
    COUNT(*) as total_stores,
    COUNT(CASE WHEN is_verified = true THEN 1 END) as verified_stores,
    AVG(rating) as average_rating,
    SUM(review_count) as total_reviews
FROM stores;