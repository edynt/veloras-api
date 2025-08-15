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

// GetPermissionsByRole implements RoleService.
func (r *roleService) GetPermissionsByRole(ctx context.Context, roleId string) ([]appDto.PermissionOutPut, error) {
	permissions, err := r.roleRepo.GetPermissionsByRole(ctx, roleId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetPermissionsByRole, err)
	}

	var permissionsOutPut []appDto.PermissionOutPut
	if err := utils.SafeCopy(&permissionsOutPut, &permissions); err != nil {
		return nil, err
	}

	return permissionsOutPut, nil
}

// AssignPermissionToRole implements RoleService.
func (r *roleService) AssignPermissionToRole(ctx context.Context, roleId, permissionId string) error {
	err := r.roleRepo.AssignPermissionToRole(ctx, roleId, permissionId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToAssignPermissionToRole, err)
	}

	return nil
}

// RemovePermissionFromRole implements RoleService.
func (r *roleService) RemovePermissionFromRole(ctx context.Context, roleId, permissionId string) error {
	// TODO: Implement this method
	return fmt.Errorf("not implemented")
}

// AssignRoleToUser implements RoleService.
func (r *roleService) AssignRoleToUser(ctx context.Context, userId, roleId string) error {
	err := r.roleRepo.AssignRoleToUser(ctx, userId, roleId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToAssignRoleToUser, err)
	}

	return nil
}

// RemoveRoleFromUser implements RoleService.
func (r *roleService) RemoveRoleFromUser(ctx context.Context, userId, roleId string) error {
	// TODO: Implement this method
	return fmt.Errorf("not implemented")
}

// GetUserRoles implements RoleService.
func (r *roleService) GetUserRoles(ctx context.Context, userId string) ([]appDto.RoleOutPut, error) {
	roles, err := r.roleRepo.GetRolesByUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetRolesByUser, err)
	}

	var rolesOutPut []appDto.RoleOutPut
	if err := utils.SafeCopy(&rolesOutPut, &roles); err != nil {
		return nil, err
	}

	return rolesOutPut, nil
}

// GetRoleWithPermissions implements RoleService.
func (r *roleService) GetRoleWithPermissions(ctx context.Context, roleId string) (*appDto.RoleWithPermissions, error) {
	// Get role
	role, err := r.GetRoleById(ctx, roleId)
	if err != nil {
		return nil, err
	}

	// Get permissions for role
	permissions, err := r.GetPermissionsByRole(ctx, roleId)
	if err != nil {
		return nil, err
	}

	return &appDto.RoleWithPermissions{
		Role:        role,
		Permissions: permissions,
	}, nil
}

// CreateUser implements RoleService.
func NewRoleService(
	roleRepo roleRepo.RoleRepository,
) RoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}
