package repository

import (
	"context"
	"fmt"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/internal/shared/gen"
	authsqlc "github.com/edynnt/veloras-api/internal/shared/gen"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepository struct {
	db *authsqlc.Queries
}

// SaveToken implements repository.AuthRepository.
func (a *authRepository) SaveToken(ctx context.Context, token *entity.Session) error {
	// Map explicitly to ensure UserID is valid
	param := gen.CreateSessionParams{
		UserID: pgtype.Int4{
			Int32: int32(token.UserID),
			Valid: true,
		},
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.ExpiresAt,
	}

	_, err := a.db.CreateSession(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToCreateSession, err)
	}

	return nil
}

// DeleteVerificationCode implements repository.AuthRepository.
func (a *authRepository) DeleteVerificationCode(ctx context.Context, userId int, code int) error {
	// Call the DB update function
	err := a.db.DeleteVerificationCode(ctx, pgtype.Int4{
		Int32: int32(userId),
		Valid: true,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdateUserStatus, err)
	}

	return nil
}

// ActiveUser implements repository.AuthRepository.
func (a *authRepository) ActiveUser(ctx context.Context, userId int) error {
	// Call the DB update function
	_, err := a.db.ActiveUser(ctx, int32(userId))
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToActiveUser, err)
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
func (a *authRepository) UpdateUserStatus(ctx context.Context, userId int, status int) error {
	// Construct the parameter object
	params := gen.UpdateUserStatusParams{
		ID: int32(userId),
		Status: pgtype.Int4{
			Int32: int32(status),
			Valid: true,
		},
	}

	// Call the DB update function
	_, err := a.db.UpdateUserStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdateUserStatus, err)
	}

	return nil
}

// GetVerificationCode implements repository.AuthRepository.
func (a *authRepository) GetVerificationCode(ctx context.Context, userId int, code int) (*entity.EmailVerification, error) {
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
func (a *authRepository) CreateUser(ctx context.Context, account *entity.Account) (int, error) {

	var param gen.CreateUserParams
	if err := utils.SafeCopy(&param, &account); err != nil {
		return 0, err
	}

	createdAccount, err := a.db.CreateUser(ctx, param)

	if err != nil {
		return 0, err
	}

	return int(createdAccount.ID), nil
}

// UsernameExists implements repository.AuthRepository.
func (a *authRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	return a.db.GetUsernameExists(ctx, username)
}

func NewAuthRepository(db *pgxpool.Pool) repository.AuthRepository {
	queries := authsqlc.New(db) // db is *pgxpool.Pool
	return &authRepository{db: queries}
}

// RefreshToken implements repository.AuthRepository.
// Note: Basic stub verification that defers cryptographic/expiry validation to service layer.
// Optionally, this could verify existence against the sessions table if a query exists.
func (a *authRepository) RefreshToken(ctx context.Context, refreshToken string) error {
	// Without an sqlc method to lookup by token, we accept the validated token from service.
	// Extend later to check presence/blacklist or rotation in DB.
	if refreshToken == "" {
		return fmt.Errorf("%s", msg.InvalidRefreshToken)
	}
	return nil
}
