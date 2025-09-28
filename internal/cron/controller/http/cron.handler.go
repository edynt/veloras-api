package http

import (
	"net/http"

	"github.com/edynnt/veloras-api/internal/cron"
	"github.com/edynnt/veloras-api/pkg/response"
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
// @Success 200 {object} response.Response{data=[]cron.ScheduledJob}
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
// @Success 200 {object} response.Response{data=string}
// @Failure 500 {object} response.Response{data=string}
// @Router /admin/cron/cleanup [post]
func (h *CronHandler) RunManualCleanup(c *gin.Context) {
	if err := h.Scheduler.RunManualCleanup(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Manual cleanup failed", err.Error())
		return
	}

	response.SuccessResponse(c, "Cleanup completed successfully")
}
