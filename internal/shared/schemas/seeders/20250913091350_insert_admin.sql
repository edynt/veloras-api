-- +goose Up
-- +goose StatementBegin
INSERT INTO users (
    email,
    password,
    is_verified,
    phone_number,
    first_name,
    last_name,
    status,
    language
) VALUES (
    'admin@verloras.com',
    '$2a$10$VF8ohSg8E48qd.jVM5U0SOes3aJGrHMXJPppvVYRh2XSGjZjOdMO6', -- hash của Abc@1234
    TRUE,
    '0000000000',
    'Tri',
    'Nguyen',
    1,
    'en'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE email = 'admin@verloras.com';
-- +goose StatementEnd
