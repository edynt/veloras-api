-- name: CreateRateLimitRecord :one
INSERT INTO rate_limits (user_id, action, ip_address, attempts, window_start)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRateLimitRecord :one
SELECT * FROM rate_limits WHERE user_id = $1 AND action = $2 AND ip_address = $3;

-- name: UpdateRateLimitAttempts :exec
UPDATE rate_limits SET attempts = $1, window_start = $2 WHERE id = $3;

-- name: DeleteExpiredRateLimits :exec
DELETE FROM rate_limits WHERE window_start < $1;

-- name: GetRateLimitByIP :one
SELECT * FROM rate_limits WHERE ip_address = $1 AND action = $2;

-- name: ResetRateLimit :exec
UPDATE rate_limits SET attempts = 0, window_start = $1 WHERE id = $2;
