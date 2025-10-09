package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/edynnt/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	LogDir string
}

func NewLogHandler(logDir string) *LogHandler {
	return &LogHandler{
		LogDir: logDir,
	}
}

// GetLogFiles returns list of available log files
// @Summary Get log files
// @Description Get list of available log files
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /admin/logs/files [get]
func (h *LogHandler) GetLogFiles(c *gin.Context) {
	files, err := h.getLogFiles()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to get log files", err.Error())
		return
	}

	response.SuccessResponse(c, files)
}

// GetLogContent returns content of a specific log file
// @Summary Get log content
// @Description Get content of a specific log file
// @Tags admin
// @Accept json
// @Produce json
// @Param filename query string true "Log filename"
// @Param lines query int false "Number of lines to return (default: 100)"
// @Success 200 {object} map[string]interface{}
// @Router /admin/logs/content [get]
func (h *LogHandler) GetLogContent(c *gin.Context) {
	filename := c.Query("filename")
	if filename == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "Filename is required", "")
		return
	}

	lines := c.DefaultQuery("lines", "100")
	lineCount := 100
	if l, err := fmt.Sscanf(lines, "%d", &lineCount); err != nil || l != 1 {
		lineCount = 100
	}

	content, err := h.getLogContent(filename, lineCount)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to read log file", err.Error())
		return
	}

	response.SuccessResponse(c, content)
}

// CleanupOldLogs removes log files older than specified days
// @Summary Cleanup old logs
// @Description Remove log files older than specified days (default: 30)
// @Tags admin
// @Accept json
// @Produce json
// @Param days query int false "Number of days to keep logs (default: 30)"
// @Success 200 {object} map[string]interface{}
// @Router /admin/logs/cleanup [post]
func (h *LogHandler) CleanupOldLogs(c *gin.Context) {
	days := c.DefaultQuery("days", "30")
	dayCount := 30
	if d, err := fmt.Sscanf(days, "%d", &dayCount); err != nil || d != 1 {
		dayCount = 30
	}

	result, err := h.cleanupOldLogs(dayCount)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to cleanup old logs", err.Error())
		return
	}

	response.SuccessResponse(c, result)
}

// getLogFiles returns list of log files in the log directory
func (h *LogHandler) getLogFiles() ([]LogFile, error) {
	files, err := ioutil.ReadDir(h.LogDir)
	if err != nil {
		return nil, err
	}

	var logFiles []LogFile
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".log") {
			logFiles = append(logFiles, LogFile{
				Name:     file.Name(),
				Size:     file.Size(),
				Modified: file.ModTime(),
			})
		}
	}

	// Sort by modification time (newest first)
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].Modified.After(logFiles[j].Modified)
	})

	return logFiles, nil
}

// getLogContent returns content of a log file
func (h *LogHandler) getLogContent(filename string, lineCount int) (*LogContent, error) {
	// Validate filename to prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return nil, fmt.Errorf("invalid filename")
	}

	filePath := filepath.Join(h.LogDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found")
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	// Get last N lines
	start := 0
	if len(lines) > lineCount {
		start = len(lines) - lineCount
	}

	selectedLines := lines[start:]

	return &LogContent{
		Filename:      filename,
		Lines:         selectedLines,
		TotalLines:    len(lines),
		ReturnedLines: len(selectedLines),
		LastModified:  time.Now(), // We could get this from file stat if needed
	}, nil
}

// cleanupOldLogs removes log files older than specified days
func (h *LogHandler) cleanupOldLogs(days int) (*LogCleanupResult, error) {
	cutoffTime := time.Now().AddDate(0, 0, -days)
	
	files, err := ioutil.ReadDir(h.LogDir)
	if err != nil {
		return nil, err
	}

	var deletedFiles []string
	var totalSizeDeleted int64
	var errors []string

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".log") {
			if file.ModTime().Before(cutoffTime) {
				filePath := filepath.Join(h.LogDir, file.Name())
				if err := os.Remove(filePath); err != nil {
					errors = append(errors, fmt.Sprintf("Failed to delete %s: %v", file.Name(), err))
				} else {
					deletedFiles = append(deletedFiles, file.Name())
					totalSizeDeleted += file.Size()
				}
			}
		}
	}

	return &LogCleanupResult{
		DeletedFiles:     deletedFiles,
		TotalFilesDeleted: len(deletedFiles),
		TotalSizeDeleted: totalSizeDeleted,
		CutoffDate:       cutoffTime,
		Errors:           errors,
	}, nil
}

type LogFile struct {
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
}

type LogContent struct {
	Filename      string    `json:"filename"`
	Lines         []string  `json:"lines"`
	TotalLines    int       `json:"total_lines"`
	ReturnedLines int       `json:"returned_lines"`
	LastModified  time.Time `json:"last_modified"`
}

type LogCleanupResult struct {
	DeletedFiles      []string  `json:"deleted_files"`
	TotalFilesDeleted int       `json:"total_files_deleted"`
	TotalSizeDeleted  int64     `json:"total_size_deleted"`
	CutoffDate        time.Time `json:"cutoff_date"`
	Errors            []string  `json:"errors"`
}
