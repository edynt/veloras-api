package service

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/user/application/service/dto"
)

type UserService interface {
	GetAllUsers(ctx context.Context, page, pageSize int32) (*dto.UserListAppDTO, error)
	GetUserByID(ctx context.Context, userID int32) (*dto.UserAppDTO, error)
	GetCurrentUser(ctx context.Context) (*dto.UserAppDTO, error)
	UpdateUserProfile(ctx context.Context, userID int32, req dto.UpdateUserProfileAppDTO) (*dto.UserAppDTO, error)
}
