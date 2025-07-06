package service

import (
	"context"

	inputDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
)

type AuthService interface {
	CreateUser(ctx context.Context, AcountInputDTO inputDto.AccountInputDTO) (int64, error)
}
