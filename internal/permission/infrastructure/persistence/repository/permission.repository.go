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

// GetPermissions implements repository.PermissisonRepository.
func (p *permissionRepository) GetPermissions(ctx context.Context) (*entity.Permission, error) {
	panic("unimplemented")
}

func NewPermissionRepository(db *pgxpool.Pool) repository.PermissisonRepository {
	queries := permissionSqlc.New(db) // db is *pgxpool.Pool
	return &permissionRepository{db: queries}
}
