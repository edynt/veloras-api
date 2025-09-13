-- +goose Up
-- +goose StatementBegin
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
JOIN roles r ON r.name = 'admin'
WHERE u.email = 'admin@verloras.com'
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM user_roles
WHERE user_id = (SELECT id FROM users WHERE email = 'admin@verloras.com')
  AND role_id = (SELECT id FROM roles WHERE name = 'admin');
-- +goose StatementEnd
