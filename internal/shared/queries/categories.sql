-- name: CreateCategory :one
INSERT INTO categories (
    name, slug, description, image, parent_id, product_count, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories WHERE id = $1;

-- name: GetCategoryBySlug :one
SELECT * FROM categories WHERE slug = $1;

-- name: ListCategories :many
SELECT * FROM categories 
WHERE is_active = true 
ORDER BY name;

-- name: ListCategoriesWithPagination :many
SELECT * FROM categories 
WHERE is_active = true 
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: ListSubCategories :many
SELECT * FROM categories 
WHERE parent_id = $1 AND is_active = true 
ORDER BY name;

-- name: UpdateCategory :one
UPDATE categories 
SET name = $2, slug = $3, description = $4, image = $5, parent_id = $6, 
    product_count = $7, is_active = $8, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: UpdateCategoryProductCount :exec
UPDATE categories 
SET product_count = product_count + $2, updated_at = NOW()
WHERE id = $1;

-- name: GetCategoryStats :one
SELECT 
    COUNT(*) as total_categories,
    COUNT(CASE WHEN parent_id IS NULL THEN 1 END) as parent_categories,
    COUNT(CASE WHEN parent_id IS NOT NULL THEN 1 END) as sub_categories,
    SUM(product_count) as total_products
FROM categories 
WHERE is_active = true;