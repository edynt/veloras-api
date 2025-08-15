package service

import (
	"context"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
)

type TwoFAService interface {
	Setup2FA(ctx context.Context, userId string) (*appDto.TwoFASetupResponse, error)
	Verify2FA(ctx context.Context, userId, code, codeType string) (*appDto.TwoFAVerifyResponse, error)
	Enable2FA(ctx context.Context, userId, secret, code string) error
	Disable2FA(ctx context.Context, userId, code string) error
	Get2FAStatus(ctx context.Context, userId string) (*appDto.TwoFAStatusResponse, error)
	Generate2FACode(ctx context.Context, userId, codeType string) (string, error)
	Validate2FACode(ctx context.Context, userId, code, codeType string) (bool, error)
}
