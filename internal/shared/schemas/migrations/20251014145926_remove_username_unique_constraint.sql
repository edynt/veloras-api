-- +goose Up
-- +goose StatementBegin
-- Remove unique constraint on username column since we no longer use username
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Re-add unique constraint on username column
ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username);
-- +goose StatementEnd
