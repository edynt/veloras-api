package repository

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/auth/domain/repository"
	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	roleSqlc "github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/edynt/chogiare/veloras-api/pkg/utils"
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
func (r *roleRepository) DeleteRole(ctx context.Context, id int) error {
	return r.db.DeleteRole(ctx, int32(id))
}

// GetRoleById implements repository.RoleRepository.
func (r *roleRepository) GetRoleById(ctx context.Context, id int) (*entity.Role, error) {
	role, err := r.db.GetRoleById(ctx, int32(id))

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

// AssignPermissions implements repository.RoleRepository.
func (r *roleRepository) AssignPermissions(ctx context.Context, rolePermission *entity.RolePermission) error {
	var param gen.AssignPermissionToRoleParams

	if err := utils.SafeCopy(&param, rolePermission); err != nil {
		return err
	}

	return r.db.AssignPermissionToRole(ctx, param)
}

func NewRoleRepository(db *pgxpool.Pool) repository.RoleRepository {
	queries := roleSqlc.New(db) // db is *pgxpool.Pool
	return &roleRepository{db: queries}
}
