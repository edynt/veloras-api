package http

import (
	"github.com/edynnt/veloras-api/internal/middleware"
	"github.com/edynnt/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, handler *AuthHandler) {
	auth := rg.Group("/auth")
	// user registration
	auth.POST("/register", response.Wrap(handler.RegisterUser))
	auth.GET("/verify/:userId/:code", response.Wrap(handler.VerifyUser))
	auth.POST("/login", response.Wrap(handler.LoginUser))
	auth.POST("/refresh", response.Wrap(handler.RefreshToken))
	auth.POST("/forgot-password", response.Wrap(handler.ForgotPassword))
	auth.POST("/reset-password", response.Wrap(handler.ResetPassword))

	// protected routes
	protected := auth.Group("")
	protected.Use(middleware.AuthenMiddleware())
	protected.POST("/logout", response.Wrap(handler.Logout))
	protected.POST("/change-password", response.Wrap(handler.ChangePassword))
}
