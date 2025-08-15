package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
)

type TwoFARepository interface {
	Create2FACode(ctx context.Context, code *entity.TwoFACode) (string, error)
	Get2FACode(ctx context.Context, userId string, code string, codeType string) (*entity.TwoFACode, error)
	Delete2FACode(ctx context.Context, userId string, codeType string) error
	DeleteExpired2FACodes(ctx context.Context, currentTime int64) error
	GetUser2FACodes(ctx context.Context, userId string, codeType string) ([]*entity.TwoFACode, error)
}
