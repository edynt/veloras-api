-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_info (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    
    full_name TEXT,
    avatar_url TEXT,
    gender TEXT CHECK (gender IN ('male', 'female', 'other')),
    date_of_birth DATE,
    phone_number TEXT,
    address TEXT,
    country TEXT, 
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_info;
-- +goose StatementEnd
