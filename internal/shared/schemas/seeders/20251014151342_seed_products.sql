-- +goose Up
-- +goose StatementBegin
-- Insert products data
INSERT INTO products (title, description, price, original_price, category_id, images, condition, tags, location, stock, seller_id, status, badges, rating, review_count, view_count, is_featured, is_promoted)
SELECT 
    product_data.title,
    product_data.description,
    product_data.price,
    product_data.original_price,
    c.id,
    product_data.images,
    product_data.condition,
    product_data.tags,
    product_data.location,
    product_data.stock,
    u.id,
    product_data.status,
    product_data.badges,
    product_data.rating,
    product_data.review_count,
    product_data.view_count,
    product_data.is_featured,
    product_data.is_promoted
FROM users u
CROSS JOIN categories c
CROSS JOIN (
    VALUES 
    -- Electronics products
    ('iPhone 15 Pro Max 256GB', 'Latest iPhone with titanium design and advanced camera system', 25000000, 28000000, 'Electronics', ARRAY['https://images.unsplash.com/photo-1592750475338-74b7b21085ab?w=500', 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=500'], 'new', ARRAY['smartphone', 'apple', 'premium'], 'Ho Chi Minh City', 5, 'active', ARRAY['featured', 'trending'], 4.8, 23, 156, true, true),
    ('MacBook Pro M3 14-inch', 'Powerful laptop with M3 chip for professionals', 35000000, 38000000, 'Electronics', ARRAY['https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=500', 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=500'], 'new', ARRAY['laptop', 'apple', 'professional'], 'Ho Chi Minh City', 3, 'active', ARRAY['featured'], 4.9, 15, 89, true, false),
    ('Samsung Galaxy S24 Ultra', 'Premium Android smartphone with S Pen', 22000000, 25000000, 'Electronics', ARRAY['https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=500'], 'new', ARRAY['smartphone', 'samsung', 'android'], 'Ho Chi Minh City', 8, 'active', ARRAY['trending'], 4.7, 31, 234, false, true),
    ('AirPods Pro 2nd Gen', 'Wireless earbuds with active noise cancellation', 5500000, 6000000, 'Electronics', ARRAY['https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=500'], 'new', ARRAY['earbuds', 'apple', 'wireless'], 'Ho Chi Minh City', 12, 'active', ARRAY['popular'], 4.6, 45, 178, false, false),
    ('Sony WH-1000XM5', 'Premium noise-canceling headphones', 8000000, 8500000, 'Electronics', ARRAY['https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=500'], 'new', ARRAY['headphones', 'sony', 'noise-canceling'], 'Ho Chi Minh City', 6, 'active', ARRAY['featured'], 4.8, 28, 145, true, false),
    
    -- Fashion products
    ('Nike Air Max 270', 'Comfortable running shoes with Air Max technology', 3200000, 3500000, 'Fashion', ARRAY['https://images.unsplash.com/photo-1549298916-b41d501d3772?w=500'], 'new', ARRAY['shoes', 'nike', 'running'], 'Ho Chi Minh City', 15, 'active', ARRAY['trending'], 4.5, 67, 312, false, true),
    ('Adidas Ultraboost 22', 'High-performance running shoes', 4500000, 4800000, 'Fashion', ARRAY['https://images.unsplash.com/photo-1549298916-b41d501d3772?w=500'], 'new', ARRAY['shoes', 'adidas', 'running'], 'Ho Chi Minh City', 10, 'active', ARRAY['featured'], 4.7, 42, 189, true, false),
    ('Levi''s 501 Original Jeans', 'Classic straight-fit jeans', 1800000, 2000000, 'Fashion', ARRAY['https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=500'], 'new', ARRAY['jeans', 'levis', 'denim'], 'Ho Chi Minh City', 25, 'active', ARRAY['classic'], 4.4, 89, 267, false, false),
    ('Zara Blazer', 'Elegant women''s blazer for office wear', 2500000, 2800000, 'Fashion', ARRAY['https://images.unsplash.com/photo-1434389677669-e08b4cac3105?w=500'], 'new', ARRAY['blazer', 'zara', 'office'], 'Ho Chi Minh City', 8, 'active', ARRAY['trending'], 4.3, 34, 156, false, true),
    ('Gucci Handbag', 'Luxury leather handbag', 15000000, 18000000, 'Fashion', ARRAY['https://images.unsplash.com/photo-1515562141207-7a88fb7ce338?w=500'], 'like_new', ARRAY['handbag', 'gucci', 'luxury'], 'Ho Chi Minh City', 2, 'active', ARRAY['luxury'], 4.9, 12, 78, true, false),
    
    -- Home & Garden products
    ('IKEA Dining Table', 'Modern wooden dining table for 6 people', 4500000, 5000000, 'Home & Garden', ARRAY['https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=500'], 'new', ARRAY['table', 'ikea', 'dining'], 'Ho Chi Minh City', 4, 'active', ARRAY['popular'], 4.6, 23, 134, false, false),
    ('Philips Air Fryer', 'Healthy cooking with air frying technology', 2800000, 3200000, 'Home & Garden', ARRAY['https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=500'], 'new', ARRAY['air-fryer', 'philips', 'kitchen'], 'Ho Chi Minh City', 7, 'active', ARRAY['trending'], 4.5, 56, 289, false, true),
    ('Garden Tool Set', 'Complete set of gardening tools', 1200000, 1500000, 'Home & Garden', ARRAY['https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=500'], 'new', ARRAY['tools', 'garden', 'set'], 'Ho Chi Minh City', 12, 'active', ARRAY['bundle'], 4.4, 18, 67, false, false),
    ('Decorative Plant Pot', 'Ceramic plant pot with modern design', 350000, 400000, 'Home & Garden', ARRAY['https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=500'], 'new', ARRAY['pot', 'ceramic', 'decorative'], 'Ho Chi Minh City', 20, 'active', ARRAY['decorative'], 4.2, 9, 45, false, false),
    
    -- Sports & Outdoors products
    ('Yoga Mat Premium', 'Non-slip yoga mat for all exercises', 800000, 950000, 'Sports & Outdoors', ARRAY['https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=500'], 'new', ARRAY['yoga', 'mat', 'exercise'], 'Ho Chi Minh City', 18, 'active', ARRAY['fitness'], 4.6, 34, 123, false, false),
    ('Tennis Racket', 'Professional tennis racket for intermediate players', 2200000, 2500000, 'Sports & Outdoors', ARRAY['https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=500'], 'new', ARRAY['tennis', 'racket', 'sports'], 'Ho Chi Minh City', 6, 'active', ARRAY['sports'], 4.5, 21, 89, false, false),
    ('Camping Tent 4-Person', 'Waterproof camping tent for outdoor adventures', 3500000, 4000000, 'Sports & Outdoors', ARRAY['https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=500'], 'new', ARRAY['tent', 'camping', 'outdoor'], 'Ho Chi Minh City', 3, 'active', ARRAY['adventure'], 4.7, 15, 67, true, false),
    
    -- Books & Media products
    ('The Psychology of Money', 'Best-selling book about money and investing', 250000, 300000, 'Books & Media', ARRAY['https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=500'], 'new', ARRAY['book', 'psychology', 'money'], 'Ho Chi Minh City', 25, 'active', ARRAY['bestseller'], 4.8, 89, 456, true, true),
    ('Atomic Habits', 'Transform your life with small changes', 280000, 320000, 'Books & Media', ARRAY['https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=500'], 'new', ARRAY['book', 'habits', 'self-help'], 'Ho Chi Minh City', 30, 'active', ARRAY['popular'], 4.9, 156, 789, true, false),
    ('Bluetooth Speaker', 'Portable wireless speaker with great sound', 1200000, 1400000, 'Books & Media', ARRAY['https://images.unsplash.com/photo-1608043152269-423dbba4e7e1?w=500'], 'new', ARRAY['speaker', 'bluetooth', 'portable'], 'Ho Chi Minh City', 15, 'active', ARRAY['audio'], 4.4, 45, 234, false, false)
) AS product_data(title, description, price, original_price, category_name, images, condition, tags, location, stock, status, badges, rating, review_count, view_count, is_featured, is_promoted)
WHERE u.email LIKE 'seller%@example.com' 
  AND c.name = product_data.category_name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete all products
DELETE FROM products;
-- +goose StatementEnd
