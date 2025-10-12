-- name: AddToWishlist :one
INSERT INTO wishlist_items (user_id, product_id)
VALUES ($1, $2)
ON CONFLICT (user_id, product_id) DO NOTHING
RETURNING *;

-- name: GetWishlistItem :one
SELECT wi.*, p.title, p.price, p.original_price, p.images, p.stock, p.status,
       u.first_name || ' ' || u.last_name as seller_name,
       s.name as store_name, s.logo as store_logo, s.rating as store_rating
FROM wishlist_items wi
LEFT JOIN products p ON wi.product_id = p.id
LEFT JOIN users u ON p.seller_id = u.id
LEFT JOIN stores s ON s.user_id = p.seller_id
WHERE wi.user_id = $1 AND wi.product_id = $2;

-- name: ListWishlistItems :many
SELECT wi.*, p.title, p.price, p.original_price, p.images, p.stock, p.status,
       u.first_name || ' ' || u.last_name as seller_name,
       s.name as store_name, s.logo as store_logo, s.rating as store_rating
FROM wishlist_items wi
LEFT JOIN products p ON wi.product_id = p.id
LEFT JOIN users u ON p.seller_id = u.id
LEFT JOIN stores s ON s.user_id = p.seller_id
WHERE wi.user_id = $1
ORDER BY wi.created_at DESC
LIMIT $2 OFFSET $3;

-- name: RemoveFromWishlist :exec
DELETE FROM wishlist_items WHERE user_id = $1 AND product_id = $2;

-- name: ClearWishlist :exec
DELETE FROM wishlist_items WHERE user_id = $1;

-- name: GetWishlistItemCount :one
SELECT COUNT(*) FROM wishlist_items WHERE user_id = $1;

-- name: IsInWishlist :one
SELECT EXISTS(
    SELECT 1 FROM wishlist_items 
    WHERE user_id = $1 AND product_id = $2
) as is_in_wishlist;