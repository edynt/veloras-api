-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_email TEXT NOT NULL UNIQUE,
    user_username TEXT NOT NULL UNIQUE,
    user_hashed_password TEXT NOT NULL,
    user_is_verified BOOLEAN DEFAULT FALSE,
    user_created_at TIMESTAMPTZ DEFAULT now(),
    user_updated_at TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd