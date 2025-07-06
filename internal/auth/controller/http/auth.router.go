package http

import (
	"github.com/edynnt/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, handler *AuthHandler) {
	auth := rg.Group("/auth")
	// user registration
	auth.POST("/register", response.Wrap(handler.RegisterUser))
}
