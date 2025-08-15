package service

import (
	"context"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
)

type RateLimitService interface {
	CheckRateLimit(ctx context.Context, req appDto.CheckRateLimitRequest) (bool, error)
	CreateRateLimitRecord(ctx context.Context, req appDto.CheckRateLimitRequest) (string, error)
	UpdateRateLimitAttempts(ctx context.Context, id string, attempts int, windowStart int64) error
	ResetRateLimit(ctx context.Context, id string, windowStart int64) error
	DeleteExpiredRateLimits(ctx context.Context, currentTime int64) error
	IsRateLimited(ctx context.Context, userId, action, ipAddress string) (bool, error)
}
