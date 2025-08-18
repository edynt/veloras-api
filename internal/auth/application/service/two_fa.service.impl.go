package service

import (
	"context"
	"fmt"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	twoFARepo "github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
)

type twoFAService struct {
	twoFARepo twoFARepo.TwoFARepository
	authRepo  twoFARepo.AuthRepository
}

// Setup2FA implements TwoFAService.
func (t *twoFAService) Setup2FA(ctx context.Context, userId string) (*appDto.TwoFASetupResponse, error) {
	// Generate 2FA secret
	secret := utils.GenerateRandomString(32)
	
	// Generate QR code (placeholder for now)
	qrCode := fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=Veloras", userId, secret)

	return &appDto.TwoFASetupResponse{
		Secret: secret,
		QRCode: qrCode,
	}, nil
}

// Verify2FA implements TwoFAService.
func (t *twoFAService) Verify2FA(ctx context.Context, userId, code, codeType string) (*appDto.TwoFAVerifyResponse, error) {
	valid, err := t.Validate2FACode(ctx, userId, code, codeType)
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

// Enable2FA implements TwoFAService.
func (t *twoFAService) Enable2FA(ctx context.Context, userId, secret, code string) error {
	// Verify the code first
	valid, err := t.Validate2FACode(ctx, userId, code, "setup")
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf(msg.TwoFACodeInvalid)
	}

	// Update user's 2FA enabled status
	err = t.authRepo.UpdateUser2FAEnabled(ctx, userId, true)
	if err != nil {
		return fmt.Errorf("failed to enable 2FA: %w", err)
	}

	return nil
}

// Disable2FA implements TwoFAService.
func (t *twoFAService) Disable2FA(ctx context.Context, userId, code string) error {
	// Verify the code first
	valid, err := t.Validate2FACode(ctx, userId, code, "setup")
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf(msg.TwoFACodeInvalid)
	}

	// Update user's 2FA enabled status
	err = t.authRepo.UpdateUser2FAEnabled(ctx, userId, false)
	if err != nil {
		return fmt.Errorf("failed to disable 2FA: %w", err)
	}

	return nil
}

// Get2FAStatus implements TwoFAService.
func (t *twoFAService) Get2FAStatus(ctx context.Context, userId string) (*appDto.TwoFAStatusResponse, error) {
	// Get user's 2FA status from the auth service
	user, err := t.authRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &appDto.TwoFAStatusResponse{
		Enabled: user.TwoFAEnabled,
		Secret:  user.TwoFASecret,
	}, nil
}

// Generate2FACode implements TwoFAService.
func (t *twoFAService) Generate2FACode(ctx context.Context, userId, codeType string) (string, error) {
	// Generate a 6-digit code
	codeInt := utils.GenerateSixDigitCode()
	code := fmt.Sprintf("%06d", codeInt) // Convert to 6-digit string with leading zeros
	
	// Set expiration time (5 minutes for most codes)
	expiresAt := utils.GetNowUnix() + (5 * 60)
	
	// Create 2FA code entity
	twoFACode := &entity.TwoFACode{
		UserID:    userId,
		Code:      code,
		ExpiresAt: expiresAt,
		Type:      codeType,
	}
	
	// Save to database
	_, err := t.twoFARepo.Create2FACode(ctx, twoFACode)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.FailedToCreate2FACode, err)
	}
	
	return code, nil
}

// Validate2FACode implements TwoFAService.
func (t *twoFAService) Validate2FACode(ctx context.Context, userId, code, codeType string) (bool, error) {
	// Get 2FA code from database
	twoFACode, err := t.twoFARepo.Get2FACode(ctx, userId, code, codeType)
	if err != nil {
		return false, fmt.Errorf("%s: %w", msg.FailedToGet2FACode, err)
	}

	// Check if code is expired
	if twoFACode.ExpiresAt < utils.GetNowUnix() {
		return false, fmt.Errorf("%s", msg.TwoFACodeExpired)
	}

	// Delete the code after use
	t.twoFARepo.Delete2FACode(ctx, userId, codeType)

	return true, nil
}

func NewTwoFAService(twoFARepo twoFARepo.TwoFARepository, authRepo twoFARepo.AuthRepository) TwoFAService {
	return &twoFAService{
		twoFARepo: twoFARepo,
		authRepo:  authRepo,
	}
}
