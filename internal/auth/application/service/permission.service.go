package service

import (
	"context"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
)

type PermissionService interface {
	// Basic Permission Management
	GetPermissions(ctx context.Context) ([]appDto.PermissionOutPut, error)
	CreatePermission(ctx context.Context, permissionAppDto appDto.PermissionAppDTO) (string, error)
	UpdatePermission(ctx context.Context, permissionAppDto appDto.PermissionAppDTO) (string, error)
	DeletePermission(ctx context.Context, id string) error
	
	// Resource-Based Permissions
	GetPermissionsByResourceType(ctx context.Context, resourceType string) ([]appDto.PermissionOutPut, error)
	GetPermissionsByOrganization(ctx context.Context, organizationId string) ([]appDto.PermissionOutPut, error)
	CheckUserPermission(ctx context.Context, userId, resourceType, resourceAction string) (bool, error)
	GetUserPermissions(ctx context.Context, userId string) ([]appDto.PermissionOutPut, error)
}
