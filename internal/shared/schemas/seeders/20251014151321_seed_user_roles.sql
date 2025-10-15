-- +goose Up
-- +goose StatementBegin
-- Assign roles to users
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE 
    -- Admin users get admin role
    (u.email IN ('admin@example.com', 'superadmin@example.com') AND r.name = 'admin')
    OR
    -- Sellers get seller role
    (u.email LIKE 'seller%@example.com' AND r.name = 'seller')
    OR
    -- All other users get user role
    (u.email NOT IN ('admin@example.com', 'superadmin@example.com') 
     AND u.email NOT LIKE 'seller%@example.com' 
     AND r.name = 'user')
ON CONFLICT (user_id, role_id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove all user role assignments
DELETE FROM user_roles;
-- +goose StatementEnd
