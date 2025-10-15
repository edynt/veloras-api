-- +goose Up
-- +goose StatementBegin
-- Insert a simple test user
INSERT INTO users (email, username, password, phone_number, first_name, last_name, is_verified, status, language) 
VALUES ('test@example.com', 'testuser', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '+84901234599', 'Test', 'User', true, 1, 'en')
ON CONFLICT (email) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete test user
DELETE FROM users WHERE email = 'test@example.com';
-- +goose StatementEnd
