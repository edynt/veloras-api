package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
)

type RateLimitRepository interface {
	CreateRateLimitRecord(ctx context.Context, record *entity.RateLimit) (string, error)
	GetRateLimitRecord(ctx context.Context, userId string, action string, ipAddress string) (*entity.RateLimit, error)
	UpdateRateLimitAttempts(ctx context.Context, id string, attempts int, windowStart int64) error
	DeleteExpiredRateLimits(ctx context.Context, currentTime int64) error
	GetRateLimitByIP(ctx context.Context, action string, ipAddress string) (*entity.RateLimit, error)
	ResetRateLimit(ctx context.Context, id string, windowStart int64) error
}
