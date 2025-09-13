-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    created_at BIGINT DEFAULT extract(epoch from now())
);

INSERT INTO roles (name, description)
VALUES 
    ('admin', 'Administrator role with full permissions'),
    ('user', 'Standard user role with limited permissions')
ON CONFLICT (name) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
