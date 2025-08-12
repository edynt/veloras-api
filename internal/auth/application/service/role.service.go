package service

import (
	"context"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
)

type RoleService interface {
	CreateRole(ctx context.Context, roleAppDTO appDto.RoleAppDTO) (string, error)
	UpdateRole(ctx context.Context, roleAppDTO appDto.RoleAppDTO) (string, error)
	DeleteRole(ctx context.Context, id string) error
	GetRoles(ctx context.Context) ([]appDto.RoleOutPut, error)
	GetRoleById(ctx context.Context, id string) (appDto.RoleOutPut, error)
}
