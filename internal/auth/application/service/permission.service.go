package service

import (
	"context"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
)

type PermissionService interface {
	GetPermissions(ctx context.Context) ([]appDto.PermissionOutPut, error)
	CreatePermission(ctx context.Context, permissionAppDto appDto.PermissionAppDTO) (string, error)
	UpdatePermission(ctx context.Context, permissionAppDto appDto.PermissionAppDTO) (string, error)
	DeletePermission(ctx context.Context, id string) error
}
