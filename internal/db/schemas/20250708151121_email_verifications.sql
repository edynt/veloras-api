-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS email_verifications (
    ev_id SERIAL PRIMARY KEY,
    ev_user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    ev_code TEXT NOT NULL,
    ev_expires_at TIMESTAMPTZ NOT NULL,
    ev_created_at TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS email_verifications;
-- +goose StatementEnd
