-- name: CreateCoupon :one
INSERT INTO coupons (
    code, description, discount_type, discount_value, min_order_amount,
    max_discount_amount, usage_limit, is_active, valid_from, valid_until
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetCoupon :one
SELECT * FROM coupons WHERE id = $1;

-- name: GetCouponByCode :one
SELECT * FROM coupons WHERE code = $1;

-- name: ListActiveCoupons :many
SELECT * FROM coupons 
WHERE is_active = true 
  AND valid_from <= NOW() 
  AND (valid_until IS NULL OR valid_until >= NOW())
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListCoupons :many
SELECT * FROM coupons 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateCoupon :one
UPDATE coupons 
SET code = $2, description = $3, discount_type = $4, discount_value = $5,
    min_order_amount = $6, max_discount_amount = $7, usage_limit = $8,
    is_active = $9, valid_from = $10, valid_until = $11, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeactivateCoupon :one
UPDATE coupons 
SET is_active = false, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCoupon :exec
DELETE FROM coupons WHERE id = $1;

-- name: ValidateCoupon :one
SELECT c.*, 
       CASE 
           WHEN NOT c.is_active THEN 'inactive'
           WHEN c.valid_from > NOW() THEN 'not_started'
           WHEN c.valid_until IS NOT NULL AND c.valid_until < NOW() THEN 'expired'
           WHEN c.usage_limit IS NOT NULL AND c.used_count >= c.usage_limit THEN 'limit_reached'
           ELSE 'valid'
       END as validation_status
FROM coupons c
WHERE c.code = $1;

-- name: UseCoupon :one
INSERT INTO user_coupon_usage (user_id, coupon_id, order_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: IncrementCouponUsage :exec
UPDATE coupons 
SET used_count = used_count + 1, updated_at = NOW()
WHERE id = $1;

-- name: GetCouponUsageHistory :many
SELECT ucu.*, c.code, c.description, c.discount_type, c.discount_value
FROM user_coupon_usage ucu
LEFT JOIN coupons c ON ucu.coupon_id = c.id
WHERE ucu.user_id = $1
ORDER BY ucu.used_at DESC
LIMIT $2 OFFSET $3;

-- name: GetCouponStats :one
SELECT 
    COUNT(*) as total_coupons,
    COUNT(CASE WHEN is_active = true THEN 1 END) as active_coupons,
    COUNT(CASE WHEN is_active = false THEN 1 END) as inactive_coupons,
    SUM(used_count) as total_usage,
    AVG(discount_value) as average_discount
FROM coupons;