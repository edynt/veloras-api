package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/internal/shared/gen"
	"github.com/edynnt/veloras-api/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type permissionRepository struct {
	db *gen.Queries
}

// GetPermissionByName implements repository.PermissisonRepository.
func (p *permissionRepository) GetPermissionByName(ctx context.Context, name string) (*entity.Permission, error) {
	res, err := p.db.GetPermissionByName(ctx, name)

	if err != nil {
		return nil, err
	}

	var entityResult entity.Permission
	if err := utils.SafeCopy(&entityResult, &res); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// DeletePermission implements repository.PermissisonRepository.
func (p *permissionRepository) DeletePermission(ctx context.Context, id string) error {
	convertId, err := utils.ConvertUUID(id)

	if err != nil {
		return err
	}

	return p.db.DeletePermission(ctx, convertId)
}

// GetPermissionById implements repository.PermissisonRepository.
func (p *permissionRepository) GetPermissionById(ctx context.Context, id string) (*entity.Permission, error) {
	convertId, err := utils.ConvertUUID(id)

	if err != nil {
		return nil, err
	}
	permission, err := p.db.GetPermissionById(ctx, convertId)

	if err != nil {
		return nil, err
	}

	var entityResult entity.Permission
	if err := utils.SafeCopy(&entityResult, &permission); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// UpdatePermission implements repository.PermissisonRepository.
func (p *permissionRepository) UpdatePermission(ctx context.Context, permission *entity.Permission) error {
	var param gen.UpdatePermissionParams
	if err := utils.SafeCopy(&param, permission); err != nil {
		return err
	}

	err := p.db.UpdatePermission(ctx, param)
	if err != nil {
		return err
	}

	return nil
}

// CreatePermission implements repository.PermissisonRepository.
func (p *permissionRepository) CreatePermission(ctx context.Context, permission *entity.Permission) error {
	var param gen.CreatePermissionParams
	if err := utils.SafeCopy(&param, permission); err != nil {
		return err
	}

	err := p.db.CreatePermission(ctx, param)
	if err != nil {
		return err
	}

	return nil
}

// GetPermissions implements repository.PermissisonRepository.
func (p *permissionRepository) GetPermissions(ctx context.Context) ([]*entity.Permission, error) {
	permissions, err := p.db.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	var entityResult []*entity.Permission
	if err := utils.SafeCopy(&entityResult, &permissions); err != nil {
		return nil, err
	}

	return entityResult, nil
}

// GetPermissionsByResourceType implements repository.PermissisonRepository.
func (p *permissionRepository) GetPermissionsByResourceType(ctx context.Context, resourceType string) ([]*entity.Permission, error) {
	param := pgtype.Text{String: resourceType, Valid: true}
	permissions, err := p.db.GetPermissionsByResourceType(ctx, param)
	if err != nil {
		return nil, err
	}

	var entityResult []*entity.Permission
	if err := utils.SafeCopy(&entityResult, &permissions); err != nil {
		return nil, err
	}

	return entityResult, nil
}

// GetPermissionsByOrganization implements repository.PermissisonRepository.
func (p *permissionRepository) GetPermissionsByOrganization(ctx context.Context, organizationId string) ([]*entity.Permission, error) {
	convertId, err := utils.ConvertUUID(organizationId)
	if err != nil {
		return nil, err
	}

	permissions, err := p.db.GetPermissionsByOrganization(ctx, convertId)
	if err != nil {
		return nil, err
	}

	var entityResult []*entity.Permission
	if err := utils.SafeCopy(&entityResult, &permissions); err != nil {
		return nil, err
	}

	return entityResult, nil
}

// CheckUserPermission implements repository.PermissisonRepository.
func (p *permissionRepository) CheckUserPermission(ctx context.Context, userId, resourceType, resourceAction string) (bool, error) {
	convertUserId, err := utils.ConvertUUID(userId)
	if err != nil {
		return false, err
	}

	param := gen.CheckUserPermissionParams{
		UserID:         convertUserId,
		ResourceType:   pgtype.Text{String: resourceType, Valid: true},
		ResourceAction: pgtype.Text{String: resourceAction, Valid: true},
	}

	hasPermission, err := p.db.CheckUserPermission(ctx, param)
	if err != nil {
		return false, err
	}

	return hasPermission, nil
}

// GetUserPermissions implements repository.PermissisonRepository.
func (p *permissionRepository) GetUserPermissions(ctx context.Context, userId string) ([]*entity.Permission, error) {
	convertUserId, err := utils.ConvertUUID(userId)
	if err != nil {
		return nil, err
	}

	permissions, err := p.db.GetUserPermissions(ctx, convertUserId)
	if err != nil {
		return nil, err
	}

	var entityResult []*entity.Permission
	if err := utils.SafeCopy(&entityResult, &permissions); err != nil {
		return nil, err
	}

	return entityResult, nil
}

func NewPermissionRepository(db *pgxpool.Pool) repository.PermissisonRepository {
	queries := gen.New(db) // db is *pgxpool.Pool
	return &permissionRepository{db: queries}
}
