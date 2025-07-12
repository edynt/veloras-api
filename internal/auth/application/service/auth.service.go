package service

import (
	"context"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
)

type AuthService interface {
	CreateUser(ctx context.Context, AcountInputDTO appDto.AccountAppDTO) (int64, error)
}
