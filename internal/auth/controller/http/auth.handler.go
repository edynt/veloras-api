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
// @Description Create a new account with email, password, and other details
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

// VerifyUser
// @Summary Verify user email
// @Description Verify a user's email address using the provided verification code
// @Tags Auth
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param code path string true "Verification code"
// @Success 200 {object} map[string]interface{} "Returns verification status"
// @Failure 400 {object} response.APIError "Invalid request or verification failed"
// @Router /auth/verify/{userId}/{code} [get]
func (ah *AuthHandler) VerifyUser(ctx *gin.Context) (res interface{}, err error) {

	userId := ctx.Param("userId")
	code := ctx.Param("code")

	verificationEmail := appDto.EmailVerification{
		UserID: utils.StringToInt(userId),
		Code:   utils.StringToInt(code),
	}

	isExist, err := ah.service.VerifyUser(ctx, verificationEmail)

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return isExist, nil
}

// LoginUser
// @Summary User login
// @Description Authenticate a user with email and password credentials
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.UserLoginReq true "User login request"
// @Success 200 {object} map[string]interface{} "Returns authenticated user account information"
// @Failure 400 {object} response.APIError "Invalid request format or validation errors"
// @Failure 401 {object} response.APIError "Login failed due to invalid credentials"
// @Router /auth/login [post]
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
		Email:    req.Email,
		Password: req.Password,
	}

	account, err := ah.service.LoginUser(ctx, requestAccount)

	if err != nil {
		return nil, response.NewAPIError(http.StatusUnauthorized, msg.LoginFailed, err.Error())
	}

	return account, nil
}

// RefreshToken
// @Summary Refresh access token
// @Description Exchange a refresh token for a new access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.RefreshTokenReq true "Refresh token request"
// @Success 200 {object} ctlDto.RefreshTokenRes "Returns new access token and refresh token"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Refresh failed"
// @Router /auth/refresh [post]
func (ah *AuthHandler) RefreshToken(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.RefreshTokenReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	if req.RefreshToken == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, "refresh_token is required")
	}

	out, err := ah.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, response.NewAPIError(http.StatusUnauthorized, msg.LoginFailed, err.Error())
	}

	return ctlDto.RefreshTokenRes{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}, nil
}

// Logout
// @Summary User logout
// @Description Logout the current user and invalidate their session
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Returns logout success message"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /auth/logout [post]
func (ah *AuthHandler) Logout(ctx *gin.Context) (res interface{}, err error) {
	err = ah.service.Logout(ctx)
	if err != nil {
		return nil, response.NewAPIError(http.StatusUnauthorized, msg.LogoutFailed, err.Error())
	}

	return map[string]interface{}{
		"message": "Logout successful",
	}, nil
}

// ChangePassword
// @Summary Change user password
// @Description Change the current user's password by providing current password and new password
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ctlDto.ChangePasswordReq true "Change password request"
// @Success 200 {object} map[string]interface{} "Returns success message"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized or current password incorrect"
// @Router /auth/change-password [post]
func (ah *AuthHandler) ChangePassword(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.ChangePasswordReq

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

	// Get user ID from context
	subject := ctx.Value("subjectID")
	if subject == nil {
		return nil, response.NewAPIError(http.StatusUnauthorized, msg.Unauthorized, "")
	}

	userID := utils.StringToInt(subject.(string))
	if userID == 0 {
		return nil, response.NewAPIError(http.StatusUnauthorized, msg.Unauthorized, msg.UserIdInvalid)
	}

	err = ah.service.ChangePassword(ctx, userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.ChangePasswordFailed, err.Error())
	}

	return map[string]interface{}{
		"message": msg.PasswordChangedSuccessfully,
	}, nil
}

// ForgotPassword
// @Summary Request password reset
// @Description Send password reset instructions to user's email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.ForgotPasswordReq true "Forgot password request"
// @Success 200 {object} map[string]interface{} "Returns success message"
// @Failure 400 {object} response.APIError "Invalid request"
// @Router /auth/forgot-password [post]
func (ah *AuthHandler) ForgotPassword(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.ForgotPasswordReq

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

	err = ah.service.ForgotPassword(ctx, req.Email)
	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.ForgotPasswordFailed, err.Error())
	}

	return map[string]interface{}{
		"message": msg.ForgotPasswordSuccess,
	}, nil
}

// ResetPassword
// @Summary Reset password with token
// @Description Reset user password using the token received via email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.ResetPasswordReq true "Reset password request"
// @Success 200 {object} map[string]interface{} "Returns success message"
// @Failure 400 {object} response.APIError "Invalid request or token"
// @Router /auth/reset-password [post]
func (ah *AuthHandler) ResetPassword(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.ResetPasswordReq

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

	err = ah.service.ResetPassword(ctx, req.Token, req.NewPassword)
	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.ResetPasswordFailed, err.Error())
	}

	return map[string]interface{}{
		"message": msg.ResetPasswordSuccess,
	}, nil
}
