-- +goose Up
-- +goose StatementBegin
-- Insert wishlist items data
INSERT INTO wishlist_items (user_id, product_id)
SELECT 
    u.id,
    p.id
FROM users u
CROSS JOIN products p
WHERE u.email IN ('john.doe@example.com', 'jane.smith@example.com', 'mike.wilson@example.com', 'sarah.johnson@example.com', 'david.brown@example.com')
  AND p.seller_id != u.id -- Users can't add their own products to wishlist
  AND p.status = 'active' -- Only active products
LIMIT 20; -- Limit to prevent too many wishlist items
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete all wishlist items
DELETE FROM wishlist_items;
-- +goose StatementEnd
