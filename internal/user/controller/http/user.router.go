package http

import (
	"github.com/edynt/chogiare/veloras-api/internal/middleware"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, handler *UserHandler) {
	users := rg.Group("/users")

	// Public routes (no authentication required)
	// None for now

	// Protected routes (require authentication)
	usersProtected := users.Group("")
	usersProtected.Use(middleware.AuthenMiddleware())
	{
		// Get current user profile
		usersProtected.GET("/me", response.Wrap(handler.GetCurrentUser))

		// Update current user profile
		usersProtected.PUT("/me", response.Wrap(handler.UpdateUserProfile))
	}

	// Admin routes (require admin role)
	usersAdmin := users.Group("")
	usersAdmin.Use(middleware.AuthenMiddleware())
	usersAdmin.Use(middleware.RequireRole("admin"))
	{
		// Get all users (admin only)
		usersAdmin.GET("", response.Wrap(handler.GetAllUsers))

		// Get user by ID (admin only)
		usersAdmin.GET("/:id", response.Wrap(handler.GetUserByID))
	}
}
