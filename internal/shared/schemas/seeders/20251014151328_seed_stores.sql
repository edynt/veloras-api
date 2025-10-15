-- +goose Up
-- +goose StatementBegin
-- Insert stores for sellers
INSERT INTO stores (user_id, name, description, logo, banner, address, phone, email, is_verified, rating, review_count)
SELECT 
    u.id,
    store_data.name,
    store_data.description,
    store_data.logo,
    store_data.banner,
    store_data.address,
    store_data.phone,
    store_data.email,
    store_data.is_verified,
    store_data.rating,
    store_data.review_count
FROM users u
CROSS JOIN (
    VALUES 
    ('TechStore Pro', 'Premium electronics and gadgets store', 'https://images.unsplash.com/photo-1560472354-b33ff0c44a43?w=200', 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=800', '123 Tech Street, District 1, Ho Chi Minh City', '+84901234577', 'seller1@example.com', true, 4.8, 156),
    ('Fashion Hub', 'Trendy clothing and accessories', 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=200', 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=800', '456 Fashion Avenue, District 3, Ho Chi Minh City', '+84901234578', 'seller2@example.com', true, 4.6, 89),
    ('Home & Garden Plus', 'Everything for your home and garden', 'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=200', 'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=800', '789 Garden Road, District 7, Ho Chi Minh City', '+84901234579', 'seller3@example.com', true, 4.7, 123),
    ('Sports Central', 'Sports equipment and outdoor gear', 'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=200', 'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=800', '321 Sports Plaza, District 2, Ho Chi Minh City', '+84901234580', 'seller4@example.com', true, 4.5, 67),
    ('Book World', 'Books, media, and educational materials', 'https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=200', 'https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=800', '654 Book Lane, District 5, Ho Chi Minh City', '+84901234581', 'seller5@example.com', true, 4.9, 234)
) AS store_data(name, description, logo, banner, address, phone, email, is_verified, rating, review_count)
WHERE u.email = store_data.email;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete all stores
DELETE FROM stores;
-- +goose StatementEnd
