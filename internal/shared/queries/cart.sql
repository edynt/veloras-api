-- name: AddToCart :one
INSERT INTO cart_items (user_id, product_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, product_id) 
DO UPDATE SET quantity = cart_items.quantity + $3, updated_at = NOW()
RETURNING *;

-- name: GetCartItem :one
SELECT ci.*, p.title, p.price, p.original_price, p.images, p.stock, p.status,
       u.first_name || ' ' || u.last_name as seller_name,
       s.name as store_name, s.logo as store_logo, s.rating as store_rating
FROM cart_items ci
LEFT JOIN products p ON ci.product_id = p.id
LEFT JOIN users u ON p.seller_id = u.id
LEFT JOIN stores s ON s.user_id = p.seller_id
WHERE ci.user_id = $1 AND ci.product_id = $2;

-- name: ListCartItems :many
SELECT ci.*, p.title, p.price, p.original_price, p.images, p.stock, p.status,
       u.first_name || ' ' || u.last_name as seller_name,
       s.name as store_name, s.logo as store_logo, s.rating as store_rating
FROM cart_items ci
LEFT JOIN products p ON ci.product_id = p.id
LEFT JOIN users u ON p.seller_id = u.id
LEFT JOIN stores s ON s.user_id = p.seller_id
WHERE ci.user_id = $1
ORDER BY ci.created_at DESC;

-- name: UpdateCartItemQuantity :one
UPDATE cart_items 
SET quantity = $3, updated_at = NOW()
WHERE user_id = $1 AND product_id = $2
RETURNING *;

-- name: RemoveFromCart :exec
DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE user_id = $1;

-- name: GetCartItemCount :one
SELECT COUNT(*) FROM cart_items WHERE user_id = $1;

-- name: GetCartTotal :one
SELECT COALESCE(SUM(ci.quantity * p.price), 0) as total
FROM cart_items ci
LEFT JOIN products p ON ci.product_id = p.id
WHERE ci.user_id = $1 AND p.status = 'active';