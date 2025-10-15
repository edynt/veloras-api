-- +goose Up
-- +goose StatementBegin
-- Insert sample users with unique usernames
INSERT INTO users (email, username, password, phone_number, first_name, last_name, is_verified, status, language) VALUES
('john.doe@example.com', 'johndoe123', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '+84901234567', 'John', 'Doe', true, 1, 'en'),
('jane.smith@example.com', 'janesmith456', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '+84901234568', 'Jane', 'Smith', true, 1, 'en'),
('mike.wilson@example.com', 'mikewilson789', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '+84901234569', 'Mike', 'Wilson', true, 1, 'en'),
('sarah.johnson@example.com', 'sarahjohnson101', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '+84901234570', 'Sarah', 'Johnson', true, 1, 'en'),
('david.brown@example.com', 'davidbrown202', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '+84901234571', 'David', 'Brown', true, 1, 'en'),
('seller1@example.com', 'seller1tech', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '+84901234577', 'Alex', 'Thompson', true, 1, 'en'),
('seller2@example.com', 'seller2fashion', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '+84901234578', 'Maria', 'Rodriguez', true, 1, 'en'),
('seller3@example.com', 'seller3home', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '+84901234579', 'Kevin', 'Lee', true, 1, 'en'),
('seller4@example.com', 'seller4sports', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '+84901234580', 'Anna', 'White', true, 1, 'en'),
('seller5@example.com', 'seller5books', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '+84901234581', 'Tom', 'Harris', true, 1, 'en')
ON CONFLICT (email) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete sample users
DELETE FROM users WHERE email IN (
    'john.doe@example.com', 'jane.smith@example.com', 'mike.wilson@example.com', 'sarah.johnson@example.com', 'david.brown@example.com',
    'seller1@example.com', 'seller2@example.com', 'seller3@example.com', 'seller4@example.com', 'seller5@example.com'
);
-- +goose StatementEnd
