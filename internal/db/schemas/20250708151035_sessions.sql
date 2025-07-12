-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions (
    session_id SERIAL PRIMARY KEY,
    session_user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    session_refresh_token TEXT NOT NULL,
    session_user_agent TEXT,
    session_client_ip TEXT,
    session_expires_at TIMESTAMPTZ NOT NULL,
    session_created_at TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd
