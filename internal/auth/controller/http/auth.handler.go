package http

import (
	"net/http"

	"github.com/edynnt/veloras-api/internal/auth/application/service"
	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	ctlDto "github.com/edynnt/veloras-api/internal/auth/controller/dto"
	"github.com/edynnt/veloras-api/pkg/response"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// RegisterUser
// @Summary Register a new user
// @Description Create a new account with username, password, email, and other details
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.UserRegisterReq true "User registration request"
// @Success 200 {object} map[string]interface{} "Returns created account ID"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 409 {object} response.APIError "Registration failed due to conflict"
// @Router /auth/register [post]
func (ah *AuthHandler) RegisterUser(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.UserRegisterReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	validation, exists := ctx.Get("validation")
	if !exists {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, msg.ValidationNotFoundInContext)
	}

	if apiErr := utils.ValidateStruct(req, validation.(*validator.Validate)); apiErr != nil {
		return nil, apiErr
	}

	account := appDto.AccountAppDTO{
		Username:    req.Username,
		Password:    req.Password,
		Email:       req.Email,
		Language:    req.Language,
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
	}

	accountId, err := ah.service.CreateUser(ctx, account)
	if err != nil {
		return nil, response.NewAPIError(http.StatusConflict, msg.RegistrationFailed, err.Error())
	}

	return accountId, nil
}

func (ah *AuthHandler) VerifyUser(ctx *gin.Context) (res interface{}, err error) {

	userId := ctx.Param("userId")
	code := ctx.Param("code")

	verificationEmail := appDto.EmailVerification{
		UserID: userId,
		Code:   utils.StringToInt(code),
	}

	isExist, err := ah.service.VerifyUser(ctx, verificationEmail)

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return isExist, nil
}

func (ah *AuthHandler) LoginUser(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.UserLoginReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	validation, exists := ctx.Get("validation")

	if !exists {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, msg.ValidationNotFoundInContext)
	}

	if apiErr := utils.ValidateStruct(req, validation.(*validator.Validate)); apiErr != nil {
		return nil, apiErr
	}

	requestAccount := appDto.AccountAppDTO{
		Username: req.Username,
		Password: req.Password,
	}

	account, err := ah.service.LoginUser(ctx, requestAccount)

	if err != nil {
		return nil, response.NewAPIError(http.StatusUnauthorized, msg.LoginFailed, err.Error())
	}

	return account, nil
}
