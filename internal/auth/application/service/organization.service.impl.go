package service

import (
	"context"
	"fmt"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	orgRepo "github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/pkg/response/msg"
)

type organizationService struct {
	orgRepo orgRepo.OrganizationRepository
}

// CreateOrganization implements OrganizationService.
func (o *organizationService) CreateOrganization(ctx context.Context, req appDto.CreateOrganizationRequest, createdBy string) (string, error) {
	org := &entity.Organization{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		CreatedBy:   createdBy,
	}

	orgId, err := o.orgRepo.CreateOrganization(ctx, org)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.FailedToCreateOrganization, err)
	}

	return orgId, nil
}

// GetOrganizationById implements OrganizationService.
func (o *organizationService) GetOrganizationById(ctx context.Context, id string) (*appDto.OrganizationResponse, error) {
	org, err := o.orgRepo.GetOrganizationById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetOrganization, err)
	}

	return &appDto.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description,
		ParentID:    org.ParentID,
		CreatedBy:   org.CreatedBy,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}, nil
}

// GetOrganizationsByParent implements OrganizationService.
func (o *organizationService) GetOrganizationsByParent(ctx context.Context, parentId string) ([]*appDto.OrganizationResponse, error) {
	orgs, err := o.orgRepo.GetOrganizationsByParent(ctx, parentId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetOrganization, err)
	}

	var responses []*appDto.OrganizationResponse
	for _, org := range orgs {
		responses = append(responses, &appDto.OrganizationResponse{
			ID:          org.ID,
			Name:        org.Name,
			Description: org.Description,
			ParentID:    org.ParentID,
			CreatedBy:   org.CreatedBy,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
		})
	}

	return responses, nil
}

// UpdateOrganization implements OrganizationService.
func (o *organizationService) UpdateOrganization(ctx context.Context, req appDto.UpdateOrganizationRequest) error {
	org := &entity.Organization{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	}

	err := o.orgRepo.UpdateOrganization(ctx, org)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdateOrganization, err)
	}

	return nil
}

// DeleteOrganization implements OrganizationService.
func (o *organizationService) DeleteOrganization(ctx context.Context, id string) error {
	err := o.orgRepo.DeleteOrganization(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeleteOrganization, err)
	}

	return nil
}

// GetUserOrganization implements OrganizationService.
func (o *organizationService) GetUserOrganization(ctx context.Context, userId string) (*appDto.OrganizationResponse, error) {
	org, err := o.orgRepo.GetUserOrganization(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetOrganization, err)
	}

	return &appDto.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description,
		ParentID:    org.ParentID,
		CreatedBy:   org.CreatedBy,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}, nil
}

// AssignUserToOrganization implements OrganizationService.
func (o *organizationService) AssignUserToOrganization(ctx context.Context, userId, organizationId string) error {
	err := o.orgRepo.AssignUserToOrganization(ctx, userId, organizationId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToAssignUserToOrganization, err)
	}

	return nil
}

// GetOrganizationUsers implements OrganizationService.
func (o *organizationService) GetOrganizationUsers(ctx context.Context, organizationId string) ([]*appDto.UserResponse, error) {
	users, err := o.orgRepo.GetOrganizationUsers(ctx, organizationId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetOrganization, err)
	}

	var responses []*appDto.UserResponse
	for _, user := range users {
		responses = append(responses, &appDto.UserResponse{
			ID:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			IsVerified:     user.IsVerified,
			TwoFAEnabled:   user.TwoFAEnabled,
			OrganizationID: user.OrganizationID,
		})
	}

	return responses, nil
}

func NewOrganizationService(orgRepo orgRepo.OrganizationRepository) OrganizationService {
	return &organizationService{
		orgRepo: orgRepo,
	}
}
