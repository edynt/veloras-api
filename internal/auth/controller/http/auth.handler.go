package http

import (
	"fmt"
	"net/http"

	"github.com/edynnt/veloras-api/internal/auth/application/service"
	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	ctlDto "github.com/edynnt/veloras-api/internal/auth/controller/dto"
	"github.com/edynnt/veloras-api/pkg/response"
	"github.com/edynnt/veloras-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (ah *AuthHandler) RegisterUser(ctx *gin.Context) (res interface{}, err error) {

	fmt.Println(("----> RegisterUser"))

	var req ctlDto.UserRegisterReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", err.Error())
	}

	validation, exists := ctx.Get("validation")

	if !exists {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Validation not found in context")
	}

	if apiErr := utils.ValidateStruct(req, validation.(*validator.Validate)); apiErr != nil {
		return nil, apiErr
	}

	account := appDto.AccountAppDTO{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Lang:     req.Language,
	}

	return account, nil
}
