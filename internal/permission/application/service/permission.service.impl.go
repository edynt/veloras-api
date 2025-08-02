package service

import (
	"context"
	"fmt"

	"github.com/edynnt/veloras-api/internal/permission/application/service/dto"
	"github.com/edynnt/veloras-api/internal/permission/domain/model/entity"
	permissionRepo "github.com/edynnt/veloras-api/internal/permission/domain/repository"
	"github.com/edynnt/veloras-api/pkg/response/msg"
)

type permissionService struct {
	permissionRepo permissionRepo.PermissisonRepository
}

// CreatePermission implements PermissionService.
func (p *permissionService) CreatePermission(ctx context.Context, permissionAppDto dto.PermissionAppDTO) (string, error) {

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
func (p *permissionService) GetPermissions(ctx context.Context, permissionAppDto dto.PermissionAppDTO) (string, error) {
	panic("unimplemented")
}

func NewPermissionService(
	permissionRepo permissionRepo.PermissisonRepository,
) PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
	}
}
