package http

import (
	"github.com/edynnt/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, handler *AuthHandler) {
	auth := rg.Group("/auth")
	// user registration
	auth.POST("/register", response.Wrap(handler.RegisterUser))
	auth.GET("/verify/:userId/:code", response.Wrap(handler.VerifyUser))
	auth.POST("/login", response.Wrap(handler.LoginUser))
}
