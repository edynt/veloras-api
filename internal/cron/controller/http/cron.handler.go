package http

import (
	"fmt"
	"net/http"

	"github.com/edynt/chogiare/veloras-api/internal/cron"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

type CronHandler struct {
	Scheduler *cron.Scheduler
}

func NewCronHandler(scheduler *cron.Scheduler) *CronHandler {
	return &CronHandler{
		Scheduler: scheduler,
	}
}

// GetScheduledJobs returns information about scheduled cron jobs
// @Summary Get scheduled cron jobs
// @Description Get information about all scheduled cron jobs
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /admin/cron/jobs [get]
func (h *CronHandler) GetScheduledJobs(c *gin.Context) {
	jobs := h.Scheduler.GetScheduledJobs()

	response.SuccessResponse(c, jobs)
}

// RunManualCleanup runs cleanup manually
// @Summary Run manual cleanup
// @Description Manually trigger cleanup of expired data
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/cron/cleanup [post]
func (h *CronHandler) RunManualCleanup(c *gin.Context) {
	if err := h.Scheduler.RunManualCleanup(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Manual cleanup failed", err.Error())
		return
	}

	response.SuccessResponse(c, "Cleanup completed successfully")
}

// RunManualLogCleanup runs log cleanup manually
// @Summary Run manual log cleanup
// @Description Manually trigger cleanup of old log files
// @Tags admin
// @Accept json
// @Produce json
// @Param days query int false "Number of days to keep logs (default: 30)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/cron/log-cleanup [post]
func (h *CronHandler) RunManualLogCleanup(c *gin.Context) {
	days := c.DefaultQuery("days", "30")
	dayCount := 30
	if d, err := fmt.Sscanf(days, "%d", &dayCount); err != nil || d != 1 {
		dayCount = 30
	}

	if err := h.Scheduler.RunManualLogCleanup(dayCount); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Manual log cleanup failed", err.Error())
		return
	}

	response.SuccessResponse(c, "Log cleanup completed successfully")
}
