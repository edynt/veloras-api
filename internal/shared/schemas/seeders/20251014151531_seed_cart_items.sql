-- +goose Up
-- +goose StatementBegin
-- Insert cart items data
INSERT INTO cart_items (user_id, product_id, quantity)
SELECT 
    u.id,
    p.id,
    cart_data.quantity
FROM users u
CROSS JOIN products p
CROSS JOIN (
    VALUES 
    -- Sample cart items
    (1), (2), (1), (3), (1), (2), (1), (1), (2), (1), (1), (3), (1), (2), (1)
) AS cart_data(quantity)
WHERE u.email IN ('john.doe@example.com', 'jane.smith@example.com', 'mike.wilson@example.com', 'sarah.johnson@example.com', 'david.brown@example.com')
  AND p.seller_id != u.id -- Users can't add their own products to cart
LIMIT 25; -- Limit to prevent too many cart items
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete all cart items
DELETE FROM cart_items;
-- +goose StatementEnd
