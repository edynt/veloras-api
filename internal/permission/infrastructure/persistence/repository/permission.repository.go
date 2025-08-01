package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/permission/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/permission/domain/repository"
	permissionSqlc "github.com/edynnt/veloras-api/internal/shared/gen"
	"github.com/jackc/pgx/v5/pgxpool"
)

type permissionRepository struct {
	db *permissionSqlc.Queries
}

// DeletePermission implements repository.PermissisonRepository.
func (p *permissionRepository) DeletePermission(ctx context.Context, permission *entity.Permission) error {
	panic("unimplemented")
}

// GetPermissionById implements repository.PermissisonRepository.
func (p *permissionRepository) GetPermissionById(ctx context.Context, id string) (*entity.Permission, error) {
	panic("unimplemented")
}

// UpdatePermission implements repository.PermissisonRepository.
func (p *permissionRepository) UpdatePermission(ctx context.Context, permission *entity.Permission) error {
	panic("unimplemented")
}

// CreatePermission implements repository.PermissisonRepository.
func (p *permissionRepository) CreatePermission(ctx context.Context, permission *entity.Permission) error {
	panic("unimplemented")
}

// GetPermissions implements repository.PermissisonRepository.
func (p *permissionRepository) GetPermissions(ctx context.Context) ([]*entity.Permission, error) {
	panic("unimplemented")
}

func NewPermissionRepository(db *pgxpool.Pool) repository.PermissisonRepository {
	queries := permissionSqlc.New(db) // db is *pgxpool.Pool
	return &permissionRepository{db: queries}
}
