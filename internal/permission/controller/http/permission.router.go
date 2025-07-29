package http

import (
	"github.com/edynnt/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterPermissionRoutes(rg *gin.RouterGroup, handler *PermissionHandler) {
	permission := rg.Group("/permissions")

	permission.GET("/", response.Wrap(handler.GetPermissions))
}
