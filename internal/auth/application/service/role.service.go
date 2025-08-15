package service

import (
	"context"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
)

type RoleService interface {
	// Basic Role Management
	CreateRole(ctx context.Context, roleAppDTO appDto.RoleAppDTO) (string, error)
	UpdateRole(ctx context.Context, roleAppDTO appDto.RoleAppDTO) (string, error)
	DeleteRole(ctx context.Context, id string) error
	GetRoles(ctx context.Context) ([]appDto.RoleOutPut, error)
	GetRoleById(ctx context.Context, id string) (appDto.RoleOutPut, error)
	
	// Enhanced RBAC
	GetPermissionsByRole(ctx context.Context, roleId string) ([]appDto.PermissionOutPut, error)
	AssignPermissionToRole(ctx context.Context, roleId, permissionId string) error
	RemovePermissionFromRole(ctx context.Context, roleId, permissionId string) error
	AssignRoleToUser(ctx context.Context, userId, roleId string) error
	RemoveRoleFromUser(ctx context.Context, userId, roleId string) error
	GetUserRoles(ctx context.Context, userId string) ([]appDto.RoleOutPut, error)
	GetRoleWithPermissions(ctx context.Context, roleId string) (*appDto.RoleWithPermissions, error)
}
