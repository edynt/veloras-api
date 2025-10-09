package http

import (
	"github.com/edynnt/veloras-api/internal/middleware"
	"github.com/edynnt/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, handler *UserHandler) {
	users := rg.Group("/users")
	users.Use(middleware.AuthenMiddleware())
	
	// Get all users (admin or authenticated users)
	users.GET("", response.Wrap(handler.GetAllUsers))
	
	// Get current user profile
	users.GET("/me", response.Wrap(handler.GetCurrentUser))
	
	// Update current user profile
	users.PUT("/me", response.Wrap(handler.UpdateUserProfile))
	
	// Get user by ID
	users.GET("/:id", response.Wrap(handler.GetUserByID))
}
