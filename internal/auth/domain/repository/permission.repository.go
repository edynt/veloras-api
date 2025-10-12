package repository

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/auth/domain/model/entity"
)

type PermissisonRepository interface {
	GetPermissions(ctx context.Context) ([]*entity.Permission, error)
	GetPermissionById(ctx context.Context, id int) (*entity.Permission, error)
	GetPermissionByName(ctx context.Context, name string) (*entity.Permission, error)
	CreatePermission(ctx context.Context, permission *entity.Permission) error
	UpdatePermission(ctx context.Context, permission *entity.Permission) error
	DeletePermission(ctx context.Context, id int) error
}
