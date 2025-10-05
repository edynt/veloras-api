package service

import (
	"context"
	"fmt"

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
func (as *authService) LoginUser(ctx context.Context, accountAppDTO appDto.AccountAppDTO) (appDto.UserOutPut, error) {
	// 1. check exists by email
	user, err := as.authRepo.GetUserByEmail(ctx, accountAppDTO.Email)

	if err != nil {
		return appDto.UserOutPut{}, fmt.Errorf("%s: %w", msg.FailedToCheckEmailExists, err)
	}

	if user == nil {
		return appDto.UserOutPut{}, fmt.Errorf(msg.EmailNotFound)
	}

	// 2. check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(accountAppDTO.Password)); err != nil {
		return appDto.UserOutPut{}, fmt.Errorf("Invalid password")
	}

	// 3. check verified
	if !user.IsVerified {
		return appDto.UserOutPut{}, fmt.Errorf(msg.UserIsNotVerified)
	}

	// 4. check status active
	if user.Status != constants.ACTIVE {
		return appDto.UserOutPut{}, fmt.Errorf(msg.UserIsNotActive)
	}

	// 5. Generate accessToken and refreshToken
	accessToken, err := utils.CreateToken(user.ID, false)

	if err != nil {
		return appDto.UserOutPut{}, fmt.Errorf("%s: %w", msg.FailedToCreateToken, err)
	}

	refreshToken, err := utils.CreateToken(user.ID, true)

	if err != nil {
		return appDto.UserOutPut{}, fmt.Errorf("%s: %w", msg.FailedToCreateToken, err)
	}

	// 6. Save refresh token
	tokenExpiresAt := utils.AddDays(global.Config.JWT.RefreshTokenExpire)

	err = as.authRepo.SaveToken(ctx, &entity.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    tokenExpiresAt,
	})

	if err != nil {
		return appDto.UserOutPut{}, fmt.Errorf("%s: %w", msg.FailedToSaveToken, err)
	}

	return appDto.UserOutPut{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		PhoneNumber:    user.PhoneNumber,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		IsVerified:     user.IsVerified,
		Status:         user.Status,
		Language:       user.Language,
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		TokenExpiresAt: tokenExpiresAt,
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

// Create implements AuthService.
func (as *authService) CreateUser(ctx context.Context, accountDto appDto.AccountAppDTO) (int, error) {
	//1. Check permissions -> event registered

	// 2. Check username exists
	exists, err := as.authRepo.UsernameExists(ctx, accountDto.Username)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", msg.FailedToCheckUserNameExists, err)
	}
	if exists {
		return 0, fmt.Errorf(msg.UsernameExists)
	}

	// 3. Check email exists
	exists, err = as.authRepo.EmailExists(ctx, accountDto.Email)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", msg.FailedToCheckEmailExists, err)
	}

	if exists {
		return 0, fmt.Errorf(msg.EmailExists)
	}

	// 4. GenerateFromPassword
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(accountDto.Password), bcrypt.DefaultCost)
	if err != nil {
		// log.Printf("Error hashing password for user %s: %v", accountDto.Username, err)
		return 0, fmt.Errorf("%s: %w", msg.FailedToSecurePassword, err) // Không lộ chi tiết lỗi hash
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
		return 0, fmt.Errorf("%s: %w", msg.CouldNotCreateAccount, err)
	}

	if newAccountId == 0 {
		return 0, fmt.Errorf(msg.CouldNotCreateAccount)
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

// Refresh token
func (as *authService) RefreshToken(ctx context.Context, refreshToken string) (appDto.TokenOut, error) {
	// 1. Verify refresh token (signature, expiry)
	claims, err := utils.VerifyTokenSubject(refreshToken)
	if err != nil {
		return appDto.TokenOut{}, fmt.Errorf("%s: %w", msg.InvalidRefreshToken, err)
	}

	// 2. Optional: check token not expired (VerifyTokenSubject already does Valid())

	// 3. Parse user id from subject
	userID := utils.StringToInt(claims.Subject)
	if userID == 0 {
		return appDto.TokenOut{}, fmt.Errorf(msg.InvalidRefreshToken)
	}

	// 4. Optionally, validate the refresh token exists in sessions store
	if err := as.authRepo.RefreshToken(ctx, refreshToken); err != nil {
		return appDto.TokenOut{}, fmt.Errorf("%s: %w", msg.FailedToRefreshToken, err)
	}

	// 5. Generate new access token and new refresh token
	accessToken, err := utils.CreateToken(userID, false)
	if err != nil {
		return appDto.TokenOut{}, fmt.Errorf("%s: %w", msg.FailedToCreateToken, err)
	}

	newRefreshToken, err := utils.CreateToken(userID, true)
	if err != nil {
		return appDto.TokenOut{}, fmt.Errorf("%s: %w", msg.FailedToCreateToken, err)
	}

	// 6. Save new refresh token to database
	tokenExpiresAt := utils.AddDays(global.Config.JWT.RefreshTokenExpire)
	err = as.authRepo.SaveToken(ctx, &entity.Session{
		UserID:       userID,
		RefreshToken: newRefreshToken,
		ExpiresAt:    tokenExpiresAt,
	})
	if err != nil {
		return appDto.TokenOut{}, fmt.Errorf("%s: %w", msg.FailedToSaveToken, err)
	}

	return appDto.TokenOut{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func NewAuthService(
	authRepo authRepo.AuthRepository,
) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}

// Logout implements AuthService.
func (as *authService) Logout(ctx context.Context) error {
	subject := ctx.Value("subjectID")
	if subject == nil {
		return fmt.Errorf(msg.Unauthorized)
	}

	userID := utils.StringToInt(subject.(string))
	if userID == 0 {
		return fmt.Errorf(msg.UserIdInvalid)
	}

	if err := as.authRepo.DeleteSessionsByUser(ctx, userID); err != nil {
		return fmt.Errorf("%s: %w", msg.CouldNotDeleteUser, err)
	}
	return nil
}

// ChangePassword implements AuthService.
func (as *authService) ChangePassword(ctx context.Context, userID int, currentPassword, newPassword string) error {
	// 1. Get user by ID
	user, err := as.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToGetUserById, err)
	}

	if user == nil {
		return fmt.Errorf(msg.UserIdInvalid)
	}

	// 2. Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return fmt.Errorf(msg.CurrentPasswordIncorrect)
	}

	// 3. Hash new password
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToSecurePassword, err)
	}
	hashedPassword := string(hashedPasswordBytes)

	// 4. Update password in database
	if err := as.authRepo.UpdateUserPassword(ctx, userID, hashedPassword); err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdatePassword, err)
	}

	return nil
}

// ForgotPassword implements AuthService.
func (as *authService) ForgotPassword(ctx context.Context, email string) error {
	// 1. Check if user exists with this email
	user, err := as.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToGetUserById, err)
	}

	if user == nil {
		// For security reasons, don't reveal if email exists or not
		// Just return success message
		return nil
	}

	// 2. Generate reset token
	resetToken := utils.GenerateRandomString(32)

	// 3. Set token expiration (1 hour from now)
	tokenExpiresAt := utils.AddHours(1)

	// 4. Save reset token to database
	err = as.authRepo.CreatePasswordReset(ctx, &entity.PasswordReset{
		UserID:     user.ID,
		ResetToken: resetToken,
		ExpiresAt:  tokenExpiresAt,
	})

	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToCreatePasswordReset, err)
	}

	// 5. Send email with reset link
	go utils.SendTemplateEmailOtp(
		[]string{email}, global.Config.SMTP.User,
		"reset-password.html",
		map[string]interface{}{
			"ResetToken": resetToken,
			"Username":   user.Username,
			"ExpiresAt":  tokenExpiresAt,
		},
	)

	return nil
}

// ResetPassword implements AuthService.
func (as *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// 1. Find password reset record by token
	passwordReset, err := as.authRepo.GetPasswordResetByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.InvalidResetToken, err)
	}

	if passwordReset == nil {
		return fmt.Errorf(msg.InvalidResetToken)
	}

	// 2. Check if token is expired
	now := utils.GetNowUnix()
	if passwordReset.ExpiresAt < now {
		return fmt.Errorf(msg.ResetTokenExpired)
	}

	// 3. Hash new password
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToSecurePassword, err)
	}
	hashedPassword := string(hashedPasswordBytes)

	// 4. Update user password
	if err := as.authRepo.UpdateUserPassword(ctx, passwordReset.UserID, hashedPassword); err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdatePassword, err)
	}

	// 5. Delete the reset token
	if err := as.authRepo.DeletePasswordReset(ctx, passwordReset.UserID, token); err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeletePasswordReset, err)
	}

	return nil
}
