-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_roles (
    ur_user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    ur_role_id UUID REFERENCES roles(role_id) ON DELETE CASCADE,
    PRIMARY KEY (ur_user_id, ur_role_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_roles;
-- +goose StatementEnd
