package service

import (
	"context"
	"fmt"

	appDto "github.com/edynt/chogiare/veloras-api/internal/auth/application/service/dto"
	"github.com/edynt/chogiare/veloras-api/internal/auth/domain/model/entity"
	permissionRepo "github.com/edynt/chogiare/veloras-api/internal/auth/domain/repository"
	"github.com/edynt/chogiare/veloras-api/pkg/response/msg"
	"github.com/edynt/chogiare/veloras-api/pkg/utils"
)

type permissionService struct {
	permissionRepo permissionRepo.PermissisonRepository
}

// DeletePermission implements PermissionService.
func (p *permissionService) DeletePermission(ctx context.Context, id int) error {
	exists, _ := p.permissionRepo.GetPermissionById(ctx, id)

	if exists == nil {
		return fmt.Errorf(msg.PermissionNotFound)
	}

	err := p.permissionRepo.DeletePermission(ctx, id)

	if err != nil {
		return fmt.Errorf("%s: %w", msg.CouldNotDeletePermission, err)
	}

	return nil
}

// UpdatePermission implements PermissionService.
func (p *permissionService) UpdatePermission(ctx context.Context, permissionAppDto appDto.PermissionAppDTO) (int, error) {
	exists, _ := p.permissionRepo.GetPermissionById(ctx, permissionAppDto.ID)

	if exists == nil {
		return 0, fmt.Errorf(msg.PermissionNotFound)
	}

	err := p.permissionRepo.UpdatePermission(ctx, &entity.Permission{
		ID:          permissionAppDto.ID,
		Name:        permissionAppDto.Name,
		Description: permissionAppDto.Description,
	})

	if err != nil {
		return 0, fmt.Errorf("%s: %w", msg.CouldNotUpdatePermission, err)
	}

	return permissionAppDto.ID, nil
}

// CreatePermission implements PermissionService.
func (p *permissionService) CreatePermission(ctx context.Context, permissionAppDto appDto.PermissionAppDTO) (int, error) {

	exists, _ := p.permissionRepo.GetPermissionByName(ctx, permissionAppDto.Name)

	if exists != nil {
		return 0, fmt.Errorf(msg.PermissionExists)
	}

	err := p.permissionRepo.CreatePermission(ctx, &entity.Permission{
		Name:        permissionAppDto.Name,
		Description: permissionAppDto.Description,
	})

	if err != nil {
		return 0, fmt.Errorf("%s: %w", msg.CouldNotCreatePermission, err)
	}

	// Get the created permission to return its ID
	createdPermission, err := p.permissionRepo.GetPermissionByName(ctx, permissionAppDto.Name)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", msg.CouldNotCreatePermission, err)
	}

	return createdPermission.ID, nil
}

// GetPermissions implements PermissionService.
func (p *permissionService) GetPermissions(ctx context.Context) ([]appDto.PermissionOutPut, error) {
	permissions, _ := p.permissionRepo.GetPermissions(ctx)

	if len(permissions) == 0 {
		return nil, fmt.Errorf(msg.NoPermissionsFound)
	}

	var permissionsOutPut []appDto.PermissionOutPut
	if err := utils.SafeCopy(&permissionsOutPut, &permissions); err != nil {
		return nil, err
	}

	return permissionsOutPut, nil
}

func NewPermissionService(
	permissionRepo permissionRepo.PermissisonRepository,
) PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
	}
}
