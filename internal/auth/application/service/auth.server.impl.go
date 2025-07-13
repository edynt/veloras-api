package service

import (
	"context"
	"fmt"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	authRepo "github.com/edynnt/veloras-api/internal/auth/domain/repository"
)

type authService struct {
	authRepo authRepo.AuthRepository
}

// Create implements AuthService.
func (as *authService) CreateUser(ctx context.Context, AccountAppDto appDto.AccountAppDTO) (int64, error) {
	fmt.Println("call create user auth.service.impl")

	//1. Check permissions -> event registered

	// 2. Check username exist
	_, err := as.authRepo.UsernameExists(ctx, AccountAppDto.Username)

	fmt.Println(err)

	return 1, nil
}

func NewAuthService(
	authRepo authRepo.AuthRepository,
) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}
