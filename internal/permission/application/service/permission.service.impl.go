package service

import (
	"context"
	"fmt"

	appDto "github.com/edynnt/veloras-api/internal/permission/application/service/dto"
	"github.com/edynnt/veloras-api/internal/permission/domain/model/entity"
	permissionRepo "github.com/edynnt/veloras-api/internal/permission/domain/repository"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
)

type permissionService struct {
	permissionRepo permissionRepo.PermissisonRepository
}

// DeletePermission implements PermissionService.
func (p *permissionService) DeletePermission(ctx context.Context, id string) error {
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
func (p *permissionService) UpdatePermission(ctx context.Context, permissionAppDto appDto.PermissionAppDTO) (string, error) {
	exists, _ := p.permissionRepo.GetPermissionById(ctx, permissionAppDto.ID)

	if exists == nil {
		return "", fmt.Errorf(msg.PermissionNotFound)
	}

	err := p.permissionRepo.UpdatePermission(ctx, &entity.Permission{
		ID:          permissionAppDto.ID,
		Name:        permissionAppDto.Name,
		Description: permissionAppDto.Description,
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.CouldNotUpdatePermission, err)
	}

	return msg.Success, nil
}

// CreatePermission implements PermissionService.
func (p *permissionService) CreatePermission(ctx context.Context, permissionAppDto appDto.PermissionAppDTO) (string, error) {

	exists, _ := p.permissionRepo.GetPermissionByName(ctx, permissionAppDto.Name)

	if exists != nil {
		return "", fmt.Errorf(msg.PermissionExists)
	}

	err := p.permissionRepo.CreatePermission(ctx, &entity.Permission{
		Name:        permissionAppDto.Name,
		Description: permissionAppDto.Description,
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.CouldNotCreatePermission, err)
	}

	return msg.Success, nil
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
