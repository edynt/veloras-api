package repository

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/internal/shared/gen"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type rateLimitRepository struct {
	db *gen.Queries
}

// CreateRateLimitRecord implements repository.RateLimitRepository.
func (r *rateLimitRepository) CreateRateLimitRecord(ctx context.Context, record *entity.RateLimit) (string, error) {
	convertUserId, err := utils.ConvertUUID(record.UserID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	ipAddr, err := netip.ParseAddr(record.IPAddress)
	if err != nil {
		return "", fmt.Errorf("invalid IP address: %w", err)
	}

	param := gen.CreateRateLimitRecordParams{
		UserID:    convertUserId,
		Action:    record.Action,
		IpAddress: ipAddr,
		Attempts:  pgtype.Int4{Int32: int32(record.Attempts), Valid: true},
		WindowStart: record.WindowStart,
	}

	result, err := r.db.CreateRateLimitRecord(ctx, param)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.FailedToCreateRateLimitRecord, err)
	}

	return result.ID.String(), nil
}

// GetRateLimitRecord implements repository.RateLimitRepository.
func (r *rateLimitRepository) GetRateLimitRecord(ctx context.Context, userId string, action string, ipAddress string) (*entity.RateLimit, error) {
	convertUserId, err := utils.ConvertUUID(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	ipAddr, err := netip.ParseAddr(ipAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid IP address: %w", err)
	}

	param := gen.GetRateLimitRecordParams{
		UserID:    convertUserId,
		Action:    action,
		IpAddress: ipAddr,
	}

	result, err := r.db.GetRateLimitRecord(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetRateLimitRecord, err)
	}

	var entityResult entity.RateLimit
	if err := utils.SafeCopy(&entityResult, &result); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// UpdateRateLimitAttempts implements repository.RateLimitRepository.
func (r *rateLimitRepository) UpdateRateLimitAttempts(ctx context.Context, id string, attempts int, windowStart int64) error {
	convertId, err := utils.ConvertUUID(id)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.RateLimitIdInvalid, err)
	}

	param := gen.UpdateRateLimitAttemptsParams{
		Attempts:    pgtype.Int4{Int32: int32(attempts), Valid: true},
		WindowStart: windowStart,
		ID:          convertId,
	}

	err = r.db.UpdateRateLimitAttempts(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdateRateLimitAttempts, err)
	}

	return nil
}

// DeleteExpiredRateLimits implements repository.RateLimitRepository.
func (r *rateLimitRepository) DeleteExpiredRateLimits(ctx context.Context, currentTime int64) error {
	err := r.db.DeleteExpiredRateLimits(ctx, currentTime)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeleteExpiredRateLimits, err)
	}

	return nil
}

// GetRateLimitByIP implements repository.RateLimitRepository.
func (r *rateLimitRepository) GetRateLimitByIP(ctx context.Context, action string, ipAddress string) (*entity.RateLimit, error) {
	ipAddr, err := netip.ParseAddr(ipAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid IP address: %w", err)
	}

	param := gen.GetRateLimitByIPParams{
		Action:    action,
		IpAddress: ipAddr,
	}

	result, err := r.db.GetRateLimitByIP(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetRateLimitRecord, err)
	}

	var entityResult entity.RateLimit
	if err := utils.SafeCopy(&entityResult, &result); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// ResetRateLimit implements repository.RateLimitRepository.
func (r *rateLimitRepository) ResetRateLimit(ctx context.Context, id string, windowStart int64) error {
	convertId, err := utils.ConvertUUID(id)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.RateLimitIdInvalid, err)
	}

	param := gen.ResetRateLimitParams{
		WindowStart: windowStart,
		ID:          convertId,
	}

	err = r.db.ResetRateLimit(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToResetRateLimit, err)
	}

	return nil
}

func NewRateLimitRepository(db *pgxpool.Pool) repository.RateLimitRepository {
	queries := gen.New(db)
	return &rateLimitRepository{db: queries}
}
