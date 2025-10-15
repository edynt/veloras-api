-- +goose Up
-- +goose StatementBegin
-- Insert orders data
INSERT INTO orders (buyer_id, seller_id, product_id, quantity, total_amount, status, payment_method, payment_status, shipping_address_id, notes)
SELECT 
    buyer.id,
    seller.id,
    p.id,
    order_data.quantity,
    order_data.total_amount,
    order_data.status,
    order_data.payment_method,
    order_data.payment_status,
    addr.id,
    order_data.notes
FROM users buyer
CROSS JOIN users seller
CROSS JOIN products p
CROSS JOIN addresses addr
CROSS JOIN (
    VALUES 
    -- Sample orders
    (1, 2500000, 'confirmed', 'momo', 'completed', 'Please deliver in the morning'),
    (2, 5000000, 'shipped', 'zalopay', 'completed', 'Handle with care'),
    (1, 1800000, 'delivered', 'stripe', 'completed', 'Thank you!'),
    (1, 3200000, 'pending', 'bank_transfer', 'pending', 'Will pay tomorrow'),
    (2, 4500000, 'confirmed', 'momo', 'completed', 'Fast delivery please'),
    (1, 2800000, 'cancelled', 'zalopay', 'failed', 'Changed mind'),
    (3, 7500000, 'delivered', 'stripe', 'completed', 'Great product!'),
    (1, 1200000, 'shipped', 'momo', 'completed', 'Track the package'),
    (2, 3500000, 'confirmed', 'bank_transfer', 'completed', 'Office delivery'),
    (1, 800000, 'pending', 'zalopay', 'pending', 'Waiting for payment'),
    (1, 2200000, 'delivered', 'stripe', 'completed', 'Perfect condition'),
    (2, 1500000, 'shipped', 'momo', 'completed', 'Handle carefully'),
    (1, 2800000, 'confirmed', 'zalopay', 'completed', 'Weekend delivery'),
    (1, 4500000, 'delivered', 'bank_transfer', 'completed', 'Excellent service'),
    (2, 1200000, 'pending', 'stripe', 'pending', 'Will confirm soon')
) AS order_data(quantity, total_amount, status, payment_method, payment_status, notes)
WHERE buyer.email IN ('john.doe@example.com', 'jane.smith@example.com', 'mike.wilson@example.com', 'sarah.johnson@example.com', 'david.brown@example.com')
  AND seller.email LIKE 'seller%@example.com'
  AND p.seller_id = seller.id
  AND addr.user_id = buyer.id
  AND addr.is_default = true
LIMIT 50; -- Limit to prevent too many orders
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete all orders
DELETE FROM orders;
-- +goose StatementEnd
