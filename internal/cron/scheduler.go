package cron

import (
	"context"
	"time"

	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/edynt/chogiare/veloras-api/pkg/global"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Scheduler struct {
	cron              *cron.Cron
	cleanupService    *CleanupService
	logCleanupService *LogCleanupService
	ctx               context.Context
	cancel            context.CancelFunc
}

func NewScheduler(queries *gen.Queries) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	// Create cron scheduler with timezone support
	c := cron.New(cron.WithLocation(time.UTC))

	cleanupService := NewCleanupService(queries)
	logCleanupService := NewLogCleanupService("./storage/logs")

	return &Scheduler{
		cron:              c,
		cleanupService:    cleanupService,
		logCleanupService: logCleanupService,
		ctx:               ctx,
		cancel:            cancel,
	}
}

// Start starts the cron scheduler
func (s *Scheduler) Start() error {
	if !global.Config.Cron.Enabled {
		global.Logger.Info("Cron scheduler is disabled in configuration")
		return nil
	}

	global.Logger.Info("Starting cron scheduler")

	// Schedule daily cleanup at 00:01 UTC
	_, err := s.cron.AddFunc("1 0 * * *", s.dailyCleanupJob)
	if err != nil {
		global.Logger.Error("Failed to schedule daily cleanup job", zap.Error(err))
		return err
	}

	// Schedule log cleanup at 01:00 UTC (after daily cleanup)
	_, err = s.cron.AddFunc("0 1 * * *", s.logCleanupJob)
	if err != nil {
		global.Logger.Error("Failed to schedule log cleanup job", zap.Error(err))
		return err
	}

	// Schedule weekly stats report (optional)
	_, err = s.cron.AddFunc("0 1 * * 1", s.weeklyStatsJob)
	if err != nil {
		global.Logger.Error("Failed to schedule weekly stats job", zap.Error(err))
		return err
	}

	// Start the cron scheduler
	s.cron.Start()

	global.Logger.Info("Cron scheduler started successfully",
		zap.String("daily_cleanup", "00:01 UTC"),
		zap.String("log_cleanup", "01:00 UTC"),
		zap.String("weekly_stats", "01:00 UTC Monday"))

	return nil
}

// Stop stops the cron scheduler
func (s *Scheduler) Stop() {
	global.Logger.Info("Stopping cron scheduler")

	ctx := s.cron.Stop()
	<-ctx.Done()

	s.cancel()

	global.Logger.Info("Cron scheduler stopped")
}

// dailyCleanupJob runs the daily cleanup of expired data
func (s *Scheduler) dailyCleanupJob() {
	global.Logger.Info("Starting daily cleanup job",
		zap.String("scheduled_time", time.Now().Format("2006-01-02 15:04:05")))

	// Create a new context with timeout for the cleanup job
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Minute)
	defer cancel()

	// Run cleanup
	if err := s.cleanupService.CleanupAllExpiredData(ctx); err != nil {
		global.Logger.Error("Daily cleanup job failed", zap.Error(err))
	} else {
		global.Logger.Info("Daily cleanup job completed successfully")
	}
}

// logCleanupJob runs daily log cleanup to remove old log files
func (s *Scheduler) logCleanupJob() {
	global.Logger.Info("Starting log cleanup job",
		zap.String("scheduled_time", time.Now().Format("2006-01-02 15:04:05")))

	// Cleanup logs older than 30 days
	if err := s.logCleanupService.CleanupOldLogs(30); err != nil {
		global.Logger.Error("Log cleanup job failed", zap.Error(err))
	} else {
		global.Logger.Info("Log cleanup job completed successfully")
	}
}

// weeklyStatsJob runs weekly statistics collection
func (s *Scheduler) weeklyStatsJob() {
	global.Logger.Info("Starting weekly stats job",
		zap.String("scheduled_time", time.Now().Format("2006-01-02 15:04:05")))

	// Create a new context with timeout for the stats job
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Minute)
	defer cancel()

	// Get cleanup stats
	stats, err := s.cleanupService.GetCleanupStats(ctx)
	if err != nil {
		global.Logger.Error("Weekly stats job failed", zap.Error(err))
		return
	}

	// Log weekly stats
	global.Logger.Info("Weekly cleanup statistics",
		zap.Int64("expired_sessions", stats.ExpiredSessions),
		zap.Int64("expired_email_verifications", stats.ExpiredEmailVerifications),
		zap.Int64("expired_password_resets", stats.ExpiredPasswordResets),
		zap.String("timestamp", stats.Timestamp.Format("2006-01-02 15:04:05")))
}

// RunManualCleanup runs cleanup manually (for testing or admin purposes)
func (s *Scheduler) RunManualCleanup() error {
	global.Logger.Info("Running manual cleanup",
		zap.String("start_time", time.Now().Format("2006-01-02 15:04:05")))

	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Minute)
	defer cancel()

	if err := s.cleanupService.CleanupAllExpiredData(ctx); err != nil {
		global.Logger.Error("Manual cleanup failed", zap.Error(err))
		return err
	}

	global.Logger.Info("Manual cleanup completed successfully")
	return nil
}

// RunManualLogCleanup runs log cleanup manually (for testing or admin purposes)
func (s *Scheduler) RunManualLogCleanup(days int) error {
	global.Logger.Info("Running manual log cleanup",
		zap.String("start_time", time.Now().Format("2006-01-02 15:04:05")),
		zap.Int("days", days))

	if err := s.logCleanupService.CleanupOldLogs(days); err != nil {
		global.Logger.Error("Manual log cleanup failed", zap.Error(err))
		return err
	}

	global.Logger.Info("Manual log cleanup completed successfully")
	return nil
}

// GetScheduledJobs returns information about scheduled jobs
func (s *Scheduler) GetScheduledJobs() []ScheduledJob {
	entries := s.cron.Entries()
	jobs := make([]ScheduledJob, len(entries))

	for i, entry := range entries {
		jobs[i] = ScheduledJob{
			ID:       entry.ID,
			Schedule: "Daily at 00:01 UTC", // Fixed schedule description
			Next:     entry.Next,
			Prev:     entry.Prev,
		}
	}

	return jobs
}

type ScheduledJob struct {
	ID       cron.EntryID `json:"id"`
	Schedule string       `json:"schedule"`
	Next     time.Time    `json:"next_run"`
	Prev     time.Time    `json:"last_run"`
}
