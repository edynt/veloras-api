-- +goose Up
-- +goose StatementBegin
-- USERS
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_emaill ON users(email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_is_verified ON users(is_verified);

-- ROLES
CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_name ON roles(name);

-- PERMISSIONS
CREATE UNIQUE INDEX IF NOT EXISTS idx_permissions_name ON permissions(name);

-- USER_ROLES
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);

-- ROLE_PERMISSIONS
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
