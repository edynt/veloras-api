-- name: CreateProduct :one
INSERT INTO products (
    title, description, price, original_price, category_id, images, condition,
    tags, location, stock, seller_id, status, badges, is_featured, is_promoted
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
) RETURNING *;

-- name: GetProduct :one
SELECT p.*, c.name as category_name, c.slug as category_slug,
       u.first_name || ' ' || u.last_name as seller_name, u.email as seller_email,
       s.name as store_name, s.logo as store_logo, s.rating as store_rating
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN users u ON p.seller_id = u.id
LEFT JOIN stores s ON s.user_id = p.seller_id
WHERE p.id = $1;

-- name: ListProducts :many
SELECT p.*, c.name as category_name, c.slug as category_slug,
       u.first_name || ' ' || u.last_name as seller_name,
       s.name as store_name, s.logo as store_logo, s.rating as store_rating
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN users u ON p.seller_id = u.id
LEFT JOIN stores s ON s.user_id = p.seller_id
WHERE p.status = 'active'
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListProductsByCategory :many
SELECT p.*, c.name as category_name, c.slug as category_slug,
       u.first_name || ' ' || u.last_name as seller_name,
       s.name as store_name, s.logo as store_logo, s.rating as store_rating
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN users u ON p.seller_id = u.id
LEFT JOIN stores s ON s.user_id = p.seller_id
WHERE p.category_id = $1 AND p.status = 'active'
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListProductsBySeller :many
SELECT p.*, c.name as category_name, c.slug as category_slug,
       u.first_name || ' ' || u.last_name as seller_name,
       s.name as store_name, s.logo as store_logo, s.rating as store_rating
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN users u ON p.seller_id = u.id
LEFT JOIN stores s ON s.user_id = p.seller_id
WHERE p.seller_id = $1 AND p.status = 'active'
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchProducts :many
SELECT p.*, c.name as category_name, c.slug as category_slug,
       u.first_name || ' ' || u.last_name as seller_name,
       s.name as store_name, s.logo as store_logo, s.rating as store_rating,
       ts_rank(to_tsvector('english', p.title || ' ' || p.description), plainto_tsquery('english', $1)) as rank
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN users u ON p.seller_id = u.id
LEFT JOIN stores s ON s.user_id = p.seller_id
WHERE p.status = 'active' 
  AND (to_tsvector('english', p.title || ' ' || p.description) @@ plainto_tsquery('english', $1)
       OR p.tags && string_to_array($1, ' '))
ORDER BY rank DESC, p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetFeaturedProducts :many
SELECT p.*, c.name as category_name, c.slug as category_slug,
       u.first_name || ' ' || u.last_name as seller_name,
       s.name as store_name, s.logo as store_logo, s.rating as store_rating
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN users u ON p.seller_id = u.id
LEFT JOIN stores s ON s.user_id = p.seller_id
WHERE p.is_featured = true AND p.status = 'active'
ORDER BY p.rating DESC, p.review_count DESC
LIMIT $1;

-- name: GetPromotedProducts :many
SELECT p.*, c.name as category_name, c.slug as category_slug,
       u.first_name || ' ' || u.last_name as seller_name,
       s.name as store_name, s.logo as store_logo, s.rating as store_rating
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN users u ON p.seller_id = u.id
LEFT JOIN stores s ON s.user_id = p.seller_id
WHERE p.is_promoted = true AND p.status = 'active'
ORDER BY p.created_at DESC
LIMIT $1;

-- name: UpdateProduct :one
UPDATE products 
SET title = $2, description = $3, price = $4, original_price = $5, 
    category_id = $6, images = $7, condition = $8, tags = $9, 
    location = $10, stock = $11, status = $12, badges = $13, 
    is_featured = $14, is_promoted = $15, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateProductStatus :one
UPDATE products 
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateProductStock :one
UPDATE products 
SET stock = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: IncrementProductViews :exec
UPDATE products 
SET view_count = view_count + 1
WHERE id = $1;

-- name: UpdateProductRating :exec
UPDATE products 
SET rating = $2, review_count = $3, updated_at = NOW()
WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: GetProductStats :one
SELECT 
    COUNT(*) as total_products,
    COUNT(CASE WHEN status = 'active' THEN 1 END) as active_products,
    COUNT(CASE WHEN status = 'draft' THEN 1 END) as draft_products,
    COUNT(CASE WHEN status = 'sold' THEN 1 END) as sold_products,
    COUNT(CASE WHEN is_featured = true THEN 1 END) as featured_products,
    COUNT(CASE WHEN is_promoted = true THEN 1 END) as promoted_products,
    AVG(rating) as average_rating,
    SUM(view_count) as total_views
FROM products;

-- name: GetProductsByFilters :many
SELECT p.*, c.name as category_name, c.slug as category_slug,
       u.first_name || ' ' || u.last_name as seller_name,
       s.name as store_name, s.logo as store_logo, s.rating as store_rating
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN users u ON p.seller_id = u.id
LEFT JOIN stores s ON s.user_id = p.seller_id
WHERE p.status = 'active'
  AND ($1::uuid IS NULL OR p.category_id = $1)
  AND ($2::decimal IS NULL OR p.price >= $2)
  AND ($3::decimal IS NULL OR p.price <= $3)
  AND ($4::varchar IS NULL OR p.condition = $4)
  AND ($5::varchar IS NULL OR p.location ILIKE '%' || $5 || '%')
  AND ($6::integer IS NULL OR p.seller_id = $6)
  AND ($7::decimal IS NULL OR p.rating >= $7)
ORDER BY 
  CASE WHEN $8 = 'price_asc' THEN p.price END ASC,
  CASE WHEN $8 = 'price_desc' THEN p.price END DESC,
  CASE WHEN $8 = 'rating_desc' THEN p.rating END DESC,
  CASE WHEN $8 = 'newest' THEN p.created_at END DESC,
  p.created_at DESC
LIMIT $9 OFFSET $10;