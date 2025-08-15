package repository

import (
	"context"
	"fmt"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/internal/shared/gen"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type twoFARepository struct {
	db *gen.Queries
}

// Create2FACode implements repository.TwoFARepository.
func (t *twoFARepository) Create2FACode(ctx context.Context, code *entity.TwoFACode) (string, error) {
	var param gen.Create2FACodeParams
	if err := utils.SafeCopy(&param, code); err != nil {
		return "", err
	}

	result, err := t.db.Create2FACode(ctx, param)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.FailedToCreate2FACode, err)
	}

	return result.ID.String(), nil
}

// Get2FACode implements repository.TwoFARepository.
func (t *twoFARepository) Get2FACode(ctx context.Context, userId string, code string, codeType string) (*entity.TwoFACode, error) {
	convertUserId, err := utils.ConvertUUID(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.Get2FACodeParams{
		UserID: convertUserId,
		Code:   code,
		Type:   codeType,
	}

	result, err := t.db.Get2FACode(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGet2FACode, err)
	}

	var entityResult entity.TwoFACode
	if err := utils.SafeCopy(&entityResult, &result); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// Delete2FACode implements repository.TwoFARepository.
func (t *twoFARepository) Delete2FACode(ctx context.Context, userId string, codeType string) error {
	convertUserId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.Delete2FACodeParams{
		UserID: convertUserId,
		Type:   codeType,
	}

	err = t.db.Delete2FACode(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDelete2FACode, err)
	}

	return nil
}

// DeleteExpired2FACodes implements repository.TwoFARepository.
func (t *twoFARepository) DeleteExpired2FACodes(ctx context.Context, currentTime int64) error {
	err := t.db.DeleteExpired2FACodes(ctx, currentTime)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeleteExpired2FACodes, err)
	}

	return nil
}

// GetUser2FACodes implements repository.TwoFARepository.
func (t *twoFARepository) GetUser2FACodes(ctx context.Context, userId string, codeType string) ([]*entity.TwoFACode, error) {
	convertUserId, err := utils.ConvertUUID(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.GetUser2FACodesParams{
		UserID: convertUserId,
		Type:   codeType,
	}

	results, err := t.db.GetUser2FACodes(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGet2FACode, err)
	}

	var entityResults []*entity.TwoFACode
	for _, result := range results {
		var entityResult entity.TwoFACode
		if err := utils.SafeCopy(&entityResult, &result); err != nil {
			return nil, err
		}
		entityResults = append(entityResults, &entityResult)
	}

	return entityResults, nil
}

func NewTwoFARepository(db *pgxpool.Pool) repository.TwoFARepository {
	queries := gen.New(db)
	return &twoFARepository{db: queries}
}
