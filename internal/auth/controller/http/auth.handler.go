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
		UserID: userId,
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
// @Description Authenticate a user with username and password credentials
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

	loginRequest := appDto.LoginRequest{
		Username:  req.Username,
		Password:  req.Password,
		TwoFACode: req.TwoFACode,
	}

	account, err := ah.service.LoginUser(ctx, loginRequest)

	if err != nil {
		return nil, response.NewAPIError(http.StatusUnauthorized, msg.LoginFailed, err.Error())
	}

	return account, nil
}

// LogoutUser
// @Summary User logout
// @Description Logout a user by invalidating their refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.LogoutReq true "Logout request"
// @Success 200 {object} map[string]interface{} "Returns success message"
// @Failure 400 {object} response.APIError "Invalid request format"
// @Router /auth/logout [post]
func (ah *AuthHandler) LogoutUser(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.LogoutReq

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

	err = ah.service.LogoutUser(ctx, req.RefreshToken)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToLogout, err.Error())
	}

	return map[string]interface{}{
		"message": "Logged out successfully",
	}, nil
}

// RefreshToken
// @Summary Refresh access token
// @Description Get a new access token using a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.RefreshTokenReq true "Refresh token request"
// @Success 200 {object} appDto.RefreshTokenResponse "Returns new access token"
// @Failure 400 {object} response.APIError "Invalid request format"
// @Failure 401 {object} response.APIError "Invalid or expired refresh token"
// @Router /auth/refresh [post]
func (ah *AuthHandler) RefreshToken(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.RefreshTokenReq

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

	tokenResponse, err := ah.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, response.NewAPIError(http.StatusUnauthorized, msg.RefreshTokenInvalid, err.Error())
	}

	return tokenResponse, nil
}

// ForgotPassword
// @Summary Forgot password
// @Description Send password reset email to user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.ForgotPasswordReq true "Forgot password request"
// @Success 200 {object} map[string]interface{} "Returns success message"
// @Failure 400 {object} response.APIError "Invalid request format"
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
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToSendPasswordReset, err.Error())
	}

	return map[string]interface{}{
		"message": "Password reset email sent successfully",
	}, nil
}

// ResetPassword
// @Summary Reset password
// @Description Reset user password using reset token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.ResetPasswordReq true "Reset password request"
// @Success 200 {object} map[string]interface{} "Returns success message"
// @Failure 400 {object} response.APIError "Invalid request format"
// @Failure 401 {object} response.APIError "Invalid or expired reset token"
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

	err = ah.service.ResetPassword(ctx, req.Email, req.ResetToken, req.NewPassword)
	if err != nil {
		return nil, response.NewAPIError(http.StatusUnauthorized, msg.FailedToResetPassword, err.Error())
	}

	return map[string]interface{}{
		"message": "Password reset successfully",
	}, nil
}

// ChangePassword
// @Summary Change password
// @Description Change user password (requires current password)
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.ChangePasswordReq true "Change password request"
// @Success 200 {object} map[string]interface{} "Returns success message"
// @Failure 400 {object} response.APIError "Invalid request format"
// @Failure 401 {object} response.APIError "Invalid current password"
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

	// TODO: Get user ID from JWT token context
	userId := "current-user-id" // Placeholder

	err = ah.service.ChangePassword(ctx, userId, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return nil, response.NewAPIError(http.StatusUnauthorized, msg.FailedToChangePassword, err.Error())
	}

	return map[string]interface{}{
		"message": "Password changed successfully",
	}, nil
}

// Setup2FA
// @Summary Setup 2FA
// @Description Setup two-factor authentication for user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.TwoFASetupReq true "2FA setup request"
// @Success 200 {object} appDto.TwoFASetupResponse "Returns 2FA setup details"
// @Failure 400 {object} response.APIError "Invalid request format"
// @Router /auth/2fa/setup [post]
func (ah *AuthHandler) Setup2FA(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.TwoFASetupReq

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

	setupResponse, err := ah.service.Setup2FA(ctx, req.UserID)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToSetup2FA, err.Error())
	}

	return setupResponse, nil
}

// Verify2FA
// @Summary Verify 2FA code
// @Description Verify a 2FA code for user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.TwoFAVerifyReq true "2FA verification request"
// @Success 200 {object} appDto.TwoFAVerifyResponse "Returns verification result"
// @Failure 400 {object} response.APIError "Invalid request format"
// @Router /auth/2fa/verify [post]
func (ah *AuthHandler) Verify2FA(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.TwoFAVerifyReq

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

	verifyResponse, err := ah.service.Verify2FA(ctx, req.UserID, req.Code, req.Type)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToVerify2FA, err.Error())
	}

	return verifyResponse, nil
}

// Enable2FA
// @Summary Enable 2FA
// @Description Enable two-factor authentication for user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.TwoFAEnableReq true "2FA enable request"
// @Success 200 {object} map[string]interface{} "Returns success message"
// @Failure 400 {object} response.APIError "Invalid request format"
// @Router /auth/2fa/enable [post]
func (ah *AuthHandler) Enable2FA(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.TwoFAEnableReq

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

	err = ah.service.Enable2FA(ctx, req.UserID, req.Secret, req.Code)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToEnable2FA, err.Error())
	}

	return map[string]interface{}{
		"message": "2FA enabled successfully",
	}, nil
}

// Disable2FA
// @Summary Disable 2FA
// @Description Disable two-factor authentication for user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.TwoFADisableReq true "2FA disable request"
// @Success 200 {object} map[string]interface{} "Returns success message"
// @Failure 400 {object} response.APIError "Invalid request format"
// @Router /auth/2fa/disable [post]
func (ah *AuthHandler) Disable2FA(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.TwoFADisableReq

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

	err = ah.service.Disable2FA(ctx, req.UserID, req.Code)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToDisable2FA, err.Error())
	}

	return map[string]interface{}{
		"message": "2FA disabled successfully",
	}, nil
}

// Get2FAStatus
// @Summary Get 2FA status
// @Description Get two-factor authentication status for user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ctlDto.TwoFAStatusReq true "2FA status request"
// @Success 200 {object} appDto.TwoFAStatusResponse "Returns 2FA status"
// @Failure 400 {object} response.APIError "Invalid request format"
// @Router /auth/2fa/status [post]
func (ah *AuthHandler) Get2FAStatus(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.TwoFAStatusReq

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

	statusResponse, err := ah.service.Get2FAStatus(ctx, req.UserID)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToGet2FAStatus, err.Error())
	}

	return statusResponse, nil
}

// GetUserSessions
// @Summary Get user sessions
// @Description Get all active sessions for a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} appDto.SessionInfo "Returns list of user sessions"
// @Failure 400 {object} response.APIError "Invalid user ID"
// @Router /auth/sessions/{userId} [get]
func (ah *AuthHandler) GetUserSessions(ctx *gin.Context) (res interface{}, err error) {
	userId := ctx.Param("userId")

	sessions, err := ah.service.GetUserSessions(ctx, userId)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToGetUserSessions, err.Error())
	}

	return sessions, nil
}

// DeleteAllUserSessions
// @Summary Delete all user sessions
// @Description Delete all active sessions for a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} map[string]interface{} "Returns success message"
// @Failure 400 {object} response.APIError "Invalid user ID"
// @Router /auth/sessions/{userId} [delete]
func (ah *AuthHandler) DeleteAllUserSessions(ctx *gin.Context) (res interface{}, err error) {
	userId := ctx.Param("userId")

	err = ah.service.DeleteAllUserSessions(ctx, userId)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToDeleteAllUserSessions, err.Error())
	}

	return map[string]interface{}{
		"message": "All user sessions deleted successfully",
	}, nil
}
