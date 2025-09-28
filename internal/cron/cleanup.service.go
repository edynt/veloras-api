package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/edynnt/veloras-api/internal/shared/gen"
	"github.com/edynnt/veloras-api/pkg/global"
	"go.uber.org/zap"
)

type CleanupService struct {
	queries *gen.Queries
}

func NewCleanupService(queries *gen.Queries) *CleanupService {
	return &CleanupService{
		queries: queries,
	}
}

// CleanupExpiredSessions removes expired sessions
func (s *CleanupService) CleanupExpiredSessions(ctx context.Context) error {
	currentTime := time.Now().Unix()

	// Count expired sessions before deletion
	expiredCount, _ := s.queries.CountExpiredSessions(ctx, currentTime)

	err := s.queries.DeleteExpiredSessions(ctx, currentTime)
	if err != nil {
		global.Logger.Error("Failed to cleanup expired sessions", zap.Error(err))
		return err
	}

	global.Logger.Info("Cleanup expired sessions completed",
		zap.Int64("deleted_count", expiredCount),
		zap.String("cleanup_time", time.Now().Format("2006-01-02 15:04:05")))

	return nil
}

// CleanupExpiredEmailVerifications removes expired email verification codes
func (s *CleanupService) CleanupExpiredEmailVerifications(ctx context.Context) error {
	currentTime := time.Now().Unix()

	// Count expired verifications before deletion
	expiredCount, _ := s.queries.CountExpiredEmailVerifications(ctx, currentTime)

	err := s.queries.DeleteExpiredEmailVerifications(ctx, currentTime)
	if err != nil {
		global.Logger.Error("Failed to cleanup expired email verifications", zap.Error(err))
		return err
	}

	global.Logger.Info("Cleanup expired email verifications completed",
		zap.Int64("deleted_count", expiredCount),
		zap.String("cleanup_time", time.Now().Format("2006-01-02 15:04:05")))

	return nil
}

// CleanupExpiredPasswordResets removes expired password reset tokens
func (s *CleanupService) CleanupExpiredPasswordResets(ctx context.Context) error {
	currentTime := time.Now().Unix()

	// Count expired password resets before deletion
	expiredCount, _ := s.queries.CountExpiredPasswordResets(ctx, currentTime)

	err := s.queries.DeleteExpiredPasswordResets(ctx, currentTime)
	if err != nil {
		global.Logger.Error("Failed to cleanup expired password resets", zap.Error(err))
		return err
	}

	global.Logger.Info("Cleanup expired password resets completed",
		zap.Int64("deleted_count", expiredCount),
		zap.String("cleanup_time", time.Now().Format("2006-01-02 15:04:05")))

	return nil
}

// CleanupAllExpiredData runs all cleanup operations
func (s *CleanupService) CleanupAllExpiredData(ctx context.Context) error {
	global.Logger.Info("Starting daily cleanup of expired data",
		zap.String("start_time", time.Now().Format("2006-01-02 15:04:05")))

	var errors []error

	// Cleanup sessions
	if global.Config.Cron.CleanupSessions {
		if err := s.CleanupExpiredSessions(ctx); err != nil {
			errors = append(errors, fmt.Errorf("sessions cleanup failed: %w", err))
		}
	}

	// Cleanup email verifications
	if global.Config.Cron.CleanupVerifications {
		if err := s.CleanupExpiredEmailVerifications(ctx); err != nil {
			errors = append(errors, fmt.Errorf("email verifications cleanup failed: %w", err))
		}
	}

	// Cleanup password resets
	if global.Config.Cron.CleanupPasswordResets {
		if err := s.CleanupExpiredPasswordResets(ctx); err != nil {
			errors = append(errors, fmt.Errorf("password resets cleanup failed: %w", err))
		}
	}

	// Log summary
	if len(errors) > 0 {
		global.Logger.Error("Daily cleanup completed with errors",
			zap.Int("error_count", len(errors)),
			zap.String("end_time", time.Now().Format("2006-01-02 15:04:05")))
		for i, err := range errors {
			global.Logger.Error(fmt.Sprintf("Cleanup error %d", i+1), zap.Error(err))
		}
		return fmt.Errorf("cleanup completed with %d errors", len(errors))
	}

	global.Logger.Info("Daily cleanup completed successfully",
		zap.String("end_time", time.Now().Format("2006-01-02 15:04:05")))

	return nil
}

// GetCleanupStats returns statistics about expired data
func (s *CleanupService) GetCleanupStats(ctx context.Context) (*CleanupStats, error) {
	currentTime := time.Now().Unix()

	stats := &CleanupStats{
		Timestamp: time.Now(),
	}

	// Count expired sessions
	sessionCount, err := s.queries.CountExpiredSessions(ctx, currentTime)
	if err != nil {
		global.Logger.Error("Failed to count expired sessions", zap.Error(err))
	} else {
		stats.ExpiredSessions = sessionCount
	}

	// Count expired email verifications
	verificationCount, err := s.queries.CountExpiredEmailVerifications(ctx, currentTime)
	if err != nil {
		global.Logger.Error("Failed to count expired email verifications", zap.Error(err))
	} else {
		stats.ExpiredEmailVerifications = verificationCount
	}

	// Count expired password resets
	passwordResetCount, err := s.queries.CountExpiredPasswordResets(ctx, currentTime)
	if err != nil {
		global.Logger.Error("Failed to count expired password resets", zap.Error(err))
	} else {
		stats.ExpiredPasswordResets = passwordResetCount
	}

	return stats, nil
}

type CleanupStats struct {
	Timestamp                 time.Time `json:"timestamp"`
	ExpiredSessions           int64     `json:"expired_sessions"`
	ExpiredEmailVerifications int64     `json:"expired_email_verifications"`
	ExpiredPasswordResets     int64     `json:"expired_password_resets"`
}
