package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/permission/domain/model/entity"
)

type PermissisonRepository interface {
	GetPermissions(ctx context.Context) ([]*entity.Permission, error)
	GetPermissionById(ctx context.Context, id string) (*entity.Permission, error)
	GetPermissionByName(ctx context.Context, name string) (*entity.Permission, error)
	CreatePermission(ctx context.Context, permission *entity.Permission) error
	UpdatePermission(ctx context.Context, permission *entity.Permission) error
	DeletePermission(ctx context.Context, permission *entity.Permission) error
}
