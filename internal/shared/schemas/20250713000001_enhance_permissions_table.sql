-- +goose Up
-- +goose StatementBegin

-- Add new columns to permissions table
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS resource_type TEXT DEFAULT 'general';
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS resource_action TEXT DEFAULT 'access';
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES organizations(id);

-- Create index for resource-based permissions
CREATE INDEX IF NOT EXISTS idx_permissions_resource ON permissions(resource_type, resource_action);
CREATE INDEX IF NOT EXISTS idx_permissions_organization ON permissions(organization_id);

-- Insert default permissions for common resources
INSERT INTO permissions (name, description, resource_type, resource_action) VALUES
('create:post', 'Create new posts', 'post', 'create'),
('read:post', 'Read posts', 'post', 'read'),
('update:post', 'Update posts', 'post', 'update'),
('delete:post', 'Delete posts', 'post', 'delete'),
('create:user', 'Create new users', 'user', 'create'),
('read:user', 'Read user information', 'user', 'read'),
('update:user', 'Update user information', 'user', 'update'),
('delete:user', 'Delete users', 'user', 'delete'),
('create:role', 'Create new roles', 'role', 'create'),
('read:role', 'Read role information', 'role', 'read'),
('update:role', 'Update role information', 'role', 'update'),
('delete:role', 'Delete roles', 'role', 'delete'),
('create:permission', 'Create new permissions', 'permission', 'create'),
('read:permission', 'Read permission information', 'permission', 'read'),
('update:permission', 'Update permission information', 'permission', 'update'),
('delete:permission', 'Delete permissions', 'permission', 'delete'),
('create:organization', 'Create new organizations', 'organization', 'create'),
('read:organization', 'Read organization information', 'organization', 'read'),
('update:organization', 'Update organization information', 'organization', 'update'),
('delete:organization', 'Delete organizations', 'organization', 'delete')
ON CONFLICT (name) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop indexes
DROP INDEX IF EXISTS idx_permissions_resource;
DROP INDEX IF EXISTS idx_permissions_organization;

-- Remove new columns from permissions table
ALTER TABLE permissions DROP COLUMN IF EXISTS resource_type;
ALTER TABLE permissions DROP COLUMN IF EXISTS resource_action;
ALTER TABLE permissions DROP COLUMN IF EXISTS organization_id;

-- +goose StatementEnd
