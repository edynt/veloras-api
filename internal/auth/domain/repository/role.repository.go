package repository

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/auth/domain/model/entity"
)

type RoleRepository interface {
	GetRoles(ctx context.Context) ([]*entity.Role, error)
	GetRoleById(ctx context.Context, id int) (*entity.Role, error)
	GetRoleByName(ctx context.Context, name string) (*entity.Role, error)
	CreateRole(ctx context.Context, Role *entity.Role) error
	UpdateRole(ctx context.Context, Role *entity.Role) error
	DeleteRole(ctx context.Context, id int) error
	AssignPermissions(ctx context.Context, rolePermission *entity.RolePermission) error
}
