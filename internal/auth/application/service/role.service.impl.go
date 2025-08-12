package service

import (
	"context"
	"fmt"

	"github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	roleRepo "github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
)

type roleService struct {
	roleRepo roleRepo.RoleRepository
}

// GetRoleById implements RoleService.
func (r *roleService) GetRoleById(ctx context.Context, id string) (dto.RoleOutPut, error) {
	exists, _ := r.roleRepo.GetRoleById(ctx, id)

	if exists == nil {
		return dto.RoleOutPut{}, fmt.Errorf(msg.RoleNotExists)
	}

	return dto.RoleOutPut{
		ID:          exists.ID,
		Name:        exists.Name,
		Description: exists.Description,
	}, nil
}

// GetRoles implements RoleService.
func (r *roleService) GetRoles(ctx context.Context) ([]dto.RoleOutPut, error) {
	roles, _ := r.roleRepo.GetRoles(ctx)

	if len(roles) == 0 {
		return nil, fmt.Errorf(msg.NoPermissionsFound)
	}

	var rolesOutPut []appDto.RoleOutPut
	if err := utils.SafeCopy(&rolesOutPut, &roles); err != nil {
		return nil, err
	}

	return rolesOutPut, nil
}

// DeleteRole implements RoleService.
func (r *roleService) DeleteRole(ctx context.Context, id string) error {
	exists, _ := r.roleRepo.GetRoleById(ctx, id)

	if exists == nil {
		return fmt.Errorf(msg.RoleNotExists)
	}

	err := r.roleRepo.DeleteRole(ctx, id)

	if err != nil {
		return fmt.Errorf("%s: %w", msg.CouldNotDeleteRole, err)
	}

	return nil
}

// UpdateRole implements RoleService.
func (r *roleService) UpdateRole(ctx context.Context, roleAppDTO dto.RoleAppDTO) (string, error) {
	exists, _ := r.roleRepo.GetRoleById(ctx, roleAppDTO.ID)

	if exists == nil {
		return "", fmt.Errorf(msg.RoleNotExists)
	}

	err := r.roleRepo.UpdateRole(ctx, &entity.Role{
		ID:          roleAppDTO.ID,
		Name:        roleAppDTO.Name,
		Description: roleAppDTO.Description,
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.CouldNotCreateRole, err)
	}

	return msg.Success, nil
}

// CreateRole implements RoleService.
func (r *roleService) CreateRole(ctx context.Context, roleAppDTO dto.RoleAppDTO) (string, error) {
	exists, _ := r.roleRepo.GetRoleByName(ctx, roleAppDTO.Name)

	if exists != nil {
		return "", fmt.Errorf(msg.RoleExists)
	}

	err := r.roleRepo.CreateRole(ctx, &entity.Role{
		Name:        roleAppDTO.Name,
		Description: roleAppDTO.Description,
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.CouldNotCreateRole, err)
	}

	return msg.Success, nil
}

// CreateUser implements RoleService.
func NewRoleService(
	roleRepo roleRepo.RoleRepository,
) RoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}
