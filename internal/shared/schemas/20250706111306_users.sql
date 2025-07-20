-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    phone_number TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    status INT DEFAULT 0,
    language TEXT DEFAULT 'en',
    created_at BIGINT DEFAULT extract(epoch from now()),
    updated_at BIGINT DEFAULT extract(epoch from now())
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd