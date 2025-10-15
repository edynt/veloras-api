-- +goose Up
-- +goose StatementBegin
-- Remove NOT NULL constraint on username column since we no longer use username
ALTER TABLE users ALTER COLUMN username DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Re-add NOT NULL constraint on username column
ALTER TABLE users ALTER COLUMN username SET NOT NULL;
-- +goose StatementEnd
