-- +goose Up
-- +goose StatementBegin
-- USERS
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_user_email ON users(user_email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_user_username ON users(user_username);
CREATE INDEX IF NOT EXISTS idx_users_user_is_verified ON users(user_is_verified);
CREATE INDEX IF NOT EXISTS idx_users_user_created_at ON users(user_created_at);
CREATE INDEX IF NOT EXISTS idx_users_fulltext ON users USING GIN (
  to_tsvector('english', coalesce(user_username, '') || ' ' || coalesce(user_email, ''))
);

-- ROLES
CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_role_name ON roles(role_name);

-- PERMISSIONS
CREATE UNIQUE INDEX IF NOT EXISTS idx_permissions_permission_name ON permissions(permission_name);

-- USER_ROLES
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(ur_user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(ur_role_id);

-- ROLE_PERMISSIONS
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(rp_role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(rp_permission_id);


-- USER_INFO
CREATE INDEX IF NOT EXISTS idx_user_info_user_id ON user_info(uif_user_id);
CREATE INDEX IF NOT EXISTS idx_user_info_phone ON user_info(uif_phone_number);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
