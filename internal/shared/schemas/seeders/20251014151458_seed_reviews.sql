-- +goose Up
-- +goose StatementBegin
-- Insert reviews data
INSERT INTO reviews (product_id, buyer_id, seller_id, rating, comment, images, is_verified)
SELECT 
    p.id,
    buyer.id,
    seller.id,
    review_data.rating,
    review_data.comment,
    review_data.images,
    review_data.is_verified
FROM users buyer
CROSS JOIN users seller
CROSS JOIN products p
CROSS JOIN (
    VALUES 
    -- Sample reviews
    (5, 'Excellent product! Highly recommended.', ARRAY['https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=500'], true),
    (4, 'Good quality, fast shipping. Will buy again.', ARRAY['https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=500'], true),
    (5, 'Perfect! Exactly as described.', ARRAY['https://images.unsplash.com/photo-1549298916-b41d501d3772?w=500'], true),
    (3, 'Average product, could be better.', ARRAY['https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=500'], false),
    (4, 'Great value for money. Satisfied with purchase.', ARRAY['https://images.unsplash.com/photo-1434389677669-e08b4cac3105?w=500'], true),
    (5, 'Amazing quality! Seller was very helpful.', ARRAY['https://images.unsplash.com/photo-1515562141207-7a88fb7ce338?w=500'], true),
    (4, 'Good product, arrived on time.', ARRAY['https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=500'], true),
    (2, 'Not as expected. Poor quality.', ARRAY['https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=500'], false),
    (5, 'Outstanding! Best purchase ever.', ARRAY['https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=500'], true),
    (4, 'Very good product. Recommended!', ARRAY['https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=500'], true),
    (3, 'Okay product, nothing special.', ARRAY['https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=500'], false),
    (5, 'Fantastic! Exceeded my expectations.', ARRAY['https://images.unsplash.com/photo-1608043152269-423dbba4e7e1?w=500'], true),
    (4, 'Good product, good service.', ARRAY['https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=500'], true),
    (1, 'Terrible product. Waste of money.', ARRAY['https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=500'], false),
    (5, 'Perfect! Will definitely buy again.', ARRAY['https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=500'], true)
) AS review_data(rating, comment, images, is_verified)
WHERE buyer.email IN ('john.doe@example.com', 'jane.smith@example.com', 'mike.wilson@example.com', 'sarah.johnson@example.com', 'david.brown@example.com')
  AND seller.email LIKE 'seller%@example.com'
  AND p.seller_id = seller.id
LIMIT 30; -- Limit to prevent too many reviews
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete all reviews
DELETE FROM reviews;
-- +goose StatementEnd
