-- +goose Up
-- +goose StatementBegin
-- Gán tất cả các quyền hiện có trong bảng permissions cho role admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON TRUE
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Xóa tất cả quyền của role admin
DELETE FROM role_permissions
WHERE role_id = (SELECT id FROM roles WHERE name = 'admin');
-- +goose StatementEnd
