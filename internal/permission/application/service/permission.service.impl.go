package service

import (
	"context"

	"github.com/edynnt/veloras-api/internal/permission/application/service/dto"
	permissionRepo "github.com/edynnt/veloras-api/internal/permission/domain/repository"
)

type permissionService struct {
	permissionRepo permissionRepo.PermissisonRepository
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
