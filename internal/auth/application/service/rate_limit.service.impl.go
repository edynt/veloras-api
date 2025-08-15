package service

import (
	"context"
	"fmt"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	rateLimitRepo "github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
)

type rateLimitService struct {
	rateLimitRepo rateLimitRepo.RateLimitRepository
	config        appDto.RateLimitConfig
}

// CheckRateLimit implements RateLimitService.
func (r *rateLimitService) CheckRateLimit(ctx context.Context, req appDto.CheckRateLimitRequest) (bool, error) {
	// Check if rate limit is exceeded
	isLimited, err := r.IsRateLimited(ctx, req.UserID, req.Action, req.IPAddress)
	if err != nil {
		return false, err
	}

	if isLimited {
		return false, nil
	}

	return true, nil
}

// CreateRateLimitRecord implements RateLimitService.
func (r *rateLimitService) CreateRateLimitRecord(ctx context.Context, req appDto.CheckRateLimitRequest) (string, error) {
	now := utils.GetNowUnix()
	windowStart := now - (now % r.config.WindowSize)

	rateLimit := &entity.RateLimit{
		UserID:      req.UserID,
		Action:      req.Action,
		IPAddress:   req.IPAddress,
		Attempts:    1,
		WindowStart: windowStart,
	}

	recordId, err := r.rateLimitRepo.CreateRateLimitRecord(ctx, rateLimit)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.FailedToCreateRateLimitRecord, err)
	}

	return recordId, nil
}

// UpdateRateLimitAttempts implements RateLimitService.
func (r *rateLimitService) UpdateRateLimitAttempts(ctx context.Context, id string, attempts int, windowStart int64) error {
	err := r.rateLimitRepo.UpdateRateLimitAttempts(ctx, id, attempts, windowStart)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdateRateLimitAttempts, err)
	}

	return nil
}

// ResetRateLimit implements RateLimitService.
func (r *rateLimitService) ResetRateLimit(ctx context.Context, id string, windowStart int64) error {
	err := r.rateLimitRepo.ResetRateLimit(ctx, id, windowStart)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToResetRateLimit, err)
	}

	return nil
}

// DeleteExpiredRateLimits implements RateLimitService.
func (r *rateLimitService) DeleteExpiredRateLimits(ctx context.Context, currentTime int64) error {
	err := r.rateLimitRepo.DeleteExpiredRateLimits(ctx, currentTime)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeleteExpiredRateLimits, err)
	}

	return nil
}

// IsRateLimited implements RateLimitService.
func (r *rateLimitService) IsRateLimited(ctx context.Context, userId, action, ipAddress string) (bool, error) {
	now := utils.GetNowUnix()
	windowStart := now - (now % r.config.WindowSize)

	// Try to get existing rate limit record
	record, err := r.rateLimitRepo.GetRateLimitRecord(ctx, userId, action, ipAddress)
	if err != nil {
		// If no record exists, create one
		_, err = r.CreateRateLimitRecord(ctx, appDto.CheckRateLimitRequest{
			UserID:    userId,
			Action:    action,
			IPAddress: ipAddress,
		})
		if err != nil {
			return false, err
		}
		return false, nil
	}

	// Check if we're in a new time window
	if record.WindowStart < windowStart {
		// Reset for new window
		err = r.ResetRateLimit(ctx, record.ID, windowStart)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	// Check if attempts exceed limit
	if record.Attempts >= r.config.MaxAttempts {
		return true, nil
	}

	// Increment attempts
	err = r.UpdateRateLimitAttempts(ctx, record.ID, record.Attempts+1, record.WindowStart)
	if err != nil {
		return false, err
	}

	return false, nil
}

// SetRateLimitConfig sets the rate limiting configuration
func (r *rateLimitService) SetRateLimitConfig(config appDto.RateLimitConfig) {
	r.config = config
}

// GetRateLimitConfig returns the current rate limiting configuration
func (r *rateLimitService) GetRateLimitConfig() appDto.RateLimitConfig {
	return r.config
}

func NewRateLimitService(rateLimitRepo rateLimitRepo.RateLimitRepository, config appDto.RateLimitConfig) RateLimitService {
	return &rateLimitService{
		rateLimitRepo: rateLimitRepo,
		config:        config,
	}
}
