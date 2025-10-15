-- +goose Up
-- +goose StatementBegin
-- Insert coupons data
INSERT INTO coupons (code, description, discount_type, discount_value, min_order_amount, max_discount_amount, usage_limit, used_count, is_active, valid_from, valid_until)
VALUES 
-- Percentage discounts
('WELCOME10', 'Welcome discount for new users', 'percentage', 10.00, 500000, 100000, 1000, 0, true, NOW(), NOW() + INTERVAL '30 days'),
('SAVE20', 'Save 20% on your order', 'percentage', 20.00, 1000000, 200000, 500, 0, true, NOW(), NOW() + INTERVAL '15 days'),
('SUMMER15', 'Summer sale discount', 'percentage', 15.00, 800000, 150000, 200, 0, true, NOW(), NOW() + INTERVAL '7 days'),
('FLASH25', 'Flash sale - limited time offer', 'percentage', 25.00, 2000000, 500000, 100, 0, true, NOW(), NOW() + INTERVAL '3 days'),
('STUDENT10', 'Student discount', 'percentage', 10.00, 300000, 50000, 2000, 0, true, NOW(), NOW() + INTERVAL '60 days'),

-- Fixed amount discounts
('SAVE50K', 'Save 50,000 VND on your order', 'fixed', 50000, 500000, 50000, 1000, 0, true, NOW(), NOW() + INTERVAL '30 days'),
('SAVE100K', 'Save 100,000 VND on your order', 'fixed', 100000, 1000000, 100000, 500, 0, true, NOW(), NOW() + INTERVAL '15 days'),
('SAVE200K', 'Save 200,000 VND on your order', 'fixed', 200000, 2000000, 200000, 200, 0, true, NOW(), NOW() + INTERVAL '7 days'),
('FREESHIP', 'Free shipping on orders over 500K', 'fixed', 30000, 500000, 30000, 1000, 0, true, NOW(), NOW() + INTERVAL '30 days'),
('NEWUSER', 'New user special offer', 'fixed', 75000, 750000, 75000, 500, 0, true, NOW(), NOW() + INTERVAL '14 days'),

-- Special occasion coupons
('BLACKFRIDAY', 'Black Friday mega sale', 'percentage', 30.00, 3000000, 1000000, 50, 0, true, NOW(), NOW() + INTERVAL '2 days'),
('CYBERMONDAY', 'Cyber Monday special', 'percentage', 25.00, 2000000, 500000, 100, 0, true, NOW(), NOW() + INTERVAL '1 day'),
('CHRISTMAS', 'Christmas holiday discount', 'percentage', 20.00, 1500000, 300000, 200, 0, true, NOW(), NOW() + INTERVAL '10 days'),
('NEWYEAR', 'New Year celebration offer', 'percentage', 15.00, 1000000, 200000, 300, 0, true, NOW(), NOW() + INTERVAL '5 days'),
('VALENTINE', 'Valentine''s Day special', 'percentage', 18.00, 1200000, 250000, 150, 0, true, NOW(), NOW() + INTERVAL '8 days'),

-- Expired coupons (for testing)
('EXPIRED10', 'This coupon has expired', 'percentage', 10.00, 500000, 100000, 100, 0, false, NOW() - INTERVAL '10 days', NOW() - INTERVAL '1 day'),
('USEDUP', 'This coupon has been fully used', 'percentage', 15.00, 800000, 150000, 10, 10, true, NOW(), NOW() + INTERVAL '30 days');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete all coupons
DELETE FROM coupons;
-- +goose StatementEnd
