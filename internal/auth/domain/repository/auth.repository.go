package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, account *entity.Account) (int, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	CreateVerificationCode(ctx context.Context, userVerification *entity.EmailVerification) error
	GetVerificationCode(ctx context.Context, userId int, code int) (*entity.EmailVerification, error)
	UpdateUserStatus(ctx context.Context, userId int, status int) error
	GetUserByUsername(ctx context.Context, userName string) (*entity.Account, error)
	ActiveUser(ctx context.Context, userId int) error
	DeleteVerificationCode(ctx context.Context, userId int, code int) error
	SaveToken(ctx context.Context, token *entity.Session) error
	RefreshToken(ctx context.Context, refreshToken string) error
	DeleteSessionsByUser(ctx context.Context, userId int) error
}
