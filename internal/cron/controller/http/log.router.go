package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterLogRoutes(r *gin.RouterGroup, handler *LogHandler) {
	logGroup := r.Group("/admin/logs")
	{
		logGroup.GET("/files", handler.GetLogFiles)
		logGroup.GET("/content", handler.GetLogContent)
		logGroup.POST("/cleanup", handler.CleanupOldLogs)
	}
}
