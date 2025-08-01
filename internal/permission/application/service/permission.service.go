package service

import (
	"context"

	appDto "github.com/edynnt/veloras-api/internal/permission/application/service/dto"
)

type PermissionService interface {
	GetPermissions(ctx context.Context, permissionAppDto appDto.PermissionAppDTO) (string, error)
	CreatePermission(ctx context.Context, permissionAppDto appDto.PermissionAppDTO) (string, error)
}
