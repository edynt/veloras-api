package service

import (
	"context"
	"fmt"
	"log"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	authRepo "github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/pkg/constants"
	"github.com/edynnt/veloras-api/pkg/global"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	authRepo authRepo.AuthRepository
}

// LoginUser implements AuthService.
func (as *authService) LoginUser(ctx context.Context, loginReq appDto.LoginRequest) (*appDto.LoginResponse, error) {
	// 1. check exists
	user, err := as.authRepo.GetUserByUsername(ctx, loginReq.Username)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToCheckUserNameExists, err)
	}

	if user == nil {
		return nil, fmt.Errorf(msg.UsernameNotFound)
	}

	// 2. check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		// Increment failed login attempts
		as.authRepo.IncrementFailedLoginAttempts(ctx, user.ID)
		return nil, fmt.Errorf("Invalid password")
	}

	// 3. check verified
	if !user.IsVerified {
		return nil, fmt.Errorf(msg.UserIsNotVerified)
	}

	// 4. check status active
	if user.Status != constants.ACTIVE {
		return nil, fmt.Errorf(msg.UserIsNotActive)
	}

	// 5. check if account is locked
	if user.LockedUntil > 0 && user.LockedUntil > utils.GetNowUnix() {
		return nil, fmt.Errorf(msg.AccountLocked)
	}

	// 6. check 2FA if enabled
	if user.TwoFAEnabled && loginReq.TwoFACode == "" {
		return nil, fmt.Errorf(msg.TwoFACodeRequired)
	}

	if user.TwoFAEnabled && loginReq.TwoFACode != "" {
		// Verify 2FA code
		valid, err := as.Validate2FACode(ctx, user.ID, loginReq.TwoFACode, "login")
		if err != nil || !valid {
			return nil, fmt.Errorf(msg.TwoFACodeInvalid)
		}
	}

	// 7. Reset failed login attempts
	as.authRepo.ResetFailedLoginAttempts(ctx, user.ID)

	// 8. Update last login
	as.authRepo.UpdateLastLogin(ctx, user.ID)

	// 9. Generate accessToken and refreshToken
	accessToken, err := utils.CreateToken(user.ID, false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToCreateToken, err)
	}

	refreshToken, err := utils.CreateToken(user.ID, true)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToCreateToken, err)
	}

	// 10. Save refresh token
	tokenExpiresAt := utils.AddDays(global.Config.JWT.RefreshTokenExpire)

	err = as.authRepo.SaveToken(ctx, &entity.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    tokenExpiresAt,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToSaveToken, err)
	}

	return &appDto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    tokenExpiresAt,
		User: appDto.UserResponse{
			ID:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			IsVerified:     user.IsVerified,
			TwoFAEnabled:   user.TwoFAEnabled,
			OrganizationID: user.OrganizationID,
		},
	}, nil
}

// VerifyUser implements AuthService.
func (as *authService) VerifyUser(ctx context.Context, verificationEmailAppDTO appDto.EmailVerification) (bool, error) {
	existsVerificationCode, err := as.authRepo.GetVerificationCode(ctx, verificationEmailAppDTO.UserID, verificationEmailAppDTO.Code)

	if err != nil {
		return false, fmt.Errorf("%s: %w", msg.FailedToGetVerificationCode, err)
	}

	now := utils.GetNowUnix()
	if existsVerificationCode.ExpiresAt < now {
		return false, fmt.Errorf(msg.CodeExpired)
	}

	err = as.authRepo.UpdateUserStatus(ctx, verificationEmailAppDTO.UserID, constants.ACTIVE)

	if err != nil {
		return false, fmt.Errorf("%s: %w", msg.FailedToUpdateUserStatus, err)
	}

	err = as.authRepo.ActiveUser(ctx, verificationEmailAppDTO.UserID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", msg.FailedToActiveUser, err)
	}

	err = as.authRepo.DeleteVerificationCode(ctx, verificationEmailAppDTO.UserID, verificationEmailAppDTO.Code)
	if err != nil {
		return false, fmt.Errorf("%s: %w", msg.FailedToDeleteVerificationCode, err)
	}

	return true, nil
}

// CreateUser implements AuthService.
func (as *authService) CreateUser(ctx context.Context, accountDto appDto.AccountAppDTO) (string, error) {
	//1. Check permissions -> event registered

	// 2. Check username exists
	exists, err := as.authRepo.UsernameExists(ctx, accountDto.Username)

	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.FailedToCheckUserNameExists, err)
	}
	if exists {
		return "", fmt.Errorf(msg.UsernameExists)
	}

	// 3. Check email exists
	exists, err = as.authRepo.EmailExists(ctx, accountDto.Email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.FailedToCheckEmailExists, err)
	}

	if exists {
		return "", fmt.Errorf(msg.EmailExists)
	}

	// 4. GenerateFromPassword
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(accountDto.Password), bcrypt.DefaultCost)
	if err != nil {
		// log.Printf("Error hashing password for user %s: %v", accountDto.Username, err)
		return "", fmt.Errorf("%s: %w", msg.FailedToSecurePassword, err) // Không lộ chi tiết lỗi hash
	}
	hashedPassword := string(hashedPasswordBytes)

	accountDto.Password = hashedPassword

	// 5. Insert account into database
	newAccountId, err := as.authRepo.CreateUser(ctx, &entity.Account{
		Username:    accountDto.Username,
		Email:       accountDto.Email,
		Password:    accountDto.Password,
		PhoneNumber: accountDto.PhoneNumber,
		FirstName:   accountDto.FirstName,
		LastName:    accountDto.LastName,
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.CouldNotCreateAccount, err)
	}

	if newAccountId == "" {
		return "", fmt.Errorf(msg.CouldNotCreateAccount)
	}

	codeGen := utils.GenerateSixDigitCode()
	as.authRepo.CreateVerificationCode(ctx, &entity.EmailVerification{
		UserID:    newAccountId,
		Code:      codeGen,
		ExpiresAt: utils.AddHours(1),
	})

	go utils.SendTemplateEmailOtp(
		[]string{accountDto.Email}, global.Config.SMTP.User,
		"otp-auth.html",
		map[string]interface{}{"Otp": codeGen},
	)

	// 6. Return account ID
	return newAccountId, nil
}

// GetUserByEmail implements AuthService.
func (as *authService) GetUserByEmail(ctx context.Context, email string) (*appDto.UserResponse, error) {
	user, err := as.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetUser, err)
	}

	if user == nil {
		return nil, fmt.Errorf(msg.UserNotFound)
	}

	return &appDto.UserResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		IsVerified:     user.IsVerified,
		TwoFAEnabled:   user.TwoFAEnabled,
		OrganizationID: user.OrganizationID,
	}, nil
}

// GetUserById implements AuthService.
func (as *authService) GetUserById(ctx context.Context, userId string) (*appDto.UserResponse, error) {
	user, err := as.authRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetUser, err)
	}

	if user == nil {
		return nil, fmt.Errorf(msg.UserNotFound)
	}

	return &appDto.UserResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		IsVerified:     user.IsVerified,
		TwoFAEnabled:   user.TwoFAEnabled,
		OrganizationID: user.OrganizationID,
	}, nil
}

// UpdateUserProfile implements AuthService.
func (as *authService) UpdateUserProfile(ctx context.Context, userId string, firstName, lastName, phoneNumber, language string) error {
	return as.authRepo.UpdateUserProfile(ctx, userId, firstName, lastName, phoneNumber, language)
}

// UpdateUserPassword implements AuthService.
func (as *authService) UpdateUserPassword(ctx context.Context, userId string, currentPassword, newPassword string) error {
	// Get current user to verify current password
	user, err := as.authRepo.GetUserById(ctx, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToGetUser, err)
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return fmt.Errorf("Current password is incorrect")
	}

	// Hash new password
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToSecurePassword, err)
	}

	// Update password
	return as.authRepo.UpdateUserPassword(ctx, userId, string(hashedPasswordBytes))
}

// LogoutUser implements AuthService.
func (as *authService) LogoutUser(ctx context.Context, refreshToken string) error {
	return as.authRepo.DeleteSessionByRefreshToken(ctx, refreshToken)
}

// RefreshToken implements AuthService.
func (as *authService) RefreshToken(ctx context.Context, refreshToken string) (*appDto.RefreshTokenResponse, error) {
	// Get session by refresh token
	session, err := as.authRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.RefreshTokenInvalid, err)
	}

	// Check if session is expired
	if session.ExpiresAt < utils.GetNowUnix() {
		return nil, fmt.Errorf("%s", msg.RefreshTokenExpired)
	}

	// Generate new access token
	accessToken, err := utils.CreateToken(session.UserID, false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToCreateToken, err)
	}

	return &appDto.RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   utils.AddDays(global.Config.JWT.AccessTokenExpire),
	}, nil
}

// ForgotPassword implements AuthService.
func (as *authService) ForgotPassword(ctx context.Context, email string) error {
	// Get user by email
	user, err := as.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToGetUser, err)
	}

	if user == nil {
		return fmt.Errorf(msg.UserNotFound)
	}

	// Generate reset token
	resetToken := utils.GenerateRandomString(32)
	expiresAt := utils.AddHours(24) // Token expires in 24 hours

	// Create password reset record
	err = as.authRepo.CreatePasswordReset(ctx, user.ID, resetToken, expiresAt)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToCreatePasswordReset, err)
	}

	// Send email with reset link
	err = as.sendPasswordResetEmail(user.Email, resetToken)
	if err != nil {
		log.Printf("Failed to send password reset email: %v", err)
		// Don't return error here as the reset token was created successfully
	}

	return nil
}

// ResetPassword implements AuthService.
func (as *authService) ResetPassword(ctx context.Context, email, resetToken, newPassword string) error {
	// Get user by email
	user, err := as.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToGetUser, err)
	}

	if user == nil {
		return fmt.Errorf(msg.UserNotFound)
	}

	// Get password reset record
	resetRecord, err := as.authRepo.GetPasswordReset(ctx, user.ID, resetToken)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToGetPasswordReset, err)
	}

	// Check if token is expired
	if resetRecord.ExpiresAt < utils.GetNowUnix() {
		return fmt.Errorf("%s", msg.PasswordResetTokenExpired)
	}

	// Hash new password
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToSecurePassword, err)
	}

	// Update password
	err = as.authRepo.UpdateUserPassword(ctx, user.ID, string(hashedPasswordBytes))
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdatePassword, err)
	}

	// Delete password reset record
	err = as.authRepo.DeletePasswordReset(ctx, user.ID, resetToken)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeletePasswordReset, err)
	}

	return nil
}

// ChangePassword implements AuthService.
func (as *authService) ChangePassword(ctx context.Context, userId, currentPassword, newPassword string) error {
	return as.UpdateUserPassword(ctx, userId, currentPassword, newPassword)
}

// Setup2FA implements AuthService.
func (as *authService) Setup2FA(ctx context.Context, userId string) (*appDto.TwoFASetupResponse, error) {
	// Generate 2FA secret
	secret := utils.GenerateRandomString(32)

	// Store secret temporarily (not enabled yet)
	err := as.authRepo.UpdateUser2FASecret(ctx, userId, secret)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToUpdate2FASecret, err)
	}

	// Generate QR code (placeholder for now)
	qrCode := fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=Veloras", userId, secret)

	return &appDto.TwoFASetupResponse{
		Secret: secret,
		QRCode: qrCode,
	}, nil
}

// Verify2FA implements AuthService.
func (as *authService) Verify2FA(ctx context.Context, userId, code, codeType string) (*appDto.TwoFAVerifyResponse, error) {
	valid, err := as.Validate2FACode(ctx, userId, code, codeType)
	if err != nil {
		return nil, err
	}

	if valid {
		return &appDto.TwoFAVerifyResponse{
			Success: true,
			Message: "2FA code verified successfully",
		}, nil
	}

	return &appDto.TwoFAVerifyResponse{
		Success: false,
		Message: "Invalid 2FA code",
	}, nil
}

// Enable2FA implements AuthService.
func (as *authService) Enable2FA(ctx context.Context, userId, secret, code string) error {
	// Verify the code first
	valid, err := as.Validate2FACode(ctx, userId, code, "setup")
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf(msg.TwoFACodeInvalid)
	}

	// Enable 2FA
	return as.authRepo.UpdateUser2FAEnabled(ctx, userId, true)
}

// Disable2FA implements AuthService.
func (as *authService) Disable2FA(ctx context.Context, userId, code string) error {
	// Verify the code first
	valid, err := as.Validate2FACode(ctx, userId, code, "setup")
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf(msg.TwoFACodeInvalid)
	}

	// Disable 2FA
	return as.authRepo.UpdateUser2FAEnabled(ctx, userId, false)
}

// Get2FAStatus implements AuthService.
func (as *authService) Get2FAStatus(ctx context.Context, userId string) (*appDto.TwoFAStatusResponse, error) {
	user, err := as.authRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetUser, err)
	}

	return &appDto.TwoFAStatusResponse{
		Enabled: user.TwoFAEnabled,
		Secret:  user.TwoFASecret,
	}, nil
}

// GetUserSessions implements AuthService.
func (as *authService) GetUserSessions(ctx context.Context, userId string) ([]*appDto.SessionInfo, error) {
	sessions, err := as.authRepo.GetSessionsByUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetUser, err)
	}

	var sessionInfos []*appDto.SessionInfo
	for _, session := range sessions {
		sessionInfos = append(sessionInfos, &appDto.SessionInfo{
			ID:           session.ID,
			IPAddress:    session.IPAddress,
			UserAgent:    session.UserAgent,
			CreatedAt:    session.CreatedAt,
			LastActivity: session.LastActivity,
			IsActive:     session.IsActive,
		})
	}

	return sessionInfos, nil
}

// DeleteSession implements AuthService.
func (as *authService) DeleteSession(ctx context.Context, sessionId int32) error {
	return as.authRepo.DeleteSession(ctx, sessionId)
}

// DeleteAllUserSessions implements AuthService.
func (as *authService) DeleteAllUserSessions(ctx context.Context, userId string) error {
	return as.authRepo.DeleteAllUserSessions(ctx, userId)
}

// IncrementFailedLoginAttempts implements AuthService.
func (as *authService) IncrementFailedLoginAttempts(ctx context.Context, userId string) error {
	return as.authRepo.IncrementFailedLoginAttempts(ctx, userId)
}

// ResetFailedLoginAttempts implements AuthService.
func (as *authService) ResetFailedLoginAttempts(ctx context.Context, userId string) error {
	return as.authRepo.ResetFailedLoginAttempts(ctx, userId)
}

// LockUserAccount implements AuthService.
func (as *authService) LockUserAccount(ctx context.Context, userId string) error {
	// Lock account for 30 minutes
	lockedUntil := utils.GetNowUnix() + (30 * 60)
	return as.authRepo.LockUserAccount(ctx, userId, constants.LOCKED, lockedUntil)
}

// CheckAccountLocked implements AuthService.
func (as *authService) CheckAccountLocked(ctx context.Context, userId string) (bool, error) {
	user, err := as.authRepo.GetUserById(ctx, userId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", msg.FailedToGetUser, err)
	}

	if user.LockedUntil > 0 && user.LockedUntil > utils.GetNowUnix() {
		return true, nil
	}

	return false, nil
}

// Validate2FACode validates a 2FA code
func (as *authService) Validate2FACode(ctx context.Context, userId, code, codeType string) (bool, error) {
	// Get 2FA code from database
	twoFACode, err := as.authRepo.Get2FACode(ctx, userId, code, codeType)
	if err != nil {
		return false, fmt.Errorf("%s: %w", msg.FailedToGet2FACode, err)
	}

	// Check if code is expired
	if twoFACode.ExpiresAt < utils.GetNowUnix() {
		return false, fmt.Errorf("%s", msg.TwoFACodeExpired)
	}

	// Delete the code after use
	as.authRepo.Delete2FACode(ctx, userId, codeType)

	return true, nil
}

// sendPasswordResetEmail sends a password reset email to the user
func (as *authService) sendPasswordResetEmail(email, resetToken string) error {
	// TODO: Implement actual email sending logic
	// For now, just log the reset token
	log.Printf("Password reset email would be sent to %s with token: %s", email, resetToken)
	return nil
}

func NewAuthService(
	authRepo authRepo.AuthRepository,
) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}
