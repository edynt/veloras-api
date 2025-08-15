package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
)

type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, org *entity.Organization) (string, error)
	GetOrganizationById(ctx context.Context, id string) (*entity.Organization, error)
	GetOrganizationsByParent(ctx context.Context, parentId string) ([]*entity.Organization, error)
	UpdateOrganization(ctx context.Context, org *entity.Organization) error
	DeleteOrganization(ctx context.Context, id string) error
	GetUserOrganization(ctx context.Context, userId string) (*entity.Organization, error)
	AssignUserToOrganization(ctx context.Context, userId string, organizationId string) error
	GetOrganizationUsers(ctx context.Context, organizationId string) ([]*entity.Account, error)
}
