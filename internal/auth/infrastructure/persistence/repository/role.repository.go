package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/internal/shared/gen"
	roleSqlc "github.com/edynnt/veloras-api/internal/shared/gen"
	"github.com/edynnt/veloras-api/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type roleRepository struct {
	db *roleSqlc.Queries
}

// CreateRole implements repository.RoleRepository.
func (r *roleRepository) CreateRole(ctx context.Context, Role *entity.Role) error {
	var param gen.CreateRoleParams
	if err := utils.SafeCopy(&param, Role); err != nil {
		return err
	}

	err := r.db.CreateRole(ctx, param)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRole implements repository.RoleRepository.
func (r *roleRepository) DeleteRole(ctx context.Context, id string) error {
	convertId, err := utils.ConvertUUID(id)

	if err != nil {
		return err
	}

	return r.db.DeleteRole(ctx, convertId)
}

// GetRoleById implements repository.RoleRepository.
func (r *roleRepository) GetRoleById(ctx context.Context, id string) (*entity.Role, error) {
	convertId, err := utils.ConvertUUID(id)

	if err != nil {
		return nil, err
	}

	role, err := r.db.GetRoleById(ctx, convertId)

	if err != nil {
		return nil, err
	}

	var entityResult entity.Role
	if err := utils.SafeCopy(&entityResult, &role); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// GetRoleByName implements repository.RoleRepository.
func (r *roleRepository) GetRoleByName(ctx context.Context, name string) (*entity.Role, error) {
	res, err := r.db.GetRoleByName(ctx, name)

	if err != nil {
		return nil, err
	}

	var entityResult entity.Role
	if err := utils.SafeCopy(&entityResult, &res); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// GetRoles implements repository.RoleRepository.
func (r *roleRepository) GetRoles(ctx context.Context) ([]*entity.Role, error) {
	roles, err := r.db.GetRoles(ctx)
	if err != nil {
		return nil, err
	}

	var entityResult []*entity.Role
	if err := utils.SafeCopy(&entityResult, &roles); err != nil {
		return nil, err
	}

	return entityResult, nil
}

// UpdateRole implements repository.RoleRepository.
func (r *roleRepository) UpdateRole(ctx context.Context, Role *entity.Role) error {
	var param gen.UpdateRoleParams
	if err := utils.SafeCopy(&param, Role); err != nil {
		return err
	}

	err := r.db.UpdateRole(ctx, param)
	if err != nil {
		return err
	}

	return nil
}

func NewRoleRepository(db *pgxpool.Pool) repository.RoleRepository {
	queries := roleSqlc.New(db) // db is *pgxpool.Pool
	return &roleRepository{db: queries}
}
