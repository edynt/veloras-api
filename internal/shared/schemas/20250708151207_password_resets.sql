-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS password_resets (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    reset_token TEXT NOT NULL,
    expires_at BIGINT NOT NULL,
    created_at BIGINT DEFAULT extract(epoch from now())
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS password_resets;
-- +goose StatementEnd
