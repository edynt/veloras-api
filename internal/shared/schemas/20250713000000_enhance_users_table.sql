-- +goose Up
-- +goose StatementBegin

-- Add new columns to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS two_fa_secret TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS two_fa_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS failed_login_attempts INT DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS locked_until BIGINT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS organization_id UUID;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_at BIGINT;

-- Create organizations table
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    parent_id UUID REFERENCES organizations(id),
    created_by UUID REFERENCES users(id),
    created_at BIGINT DEFAULT extract(epoch from now()),
    updated_at BIGINT DEFAULT extract(epoch from now())
);

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    action TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT,
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at BIGINT DEFAULT extract(epoch from now())
);

-- Create rate_limits table
CREATE TABLE IF NOT EXISTS rate_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    action TEXT NOT NULL,
    ip_address INET NOT NULL,
    attempts INT DEFAULT 1,
    window_start BIGINT NOT NULL,
    created_at BIGINT DEFAULT extract(epoch from now()),
    updated_at BIGINT DEFAULT extract(epoch from now())
);

-- Create two_fa_codes table
CREATE TABLE IF NOT EXISTS two_fa_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) NOT NULL,
    code TEXT NOT NULL,
    expires_at BIGINT NOT NULL,
    type TEXT NOT NULL, -- 'setup', 'login', 'recovery'
    created_at BIGINT DEFAULT extract(epoch from now())
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_organization_id ON users(organization_id);
CREATE INDEX IF NOT EXISTS idx_users_failed_login_attempts ON users(failed_login_attempts);
CREATE INDEX IF NOT EXISTS idx_users_locked_until ON users(locked_until);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_rate_limits_user_action ON rate_limits(user_id, action);
CREATE INDEX IF NOT EXISTS idx_rate_limits_ip_action ON rate_limits(ip_address, action);
CREATE INDEX IF NOT EXISTS idx_two_fa_codes_user_type ON two_fa_codes(user_id, type);
CREATE INDEX IF NOT EXISTS idx_two_fa_codes_expires_at ON two_fa_codes(expires_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop indexes
DROP INDEX IF EXISTS idx_users_organization_id;
DROP INDEX IF EXISTS idx_users_failed_login_attempts;
DROP INDEX IF EXISTS idx_users_locked_until;
DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP INDEX IF EXISTS idx_audit_logs_resource;
DROP INDEX IF EXISTS idx_audit_logs_created_at;
DROP INDEX IF EXISTS idx_rate_limits_user_action;
DROP INDEX IF EXISTS idx_rate_limits_ip_action;
DROP INDEX IF EXISTS idx_two_fa_codes_user_type;
DROP INDEX IF EXISTS idx_two_fa_codes_expires_at;

-- Drop tables
DROP TABLE IF EXISTS two_fa_codes;
DROP TABLE IF EXISTS rate_limits;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS organizations;

-- Remove columns from users table
ALTER TABLE users DROP COLUMN IF EXISTS two_fa_secret;
ALTER TABLE users DROP COLUMN IF EXISTS two_fa_enabled;
ALTER TABLE users DROP COLUMN IF EXISTS failed_login_attempts;
ALTER TABLE users DROP COLUMN IF EXISTS locked_until;
ALTER TABLE users DROP COLUMN IF EXISTS organization_id;
ALTER TABLE users DROP COLUMN IF EXISTS last_login_at;

-- +goose StatementEnd
