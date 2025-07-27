package service

import (
	"context"
	"fmt"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	authRepo "github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/pkg/constants"
	"github.com/edynnt/veloras-api/pkg/global"
	"github.com/edynnt/veloras-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	authRepo authRepo.AuthRepository
}

// LoginUser implements AuthService.
func (as *authService) LoginUser(ctx context.Context, accountAppDTO appDto.AccountAppDTO) (appDto.UserOutPut, error) {
	// 1. check exists
	user, err := as.authRepo.GetUserByUsername(ctx, accountAppDTO.Username)

	if err != nil {
		return appDto.UserOutPut{}, fmt.Errorf("Failed to check username exists: %w", err)
	}

	if user == nil {
		return appDto.UserOutPut{}, fmt.Errorf("Username not found")
	}

	// 2. check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(accountAppDTO.Password)); err != nil {
		return appDto.UserOutPut{}, fmt.Errorf("Invalid password")
	}

	// 3. check verified
	if !user.IsVerified {
		return appDto.UserOutPut{}, fmt.Errorf("User is not verified")
	}

	// 4. check status active
	if user.Status != constants.ACTIVE {
		return appDto.UserOutPut{}, fmt.Errorf("User is not actived")
	}

	// 5. Generate accessToken and refreshToken
	accessToken, err := utils.CreateAccessToken(user.ID)

	if err != nil {
		return appDto.UserOutPut{}, fmt.Errorf("Failed to create token: %w", err)
	}

	refreshToken, err := utils.CreateRefreshToken(user.ID)

	// 6. Save refresh token

	return appDto.UserOutPut{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// VerifyUser implements AuthService.
func (as *authService) VerifyUser(ctx context.Context, verificationEmailAppDTO appDto.EmailVerification) (bool, error) {
	existsVerificationCode, err := as.authRepo.GetVerificationCode(ctx, verificationEmailAppDTO.UserID, verificationEmailAppDTO.Code)

	if err != nil {
		return false, fmt.Errorf("failed to get verification code: %w", err)
	}

	now := utils.GetNowUnix()
	if existsVerificationCode.ExpiresAt < now {
		return false, fmt.Errorf("verification code expired")
	}

	err = as.authRepo.UpdateUserStatus(ctx, verificationEmailAppDTO.UserID, constants.ACTIVE)

	if err != nil {
		return false, fmt.Errorf("failed to update user status: %w", err)
	}

	err = as.authRepo.ActiveUser(ctx, verificationEmailAppDTO.UserID)
	if err != nil {
		return false, fmt.Errorf("failed to update user status: %w", err)
	}

	err = as.authRepo.DeleteVerificationCode(ctx, verificationEmailAppDTO.UserID, verificationEmailAppDTO.Code)
	if err != nil {
		return false, fmt.Errorf("failed to delete verification code: %w", err)
	}

	return true, nil
}

// Create implements AuthService.
func (as *authService) CreateUser(ctx context.Context, accountDto appDto.AccountAppDTO) (string, error) {
	//1. Check permissions -> event registered

	// 2. Check username exists
	exists, err := as.authRepo.UsernameExists(ctx, accountDto.Username)

	if err != nil {
		return "", fmt.Errorf("failed to check username exist: %w", err)
	}
	if exists {
		return "", fmt.Errorf("username already exists")
	}

	// 3. Check email exists
	exists, err = as.authRepo.EmailExists(ctx, accountDto.Email)
	if err != nil {
		return "", fmt.Errorf("failed to check email exist: %w", err)
	}

	if exists {
		return "", fmt.Errorf("email already exists")
	}

	// 4. GenerateFromPassword
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(accountDto.Password), bcrypt.DefaultCost)
	if err != nil {
		// log.Printf("Error hashing password for user %s: %v", accountDto.Username, err)
		return "", fmt.Errorf("failed to secure password: %w", err) // Không lộ chi tiết lỗi hash
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
		// log.Printf("Error creating user %s in database: %v", accountToCreate.Username, err)
		return "", fmt.Errorf("could not create account: %w", err)
	}

	if newAccountId == "" {
		// log.Printf("Repository returned invalid account ID %d for user %s", newAccountId, accountToCreate.Username)
		return "", fmt.Errorf("account creation failed to return a valid ID")
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

func NewAuthService(
	authRepo authRepo.AuthRepository,
) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}
