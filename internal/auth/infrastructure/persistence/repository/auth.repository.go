package repository

import (
	"context"
	"fmt"

	"github.com/edynt/chogiare/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/auth/domain/repository"
	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	authsqlc "github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/edynt/chogiare/veloras-api/pkg/response/msg"
	"github.com/edynt/chogiare/veloras-api/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepository struct {
	db   *authsqlc.Queries
	pool *pgxpool.Pool
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
	return a.db.GetUserEmailExists(ctx, email)
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

func NewAuthRepository(db *pgxpool.Pool) repository.AuthRepository {
	queries := authsqlc.New(db) // db is *pgxpool.Pool
	return &authRepository{db: queries, pool: db}
}

// RefreshToken implements repository.AuthRepository.
func (a *authRepository) RefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return fmt.Errorf("%s", msg.InvalidRefreshToken)
	}

	// Check if refresh token exists in database
	session, err := a.db.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.InvalidRefreshToken, err)
	}

	// Check if session is expired
	now := utils.GetNowUnix()
	if session.ExpiresAt < now {
		return fmt.Errorf("%s", msg.RefreshTokenExpired)
	}

	return nil
}

// DeleteSessionsByUser implements repository.AuthRepository.
func (a *authRepository) DeleteSessionsByUser(ctx context.Context, userId int) error {
	if _, err := a.pool.Exec(ctx, "DELETE FROM sessions WHERE user_id = $1", userId); err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeleteSession, err)
	}
	return nil
}

// DeleteSessionByRefreshToken implements repository.AuthRepository.
func (a *authRepository) DeleteSessionByRefreshToken(ctx context.Context, refreshToken string) error {
	err := a.db.DeleteSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeleteSession, err)
	}
	return nil
}

// GetUserByID implements repository.AuthRepository.
func (a *authRepository) GetUserByID(ctx context.Context, userId int) (*entity.Account, error) {
	res, err := a.db.GetUserByID(ctx, int32(userId))
	if err != nil {
		return nil, err
	}

	var entityResult entity.Account
	if err := utils.SafeCopy(&entityResult, &res); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// UpdateUserPassword implements repository.AuthRepository.
func (a *authRepository) UpdateUserPassword(ctx context.Context, userId int, hashedPassword string) error {
	params := gen.UpdateUserPasswordParams{
		ID:       int32(userId),
		Password: hashedPassword,
	}

	_, err := a.db.UpdateUserPassword(ctx, params)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdatePassword, err)
	}

	return nil
}

// GetUserByEmail implements repository.AuthRepository.
func (a *authRepository) GetUserByEmail(ctx context.Context, email string) (*entity.Account, error) {
	res, err := a.db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	var entityResult entity.Account
	if err := utils.SafeCopy(&entityResult, &res); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// CreatePasswordReset implements repository.AuthRepository.
func (a *authRepository) CreatePasswordReset(ctx context.Context, passwordReset *entity.PasswordReset) error {
	param := gen.CreatePasswordResetParams{
		UserID: pgtype.Int4{
			Int32: int32(passwordReset.UserID),
			Valid: true,
		},
		ResetToken: passwordReset.ResetToken,
		ExpiresAt:  passwordReset.ExpiresAt,
	}

	err := a.db.CreatePasswordReset(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToCreatePasswordReset, err)
	}

	return nil
}

// GetPasswordReset implements repository.AuthRepository.
func (a *authRepository) GetPasswordReset(ctx context.Context, userId int, token string) (*entity.PasswordReset, error) {
	res, err := a.db.GetPasswordReset(ctx, gen.GetPasswordResetParams{
		UserID: pgtype.Int4{
			Int32: int32(userId),
			Valid: true,
		},
		ResetToken: token,
	})
	if err != nil {
		return nil, err
	}

	var entityResult entity.PasswordReset
	if err := utils.SafeCopy(&entityResult, &res); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// GetPasswordResetByToken implements repository.AuthRepository.
func (a *authRepository) GetPasswordResetByToken(ctx context.Context, token string) (*entity.PasswordReset, error) {
	res, err := a.db.GetPasswordResetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	var entityResult entity.PasswordReset
	if err := utils.SafeCopy(&entityResult, &res); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// DeletePasswordReset implements repository.AuthRepository.
func (a *authRepository) DeletePasswordReset(ctx context.Context, userId int, token string) error {
	err := a.db.DeletePasswordReset(ctx, gen.DeletePasswordResetParams{
		UserID: pgtype.Int4{
			Int32: int32(userId),
			Valid: true,
		},
		ResetToken: token,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeletePasswordReset, err)
	}

	return nil
}

// UserHasPermission implements repository.AuthRepository.
func (a *authRepository) UserHasPermission(ctx context.Context, userID int, permissionName string) (bool, error) {
	hasPermission, err := a.db.UserHasPermission(ctx, gen.UserHasPermissionParams{
		UserID: int32(userID),
		Name:   permissionName,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check user permission: %w", err)
	}
	return hasPermission, nil
}

// UserHasRole implements repository.AuthRepository.
func (a *authRepository) UserHasRole(ctx context.Context, userID int, roleName string) (bool, error) {
	hasRole, err := a.db.UserHasRole(ctx, gen.UserHasRoleParams{
		UserID: int32(userID),
		Name:   roleName,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check user role: %w", err)
	}
	return hasRole, nil
}
