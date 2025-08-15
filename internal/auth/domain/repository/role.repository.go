package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
)

type RoleRepository interface {
	GetRoles(ctx context.Context) ([]*entity.Role, error)
	GetRoleById(ctx context.Context, id string) (*entity.Role, error)
	GetRoleByName(ctx context.Context, name string) (*entity.Role, error)
	CreateRole(ctx context.Context, Role *entity.Role) error
	UpdateRole(ctx context.Context, Role *entity.Role) error
	DeleteRole(ctx context.Context, id string) error
	
	// Enhanced RBAC
	GetPermissionsByRole(ctx context.Context, roleId string) ([]*entity.Permission, error)
	AssignPermissionToRole(ctx context.Context, roleId, permissionId string) error
	RemovePermissionFromRole(ctx context.Context, roleId, permissionId string) error
	AssignRoleToUser(ctx context.Context, userId, roleId string) error
	RemoveRoleFromUser(ctx context.Context, userId, roleId string) error
	GetRolesByUser(ctx context.Context, userId string) ([]*entity.Role, error)
}
