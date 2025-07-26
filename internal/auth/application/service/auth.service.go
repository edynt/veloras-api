package service

import (
	"context"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
)

type AuthService interface {
	CreateUser(ctx context.Context, accountAppDTO appDto.AccountAppDTO) (string, error)
	VerifyUser(ctx context.Context, verificationEmailAppDTO appDto.EmailVerification) (bool, error)
	LoginUser(ctx context.Context, accountAppDTO appDto.AccountAppDTO) (appDto.UserOutPut, error)
}
