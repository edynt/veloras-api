-- name: CreateReview :one
INSERT INTO reviews (
    product_id, buyer_id, seller_id, rating, comment, images, is_verified
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetReview :one
SELECT r.*, 
       p.title as product_title,
       b.first_name || ' ' || b.last_name as buyer_name,
       s.first_name || ' ' || s.last_name as seller_name
FROM reviews r
LEFT JOIN products p ON r.product_id = p.id
LEFT JOIN users b ON r.buyer_id = b.id
LEFT JOIN users s ON r.seller_id = s.id
WHERE r.id = $1;

-- name: ListReviewsByProduct :many
SELECT r.*, 
       b.first_name || ' ' || b.last_name as buyer_name,
       b.email as buyer_email
FROM reviews r
LEFT JOIN users b ON r.buyer_id = b.id
WHERE r.product_id = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListReviewsByBuyer :many
SELECT r.*, 
       p.title as product_title, p.images as product_images,
       s.first_name || ' ' || s.last_name as seller_name
FROM reviews r
LEFT JOIN products p ON r.product_id = p.id
LEFT JOIN users s ON r.seller_id = s.id
WHERE r.buyer_id = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListReviewsBySeller :many
SELECT r.*, 
       p.title as product_title, p.images as product_images,
       b.first_name || ' ' || b.last_name as buyer_name
FROM reviews r
LEFT JOIN products p ON r.product_id = p.id
LEFT JOIN users b ON r.buyer_id = b.id
WHERE r.seller_id = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateReview :one
UPDATE reviews 
SET rating = $2, comment = $3, images = $4, is_verified = $5
WHERE id = $1
RETURNING *;

-- name: DeleteReview :exec
DELETE FROM reviews WHERE id = $1;

-- name: GetProductRatingStats :one
SELECT 
    COUNT(*) as total_reviews,
    AVG(rating) as average_rating,
    COUNT(CASE WHEN rating = 5 THEN 1 END) as five_star,
    COUNT(CASE WHEN rating = 4 THEN 1 END) as four_star,
    COUNT(CASE WHEN rating = 3 THEN 1 END) as three_star,
    COUNT(CASE WHEN rating = 2 THEN 1 END) as two_star,
    COUNT(CASE WHEN rating = 1 THEN 1 END) as one_star
FROM reviews 
WHERE product_id = $1;

-- name: GetSellerRatingStats :one
SELECT 
    COUNT(*) as total_reviews,
    AVG(rating) as average_rating,
    COUNT(CASE WHEN rating = 5 THEN 1 END) as five_star,
    COUNT(CASE WHEN rating = 4 THEN 1 END) as four_star,
    COUNT(CASE WHEN rating = 3 THEN 1 END) as three_star,
    COUNT(CASE WHEN rating = 2 THEN 1 END) as two_star,
    COUNT(CASE WHEN rating = 1 THEN 1 END) as one_star
FROM reviews 
WHERE seller_id = $1;

-- name: CheckUserCanReview :one
SELECT EXISTS(
    SELECT 1 FROM orders o
    WHERE o.buyer_id = $1 
      AND o.product_id = $2 
      AND o.status = 'delivered'
      AND NOT EXISTS(
          SELECT 1 FROM reviews r 
          WHERE r.buyer_id = $1 AND r.product_id = $2
      )
) as can_review;