package service

import (
	"context"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
)

type RoleService interface {
	CreateRole(ctx context.Context, roleAppDTO appDto.RoleAppDTO) (int, error)
	UpdateRole(ctx context.Context, roleAppDTO appDto.RoleAppDTO) (int, error)
	DeleteRole(ctx context.Context, id int) error
	GetRoles(ctx context.Context) ([]appDto.RoleOutPut, error)
	GetRoleById(ctx context.Context, id int) (appDto.RoleOutPut, error)
	AssignPermissions(ctx context.Context, rolePermissionAppDto appDto.RolePermissionAppDTO) error
}
