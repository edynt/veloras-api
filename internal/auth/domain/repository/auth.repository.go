package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, account *entity.Account) (string, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	CreateVerificationCode(ctx context.Context, userVerification *entity.EmailVerification) error
	GetVerificationCode(ctx context.Context, userId string, code int) (*entity.EmailVerification, error)
	UpdateUserStatus(ctx context.Context, userId string, status int) error
	GetUserByUsername(ctx context.Context, userName string) (*entity.Account, error)
	ActiveUser(ctx context.Context, userId string) error
	DeleteVerificationCode(ctx context.Context, userId string, code int) error
	SaveToken(ctx context.Context, token *entity.Session) error
}
