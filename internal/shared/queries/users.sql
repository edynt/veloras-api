-- name: CreateUser :one
INSERT INTO users (email, username, password, phone_number, first_name, last_name)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, email, status;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: VerifyUser :exec
UPDATE users SET is_verified = TRUE WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: GetUserEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: GetUsernameExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);

-- name: UpdateUserStatus :one
UPDATE users SET status = $1 WHERE id = $2 RETURNING id, email, status;

-- name: ActiveUser :one
UPDATE users SET is_verified = TRUE WHERE id = $1 RETURNING id, email, status;

-- name: DeleteVerificationCode :exec
DELETE FROM email_verifications WHERE user_id = $1;

-- name: UpdateUserPassword :exec
UPDATE users SET password = $1 WHERE id = $2;

-- name: UpdateUser2FASecret :exec
UPDATE users SET two_fa_secret = $1 WHERE id = $2;

-- name: UpdateUser2FAEnabled :exec
UPDATE users SET two_fa_enabled = $1 WHERE id = $2;

-- name: IncrementFailedLoginAttempts :exec
UPDATE users SET failed_login_attempts = failed_login_attempts + 1 WHERE id = $1;

-- name: ResetFailedLoginAttempts :exec
UPDATE users SET failed_login_attempts = 0 WHERE id = $1;

-- name: LockUserAccount :exec
UPDATE users SET status = $1, locked_until = $2 WHERE id = $3;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserProfile :exec
UPDATE users SET first_name = $1, last_name = $2, phone_number = $3, language = $4 WHERE id = $5;

-- name: GetUsersByOrganization :many
SELECT * FROM users WHERE organization_id = $1;

-- name: GetUsersByRole :many
SELECT u.* FROM users u
JOIN user_roles ur ON u.id = ur.user_id
WHERE ur.role_id = $1;

-- name: UpdateLastLogin :exec
UPDATE users SET last_login_at = $1 WHERE id = $2;