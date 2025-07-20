package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, account *entity.Account) (string, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	CreateVerificationCode(ctx context.Context, userVerification *entity.EmailVerification) (string, error)
}
