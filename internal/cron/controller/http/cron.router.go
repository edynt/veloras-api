package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterCronRoutes(r *gin.RouterGroup, handler *CronHandler) {
	cronGroup := r.Group("/admin/cron")
	{
		cronGroup.GET("/jobs", handler.GetScheduledJobs)
		cronGroup.POST("/cleanup", handler.RunManualCleanup)
		cronGroup.POST("/log-cleanup", handler.RunManualLogCleanup)
	}
}
