package service

import (
	"context"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
)

type OrganizationService interface {
	CreateOrganization(ctx context.Context, req appDto.CreateOrganizationRequest, createdBy string) (string, error)
	GetOrganizationById(ctx context.Context, id string) (*appDto.OrganizationResponse, error)
	GetOrganizationsByParent(ctx context.Context, parentId string) ([]*appDto.OrganizationResponse, error)
	UpdateOrganization(ctx context.Context, req appDto.UpdateOrganizationRequest) error
	DeleteOrganization(ctx context.Context, id string) error
	GetUserOrganization(ctx context.Context, userId string) (*appDto.OrganizationResponse, error)
	AssignUserToOrganization(ctx context.Context, userId, organizationId string) error
	GetOrganizationUsers(ctx context.Context, organizationId string) ([]*appDto.UserResponse, error)
}
