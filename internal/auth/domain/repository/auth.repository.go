package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
)

type AuthRepository interface {
	CreatUser(ctx context.Context, account *entity.Account) (int64, error)
}
