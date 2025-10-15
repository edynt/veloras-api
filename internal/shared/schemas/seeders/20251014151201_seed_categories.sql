-- +goose Up
-- +goose StatementBegin
-- Insert categories data
INSERT INTO categories (name, slug, description, image, parent_id, product_count, is_active) VALUES
-- Main categories
('Electronics', 'electronics', 'Electronic devices and gadgets', 'https://images.unsplash.com/photo-1498049794561-7780e7231661?w=500', NULL, 0, true),
('Fashion', 'fashion', 'Clothing, shoes, and accessories', 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=500', NULL, 0, true),
('Home & Garden', 'home-garden', 'Home improvement and garden supplies', 'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=500', NULL, 0, true),
('Sports & Outdoors', 'sports-outdoors', 'Sports equipment and outdoor gear', 'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=500', NULL, 0, true),
('Books & Media', 'books-media', 'Books, movies, and music', 'https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=500', NULL, 0, true),
('Health & Beauty', 'health-beauty', 'Health and beauty products', 'https://images.unsplash.com/photo-1596462502278-27bfdc403348?w=500', NULL, 0, true),
('Automotive', 'automotive', 'Car parts and accessories', 'https://images.unsplash.com/photo-1492144534655-ae79c964c9d7?w=500', NULL, 0, true),
('Toys & Games', 'toys-games', 'Toys and gaming equipment', 'https://images.unsplash.com/photo-1606092195730-5d7b9af1efc5?w=500', NULL, 0, true)
ON CONFLICT (slug) DO NOTHING;

-- Get category IDs for subcategories
WITH category_ids AS (
    SELECT id, name FROM categories WHERE parent_id IS NULL
)
INSERT INTO categories (name, slug, description, image, parent_id, product_count, is_active)
SELECT 
    subcat.name,
    subcat.slug,
    subcat.description,
    subcat.image,
    ci.id,
    0,
    true
FROM category_ids ci
CROSS JOIN (
    VALUES 
    -- Electronics subcategories
    ('Smartphones', 'smartphones', 'Mobile phones and accessories', 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=500'),
    ('Laptops', 'laptops', 'Laptops and notebooks', 'https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=500'),
    ('Tablets', 'tablets', 'Tablets and e-readers', 'https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?w=500'),
    ('Audio', 'audio', 'Headphones, speakers, and audio equipment', 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=500'),
    ('Cameras', 'cameras', 'Digital cameras and photography equipment', 'https://images.unsplash.com/photo-1502920917128-1aa500764cbd?w=500'),
    
    -- Fashion subcategories
    ('Men''s Clothing', 'mens-clothing', 'Men''s apparel and accessories', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=500'),
    ('Women''s Clothing', 'womens-clothing', 'Women''s apparel and accessories', 'https://images.unsplash.com/photo-1434389677669-e08b4cac3105?w=500'),
    ('Shoes', 'shoes', 'Footwear for all occasions', 'https://images.unsplash.com/photo-1549298916-b41d501d3772?w=500'),
    ('Accessories', 'accessories', 'Bags, jewelry, and fashion accessories', 'https://images.unsplash.com/photo-1515562141207-7a88fb7ce338?w=500'),
    
    -- Home & Garden subcategories
    ('Furniture', 'furniture', 'Home and office furniture', 'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=500'),
    ('Kitchen & Dining', 'kitchen-dining', 'Kitchen appliances and dining sets', 'https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=500'),
    ('Garden Tools', 'garden-tools', 'Gardening equipment and tools', 'https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=500'),
    ('Home Decor', 'home-decor', 'Decorative items for home', 'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=500')
) AS subcat(name, slug, description, image)
WHERE ci.name IN ('Electronics', 'Fashion', 'Home & Garden')
ON CONFLICT (slug) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete all categories
DELETE FROM categories;
-- +goose StatementEnd
