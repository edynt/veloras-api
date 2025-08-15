package http

import (
	"github.com/edynnt/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func SetupOrganizationRoutes(router *gin.RouterGroup, handler *OrganizationHandler) {
	organizations := router.Group("/organizations")
	{
		// Organization CRUD operations
		organizations.POST("", response.Wrap(handler.CreateOrganization))
		organizations.GET("/:id", response.Wrap(handler.GetOrganizationById))
		organizations.PUT("/:id", response.Wrap(handler.UpdateOrganization))
		organizations.DELETE("/:id", response.Wrap(handler.DeleteOrganization))
		
		// Organization hierarchy
		organizations.GET("/parent/:parentId", response.Wrap(handler.GetOrganizationsByParent))
		
		// User management
		organizations.POST("/assign-user", response.Wrap(handler.AssignUserToOrganization))
		organizations.GET("/:organizationId/users", response.Wrap(handler.GetOrganizationUsers))
	}
}
