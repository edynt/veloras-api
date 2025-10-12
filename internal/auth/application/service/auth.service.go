package service

import (
	"context"

	appDto "github.com/edynt/chogiare/veloras-api/internal/auth/application/service/dto"
)

type AuthService interface {
	CreateUser(ctx context.Context, accountAppDTO appDto.AccountAppDTO) (int, error)
	VerifyUser(ctx context.Context, verificationEmailAppDTO appDto.EmailVerification) (bool, error)
	LoginUser(ctx context.Context, accountAppDTO appDto.AccountAppDTO) (appDto.UserOutPut, error)
	RefreshToken(ctx context.Context, refreshToken string) (appDto.TokenOut, error)
	Logout(ctx context.Context) error
	ChangePassword(ctx context.Context, userID int, currentPassword, newPassword string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}
