package repository

import (
	"context"
	"fmt"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/internal/shared/gen"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type organizationRepository struct {
	db *gen.Queries
}

// CreateOrganization implements repository.OrganizationRepository.
func (o *organizationRepository) CreateOrganization(ctx context.Context, org *entity.Organization) (string, error) {
	var param gen.CreateOrganizationParams
	if err := utils.SafeCopy(&param, org); err != nil {
		return "", err
	}

	result, err := o.db.CreateOrganization(ctx, param)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.FailedToCreateOrganization, err)
	}

	return result.ID.String(), nil
}

// GetOrganizationById implements repository.OrganizationRepository.
func (o *organizationRepository) GetOrganizationById(ctx context.Context, id string) (*entity.Organization, error) {
	convertId, err := utils.ConvertUUID(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.OrganizationIdInvalid, err)
	}

	result, err := o.db.GetOrganizationById(ctx, convertId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetOrganization, err)
	}

	var entityResult entity.Organization
	if err := utils.SafeCopy(&entityResult, &result); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// GetOrganizationsByParent implements repository.OrganizationRepository.
func (o *organizationRepository) GetOrganizationsByParent(ctx context.Context, parentId string) ([]*entity.Organization, error) {
	convertId, err := utils.ConvertUUID(parentId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.OrganizationIdInvalid, err)
	}

	results, err := o.db.GetOrganizationsByParent(ctx, convertId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetOrganization, err)
	}

	var entityResults []*entity.Organization
	for _, result := range results {
		var entityResult entity.Organization
		if err := utils.SafeCopy(&entityResult, &result); err != nil {
			return nil, err
		}
		entityResults = append(entityResults, &entityResult)
	}

	return entityResults, nil
}

// UpdateOrganization implements repository.OrganizationRepository.
func (o *organizationRepository) UpdateOrganization(ctx context.Context, org *entity.Organization) error {
	convertId, err := utils.ConvertUUID(org.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.OrganizationIdInvalid, err)
	}

	now := utils.GetNowUnix()
	param := gen.UpdateOrganizationParams{
		Name:        org.Name,
		Description: pgtype.Text{String: org.Description, Valid: true},
		UpdatedAt:   pgtype.Int8{Int64: now, Valid: true},
		ID:          convertId,
	}

	err = o.db.UpdateOrganization(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdateOrganization, err)
	}

	return nil
}

// DeleteOrganization implements repository.OrganizationRepository.
func (o *organizationRepository) DeleteOrganization(ctx context.Context, id string) error {
	convertId, err := utils.ConvertUUID(id)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.OrganizationIdInvalid, err)
	}

	err = o.db.DeleteOrganization(ctx, convertId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeleteOrganization, err)
	}

	return nil
}

// GetUserOrganization implements repository.OrganizationRepository.
func (o *organizationRepository) GetUserOrganization(ctx context.Context, userId string) (*entity.Organization, error) {
	convertUserId, err := utils.ConvertUUID(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	result, err := o.db.GetUserOrganization(ctx, convertUserId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetOrganization, err)
	}

	var entityResult entity.Organization
	if err := utils.SafeCopy(&entityResult, &result); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// AssignUserToOrganization implements repository.OrganizationRepository.
func (o *organizationRepository) AssignUserToOrganization(ctx context.Context, userId string, organizationId string) error {
	convertUserId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	convertOrgId, err := utils.ConvertUUID(organizationId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.OrganizationIdInvalid, err)
	}

	param := gen.AssignUserToOrganizationParams{
		OrganizationID: convertOrgId,
		ID:             convertUserId,
	}

	err = o.db.AssignUserToOrganization(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToAssignUserToOrganization, err)
	}

	return nil
}

// GetOrganizationUsers implements repository.OrganizationRepository.
func (o *organizationRepository) GetOrganizationUsers(ctx context.Context, organizationId string) ([]*entity.Account, error) {
	convertId, err := utils.ConvertUUID(organizationId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.OrganizationIdInvalid, err)
	}

	results, err := o.db.GetOrganizationUsers(ctx, convertId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetOrganization, err)
	}

	var entityResults []*entity.Account
	for _, result := range results {
		var entityResult entity.Account
		if err := utils.SafeCopy(&entityResult, &result); err != nil {
			return nil, err
		}
		entityResults = append(entityResults, &entityResult)
	}

	return entityResults, nil
}

func NewOrganizationRepository(db *pgxpool.Pool) repository.OrganizationRepository {
	queries := gen.New(db)
	return &organizationRepository{db: queries}
}
