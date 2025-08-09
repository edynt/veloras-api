package service

import (
	"context"
	"fmt"

	"github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	roleRepo "github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/pkg/response/msg"
)

type roleService struct {
	roleRepo roleRepo.RoleRepository
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
