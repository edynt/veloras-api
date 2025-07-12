-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS password_resets (
    ps_id SERIAL PRIMARY KEY,
    ps_user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    ps_reset_token TEXT NOT NULL,
    ps_expires_at TIMESTAMPTZ NOT NULL,
    ps_created_at TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS password_resets;
-- +goose StatementEnd
