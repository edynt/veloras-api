package http

import (
	"github.com/edynt/chogiare/veloras-api/internal/middleware"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(rg *gin.RouterGroup, handler *RoleHandler) {
	role := rg.Group("/roles")

	role.Use(middleware.AuthenMiddleware())
	role.POST("/", response.Wrap(handler.CreateRole))
	role.PUT("/:id", response.Wrap(handler.UpdateRole))
	role.DELETE("/:id", response.Wrap(handler.DeleteRole))
	role.GET("/:id", response.Wrap(handler.GetRole))
	role.GET("/", response.Wrap(handler.GetRoles))
	role.POST("/assign-permissions", response.Wrap(handler.AssignPermissions))
}
