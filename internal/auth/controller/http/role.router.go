package http

import (
	"github.com/edynnt/veloras-api/internal/middleware"
	"github.com/edynnt/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(rg *gin.RouterGroup, handler *RoleHandler) {
	role := rg.Group("/roles")

	role.Use(middleware.AuthenMiddleware())
	role.POST("/", response.Wrap(handler.CreateRole))
	role.PUT("/:id", response.Wrap(handler.UpdateRole))
}
