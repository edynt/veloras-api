package repository

import (
	"context"
	"fmt"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/internal/shared/gen"
	authsqlc "github.com/edynnt/veloras-api/internal/shared/gen"
	"github.com/edynnt/veloras-api/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
)

type authRepository struct {
	db *authsqlc.Queries
}

// UpdateUserStatus implements repository.AuthRepository.
func (a *authRepository) UpdateUserStatus(ctx context.Context, userID string, status int) (string, error) {
	// Parse the userID string to UUID
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", fmt.Errorf("invalid userID: %w", err)
	}

	// Construct the parameter object
	params := gen.UpdateUserStatusParams{
		ID: pgtype.UUID{
			Bytes: parsedUUID,
			Valid: true,
		},
		Status: pgtype.Int4{
			Int32: int32(status),
			Valid: true,
		},
	}

	// Call the DB update function
	result, err := a.db.UpdateUserStatus(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to update user status: %w", err)
	}

	// Return the updated ID
	return result.ID.String(), nil
}

// GetVerificationCode implements repository.AuthRepository.
func (a *authRepository) GetVerificationCode(ctx context.Context, userID string, code int) (*entity.EmailVerification, error) {
	var param gen.GetEmailVerificationParams
	copier.Copy(&param, &entity.EmailVerification{UserID: userID, Code: code})

	result, err := a.db.GetEmailVerification(ctx, param)

	if err != nil {
		return nil, err
	}

	var entityResult entity.EmailVerification
	copier.Copy(&entityResult, &result)

	return &entityResult, nil
}

// CreateVerificationCode implements repository.AuthRepository.
func (a *authRepository) CreateVerificationCode(ctx context.Context, userVerification *entity.EmailVerification) (string, error) {
	var param gen.CreateEmailVerificationParams
	copier.Copy(&param, &userVerification)
	createdEmailVerification, err := a.db.CreateEmailVerification(ctx, param)

	if err != nil {
		return "", err
	}

	return utils.Int32ToString(createdEmailVerification.ID), nil

}

// EmailExists implements repository.AuthRepository.
func (a *authRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return a.db.GetUsernameExists(ctx, email)
}

// CreateUser implements repository.AuthRepository.
func (a *authRepository) CreateUser(ctx context.Context, account *entity.Account) (string, error) {

	var param gen.CreateUserParams
	copier.Copy(&param, &account)

	createdAccount, err := a.db.CreateUser(ctx, param)

	if err != nil {
		return "", err
	}

	return createdAccount.ID.String(), nil
}

// UsernameExists implements repository.AuthRepository.
func (a *authRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	return a.db.GetUsernameExists(ctx, username)
}

func NewAuthRepository(db *pgxpool.Pool) repository.AuthRepository {
	queries := authsqlc.New(db) // db is *pgxpool.Pool
	return &authRepository{db: queries}
}
