package repository

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/user/domain/model/entity"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context, limit, offset int32) ([]entity.User, error)
	CountAllUsers(ctx context.Context) (int64, error)
	GetUserByID(ctx context.Context, userID int32) (*entity.User, error)
	UpdateUserProfile(ctx context.Context, userID int32, username, phoneNumber, firstName, lastName, language string) (*entity.User, error)
}
