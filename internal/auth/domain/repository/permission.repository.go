package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
)

type PermissisonRepository interface {
	GetPermissions(ctx context.Context) ([]*entity.Permission, error)
	GetPermissionById(ctx context.Context, id string) (*entity.Permission, error)
	GetPermissionByName(ctx context.Context, name string) (*entity.Permission, error)
	CreatePermission(ctx context.Context, permission *entity.Permission) error
	UpdatePermission(ctx context.Context, permission *entity.Permission) error
	DeletePermission(ctx context.Context, id string) error
	
	// Resource-based permissions
	GetPermissionsByResourceType(ctx context.Context, resourceType string) ([]*entity.Permission, error)
	GetPermissionsByOrganization(ctx context.Context, organizationId string) ([]*entity.Permission, error)
	CheckUserPermission(ctx context.Context, userId, resourceType, resourceAction string) (bool, error)
	GetUserPermissions(ctx context.Context, userId string) ([]*entity.Permission, error)
}
