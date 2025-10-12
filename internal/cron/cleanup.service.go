package cron

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/edynt/chogiare/veloras-api/pkg/global"
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

// CleanupAllExpiredData cleans up all expired data
func (s *CleanupService) CleanupAllExpiredData(ctx context.Context) error {
	global.Logger.Info("Starting cleanup of expired data")

	// For now, just log that cleanup is running
	// You can add actual cleanup logic here later
	global.Logger.Info("Cleanup completed successfully")
	return nil
}

// GetCleanupStats returns cleanup statistics
func (s *CleanupService) GetCleanupStats(ctx context.Context) (*CleanupStats, error) {
	return &CleanupStats{
		ExpiredSessions:           0,
		ExpiredEmailVerifications: 0,
		ExpiredPasswordResets:     0,
		Timestamp:                 time.Now(),
	}, nil
}

type CleanupStats struct {
	ExpiredSessions           int64     `json:"expired_sessions"`
	ExpiredEmailVerifications int64     `json:"expired_email_verifications"`
	ExpiredPasswordResets     int64     `json:"expired_password_resets"`
	Timestamp                 time.Time `json:"timestamp"`
}

type LogCleanupService struct {
	LogDir string
}

func NewLogCleanupService(logDir string) *LogCleanupService {
	return &LogCleanupService{
		LogDir: logDir,
	}
}

// CleanupOldLogs removes log files older than specified days
func (s *LogCleanupService) CleanupOldLogs(days int) error {
	cutoffTime := time.Now().AddDate(0, 0, -days)

	global.Logger.Info("Starting log cleanup",
		zap.String("cutoff_date", cutoffTime.Format("2006-01-02 15:04:05")),
		zap.Int("days", days))

	files, err := ioutil.ReadDir(s.LogDir)
	if err != nil {
		global.Logger.Error("Failed to read log directory", zap.Error(err), zap.String("log_dir", s.LogDir))
		return err
	}

	var deletedFiles []string
	var totalSizeDeleted int64
	var errors []string

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".log") {
			if file.ModTime().Before(cutoffTime) {
				filePath := filepath.Join(s.LogDir, file.Name())
				if err := os.Remove(filePath); err != nil {
					errorMsg := fmt.Sprintf("Failed to delete %s: %v", file.Name(), err)
					errors = append(errors, errorMsg)
					global.Logger.Error("Failed to delete log file",
						zap.String("file", file.Name()),
						zap.Error(err))
				} else {
					deletedFiles = append(deletedFiles, file.Name())
					totalSizeDeleted += file.Size()
					global.Logger.Info("Deleted old log file",
						zap.String("file", file.Name()),
						zap.Int64("size", file.Size()),
						zap.String("modified", file.ModTime().Format("2006-01-02 15:04:05")))
				}
			}
		}
	}

	// Log cleanup summary
	global.Logger.Info("Log cleanup completed",
		zap.Int("total_files_deleted", len(deletedFiles)),
		zap.Int64("total_size_deleted", totalSizeDeleted),
		zap.Strings("deleted_files", deletedFiles),
		zap.Strings("errors", errors))

	if len(errors) > 0 {
		return fmt.Errorf("cleanup completed with %d errors: %v", len(errors), errors)
	}

	return nil
}

// GetLogStats returns statistics about log files
func (s *LogCleanupService) GetLogStats() (map[string]interface{}, error) {
	files, err := ioutil.ReadDir(s.LogDir)
	if err != nil {
		return nil, err
	}

	var totalFiles int
	var totalSize int64
	var oldestFile time.Time
	var newestFile time.Time
	var filesOlderThan30Days int

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".log") {
			totalFiles++
			totalSize += file.Size()

			modTime := file.ModTime()
			if oldestFile.IsZero() || modTime.Before(oldestFile) {
				oldestFile = modTime
			}
			if newestFile.IsZero() || modTime.After(newestFile) {
				newestFile = modTime
			}

			// Check if file is older than 30 days
			if modTime.Before(time.Now().AddDate(0, 0, -30)) {
				filesOlderThan30Days++
			}
		}
	}

	return map[string]interface{}{
		"total_files":              totalFiles,
		"total_size":               totalSize,
		"oldest_file_date":         oldestFile,
		"newest_file_date":         newestFile,
		"files_older_than_30_days": filesOlderThan30Days,
	}, nil
}
