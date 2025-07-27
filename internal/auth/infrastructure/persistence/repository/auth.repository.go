package repository

import (
	"context"
	"fmt"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/internal/shared/gen"
	authsqlc "github.com/edynnt/veloras-api/internal/shared/gen"
	"github.com/edynnt/veloras-api/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepository struct {
	db *authsqlc.Queries
}

// SaveToken implements repository.AuthRepository.
func (a *authRepository) SaveToken(ctx context.Context, token *entity.Session) error {
	var param gen.CreateSessionParams
	if err := utils.SafeCopy(&param, token); err != nil {
		return err
	}

	_, err := a.db.CreateSession(ctx, param)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// DeleteVerificationCode implements repository.AuthRepository.
func (a *authRepository) DeleteVerificationCode(ctx context.Context, userId string, code int) error {
	// Parse the userId string to UUID
	convertId, err := utils.ConvertUUID(userId)

	if err != nil {
		return fmt.Errorf("invalid userId: %w", err)
	}

	// Call the DB update function
	err = a.db.DeleteVerificationCode(ctx, convertId)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	return nil
}

// ActiveUser implements repository.AuthRepository.
func (a *authRepository) ActiveUser(ctx context.Context, userId string) error {
	// Parse the userId string to UUID
	convertId, err := utils.ConvertUUID(userId)

	if err != nil {
		return fmt.Errorf("invalid userId: %w", err)
	}

	// Call the DB update function
	_, err = a.db.ActiveUser(ctx, convertId)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	return nil
}

// GetUserByUsername implements repository.AuthRepository.
func (a *authRepository) GetUserByUsername(ctx context.Context, userName string) (*entity.Account, error) {
	res, err := a.db.GetUserByUsername(ctx, userName)

	if err != nil {
		return nil, err
	}

	var entityResult entity.Account
	if err := utils.SafeCopy(&entityResult, &res); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// UpdateUserStatus implements repository.AuthRepository.
func (a *authRepository) UpdateUserStatus(ctx context.Context, userId string, status int) error {
	// Parse the userId string to UUID
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("invalid userId: %w", err)
	}

	// Construct the parameter object
	params := gen.UpdateUserStatusParams{
		ID: convertId,
		Status: pgtype.Int4{
			Int32: int32(status),
			Valid: true,
		},
	}

	// Call the DB update function
	_, err = a.db.UpdateUserStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	return nil
}

// GetVerificationCode implements repository.AuthRepository.
func (a *authRepository) GetVerificationCode(ctx context.Context, userId string, code int) (*entity.EmailVerification, error) {
	var param gen.GetEmailVerificationParams
	if err := utils.SafeCopy(&param, &entity.EmailVerification{UserID: userId, Code: code}); err != nil {
		return nil, err
	}

	result, err := a.db.GetEmailVerification(ctx, param)

	if err != nil {
		return nil, err
	}

	var entityResult entity.EmailVerification
	if err := utils.SafeCopy(&entityResult, &result); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// CreateVerificationCode implements repository.AuthRepository.
func (a *authRepository) CreateVerificationCode(ctx context.Context, userVerification *entity.EmailVerification) error {
	var param gen.CreateEmailVerificationParams
	if err := utils.SafeCopy(&param, &userVerification); err != nil {
		return err
	}

	_, err := a.db.CreateEmailVerification(ctx, param)

	if err != nil {
		return err
	}

	return nil

}

// EmailExists implements repository.AuthRepository.
func (a *authRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return a.db.GetUsernameExists(ctx, email)
}

// CreateUser implements repository.AuthRepository.
func (a *authRepository) CreateUser(ctx context.Context, account *entity.Account) (string, error) {

	var param gen.CreateUserParams
	if err := utils.SafeCopy(&param, &account); err != nil {
		return "", err
	}

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
