-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_info (
    uif_id SERIAL PRIMARY KEY,
    uif_user_id UUID NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
    
    uif_full_name TEXT,
    uif_avatar_url TEXT,
    uif_gender TEXT CHECK (uif_gender IN ('male', 'female', 'other')),
    uif_date_of_birth DATE,
    uif_phone_number TEXT,
    uif_address TEXT,
    uif_country TEXT, 
    uif_created_at TIMESTAMPTZ DEFAULT now(),
    uif_updated_at TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
