-- name: CreateOrder :one
INSERT INTO orders (
    buyer_id, seller_id, product_id, quantity, total_amount, status,
    payment_method, payment_status, shipping_address_id, notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetOrder :one
SELECT o.*, 
       p.title as product_title, p.price as product_price, p.images as product_images,
       b.first_name || ' ' || b.last_name as buyer_name, b.email as buyer_email,
       s.first_name || ' ' || s.last_name as seller_name, s.email as seller_email,
       a.street, a.city, a.state, a.zip_code, a.country
FROM orders o
LEFT JOIN products p ON o.product_id = p.id
LEFT JOIN users b ON o.buyer_id = b.id
LEFT JOIN users s ON o.seller_id = s.id
LEFT JOIN addresses a ON o.shipping_address_id = a.id
WHERE o.id = $1;

-- name: ListOrdersByBuyer :many
SELECT o.*, 
       p.title as product_title, p.price as product_price, p.images as product_images,
       s.first_name || ' ' || s.last_name as seller_name,
       st.name as store_name, st.logo as store_logo
FROM orders o
LEFT JOIN products p ON o.product_id = p.id
LEFT JOIN users s ON o.seller_id = s.id
LEFT JOIN stores st ON st.user_id = o.seller_id
WHERE o.buyer_id = $1
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListOrdersBySeller :many
SELECT o.*, 
       p.title as product_title, p.price as product_price, p.images as product_images,
       b.first_name || ' ' || b.last_name as buyer_name, b.email as buyer_email
FROM orders o
LEFT JOIN products p ON o.product_id = p.id
LEFT JOIN users b ON o.buyer_id = b.id
WHERE o.seller_id = $1
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListOrdersByStatus :many
SELECT o.*, 
       p.title as product_title, p.price as product_price, p.images as product_images,
       b.first_name || ' ' || b.last_name as buyer_name,
       s.first_name || ' ' || s.last_name as seller_name
FROM orders o
LEFT JOIN products p ON o.product_id = p.id
LEFT JOIN users b ON o.buyer_id = b.id
LEFT JOIN users s ON o.seller_id = s.id
WHERE o.status = $1
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateOrderStatus :one
UPDATE orders 
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateOrderPaymentStatus :one
UPDATE orders 
SET payment_status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateOrder :one
UPDATE orders 
SET status = $2, payment_status = $3, notes = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = $1;

-- name: GetOrderStats :one
SELECT 
    COUNT(*) as total_orders,
    COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_orders,
    COUNT(CASE WHEN status = 'confirmed' THEN 1 END) as confirmed_orders,
    COUNT(CASE WHEN status = 'shipped' THEN 1 END) as shipped_orders,
    COUNT(CASE WHEN status = 'delivered' THEN 1 END) as delivered_orders,
    COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled_orders,
    COUNT(CASE WHEN payment_status = 'completed' THEN 1 END) as completed_payments,
    SUM(CASE WHEN payment_status = 'completed' THEN total_amount ELSE 0 END) as total_revenue
FROM orders;

-- name: GetOrdersByDateRange :many
SELECT o.*, 
       p.title as product_title, p.price as product_price,
       b.first_name || ' ' || b.last_name as buyer_name,
       s.first_name || ' ' || s.last_name as seller_name
FROM orders o
LEFT JOIN products p ON o.product_id = p.id
LEFT JOIN users b ON o.buyer_id = b.id
LEFT JOIN users s ON o.seller_id = s.id
WHERE o.created_at >= $1 AND o.created_at <= $2
ORDER BY o.created_at DESC
LIMIT $3 OFFSET $4;