package service

import (
	"context"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	authRepo "github.com/edynnt/veloras-api/internal/auth/domain/repository"
)

type authService struct {
	authRepo authRepo.AuthRepository
}

// Create implements AuthService.
func (a *authService) CreateUser(ctx context.Context, AccountAppDto appDto.AccountAppDTO) (int64, error) {
	panic("unimplemented")
}

func NewAuthService(
	authRepo authRepo.AuthRepository,
) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}
