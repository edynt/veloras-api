package http

import (
	"github.com/edynnt/veloras-api/internal/auth/application/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (ah *AuthHandler) RegisterUser(ctx *gin.Context) (res interface{}, err error) {
	return nil, nil
}
