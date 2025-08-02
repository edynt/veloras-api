package http

import (
	"github.com/edynnt/veloras-api/internal/middleware"
	"github.com/edynnt/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterPermissionRoutes(rg *gin.RouterGroup, handler *PermissionHandler) {
	permission := rg.Group("/permissions")

	permission.Use(middleware.AuthenMiddleware())
	permission.GET("/", response.Wrap(handler.GetPermissions))
	permission.POST("/", response.Wrap(handler.CreatePermission))
	permission.PUT("/:id", response.Wrap(handler.UpdatePermission))
}
