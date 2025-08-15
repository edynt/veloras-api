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
	var param gen.CreateSessionParams
	if err := utils.SafeCopy(&param, token); err != nil {
		return err
	}

	_, err := a.db.CreateSession(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToCreateSession, err)
	}

	return nil
}

// DeleteVerificationCode implements repository.AuthRepository.
func (a *authRepository) DeleteVerificationCode(ctx context.Context, userId string, code int) error {
	// Parse the userId string to UUID
	convertId, err := utils.ConvertUUID(userId)

	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	// Call the DB update function
	err = a.db.DeleteVerificationCode(ctx, convertId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdateUserStatus, err)
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
func (a *authRepository) UpdateUserStatus(ctx context.Context, userId string, status int) error {
	// Parse the userId string to UUID
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
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
		return fmt.Errorf("%s: %w", msg.FailedToUpdateUserStatus, err)
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

// GetUserById implements repository.AuthRepository.
func (a *authRepository) GetUserById(ctx context.Context, userId string) (*entity.Account, error) {
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	res, err := a.db.GetUserById(ctx, convertId)
	if err != nil {
		return nil, err
	}

	var entityResult entity.Account
	if err := utils.SafeCopy(&entityResult, &res); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// UpdateUserProfile implements repository.AuthRepository.
func (a *authRepository) UpdateUserProfile(ctx context.Context, userId string, firstName, lastName, phoneNumber, language string) error {
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.UpdateUserProfileParams{
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: phoneNumber,
		Language:    pgtype.Text{String: language, Valid: true},
		ID:          convertId,
	}

	err = a.db.UpdateUserProfile(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdateUserProfile, err)
	}

	return nil
}

// UpdateUserPassword implements repository.AuthRepository.
func (a *authRepository) UpdateUserPassword(ctx context.Context, userId string, newPassword string) error {
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.UpdateUserPasswordParams{
		Password: newPassword,
		ID:       convertId,
	}

	err = a.db.UpdateUserPassword(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdatePassword, err)
	}

	return nil
}

// UpdateUser2FASecret implements repository.AuthRepository.
func (a *authRepository) UpdateUser2FASecret(ctx context.Context, userId string, secret string) error {
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.UpdateUser2FASecretParams{
		TwoFaSecret: pgtype.Text{String: secret, Valid: true},
		ID:          convertId,
	}

	err = a.db.UpdateUser2FASecret(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdate2FASecret, err)
	}

	return nil
}

// UpdateUser2FAEnabled implements repository.AuthRepository.
func (a *authRepository) UpdateUser2FAEnabled(ctx context.Context, userId string, enabled bool) error {
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.UpdateUser2FAEnabledParams{
		TwoFaEnabled: pgtype.Bool{Bool: enabled, Valid: true},
		ID:           convertId,
	}

	err = a.db.UpdateUser2FAEnabled(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdate2FAEnabled, err)
	}

	return nil
}

// UpdateLastLogin implements repository.AuthRepository.
func (a *authRepository) UpdateLastLogin(ctx context.Context, userId string) error {
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	now := utils.GetNowUnix()
	param := gen.UpdateLastLoginParams{
		LastLoginAt: pgtype.Int8{Int64: now, Valid: true},
		ID:          convertId,
	}

	err = a.db.UpdateLastLogin(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdateLastLogin, err)
	}

	return nil
}

// IncrementFailedLoginAttempts implements repository.AuthRepository.
func (a *authRepository) IncrementFailedLoginAttempts(ctx context.Context, userId string) error {
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	err = a.db.IncrementFailedLoginAttempts(ctx, convertId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToIncrementFailedLoginAttempts, err)
	}

	return nil
}

// ResetFailedLoginAttempts implements repository.AuthRepository.
func (a *authRepository) ResetFailedLoginAttempts(ctx context.Context, userId string) error {
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	err = a.db.ResetFailedLoginAttempts(ctx, convertId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToResetFailedLoginAttempts, err)
	}

	return nil
}

// LockUserAccount implements repository.AuthRepository.
func (a *authRepository) LockUserAccount(ctx context.Context, userId string, status int, lockedUntil int64) error {
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.LockUserAccountParams{
		Status:     pgtype.Int4{Int32: int32(status), Valid: true},
		LockedUntil: pgtype.Int8{Int64: lockedUntil, Valid: true},
		ID:         convertId,
	}

	err = a.db.LockUserAccount(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToLockUserAccount, err)
	}

	return nil
}

// CreatePasswordReset implements repository.AuthRepository.
func (a *authRepository) CreatePasswordReset(ctx context.Context, userID string, resetToken string, expiresAt int64) error {
	convertId, err := utils.ConvertUUID(userID)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.CreatePasswordResetParams{
		UserID:     convertId,
		ResetToken: resetToken,
		ExpiresAt:  expiresAt,
	}

	err = a.db.CreatePasswordReset(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToCreatePasswordReset, err)
	}

	return nil
}

// GetPasswordReset implements repository.AuthRepository.
func (a *authRepository) GetPasswordReset(ctx context.Context, userID string, resetToken string) (*entity.PasswordReset, error) {
	convertId, err := utils.ConvertUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.GetPasswordResetParams{
		UserID:     convertId,
		ResetToken: resetToken,
	}

	result, err := a.db.GetPasswordReset(ctx, param)
	if err != nil {
		return nil, err
	}

	var entityResult entity.PasswordReset
	if err := utils.SafeCopy(&entityResult, &result); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// DeletePasswordReset implements repository.AuthRepository.
func (a *authRepository) DeletePasswordReset(ctx context.Context, userID string, resetToken string) error {
	convertId, err := utils.ConvertUUID(userID)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.DeletePasswordResetParams{
		UserID:     convertId,
		ResetToken: resetToken,
	}

	err = a.db.DeletePasswordReset(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeletePasswordReset, err)
	}

	return nil
}

// GetSessionByRefreshToken implements repository.AuthRepository.
func (a *authRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*entity.Session, error) {
	result, err := a.db.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	var entityResult entity.Session
	if err := utils.SafeCopy(&entityResult, &result); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// GetSessionsByUser implements repository.AuthRepository.
func (a *authRepository) GetSessionsByUser(ctx context.Context, userId string) ([]*entity.Session, error) {
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	results, err := a.db.GetSessionsByUser(ctx, convertId)
	if err != nil {
		return nil, err
	}

	var entityResults []*entity.Session
	for _, result := range results {
		var entityResult entity.Session
		if err := utils.SafeCopy(&entityResult, &result); err != nil {
			return nil, err
		}
		entityResults = append(entityResults, &entityResult)
	}

	return entityResults, nil
}

// DeleteSessionByRefreshToken implements repository.AuthRepository.
func (a *authRepository) DeleteSessionByRefreshToken(ctx context.Context, refreshToken string) error {
	err := a.db.DeleteSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeleteSession, err)
	}

	return nil
}

// DeleteAllUserSessions implements repository.AuthRepository.
func (a *authRepository) DeleteAllUserSessions(ctx context.Context, userId string) error {
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	err = a.db.DeleteAllUserSessions(ctx, convertId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeleteAllUserSessions, err)
	}

	return nil
}

// DeleteExpiredSessions implements repository.AuthRepository.
func (a *authRepository) DeleteExpiredSessions(ctx context.Context, currentTime int64) error {
	err := a.db.DeleteExpiredSessions(ctx, currentTime)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeleteExpiredSessions, err)
	}

	return nil
}

// UpdateSessionExpiry implements repository.AuthRepository.
func (a *authRepository) UpdateSessionExpiry(ctx context.Context, sessionId int32, expiresAt int64) error {
	param := gen.UpdateSessionExpiryParams{
		ExpiresAt: expiresAt,
		ID:        sessionId,
	}

	err := a.db.UpdateSessionExpiry(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToUpdateSessionExpiry, err)
	}

	return nil
}

// Create2FACode implements repository.AuthRepository.
func (a *authRepository) Create2FACode(ctx context.Context, userId string, code string, expiresAt int64, codeType string) error {
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.Create2FACodeParams{
		UserID: convertId,
		Code:   code,
		ExpiresAt: expiresAt,
		Type:   codeType,
	}

	_, err = a.db.Create2FACode(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToCreate2FACode, err)
	}

	return nil
}

// Get2FACode implements repository.AuthRepository.
func (a *authRepository) Get2FACode(ctx context.Context, userId string, code string, codeType string) (*entity.TwoFACode, error) {
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.Get2FACodeParams{
		UserID: convertId,
		Code:   code,
		Type:   codeType,
	}

	result, err := a.db.Get2FACode(ctx, param)
	if err != nil {
		return nil, err
	}

	var entityResult entity.TwoFACode
	if err := utils.SafeCopy(&entityResult, &result); err != nil {
		return nil, err
	}

	return &entityResult, nil
}

// Delete2FACode implements repository.AuthRepository.
func (a *authRepository) Delete2FACode(ctx context.Context, userId string, codeType string) error {
	convertId, err := utils.ConvertUUID(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.Delete2FACodeParams{
		UserID: convertId,
		Type:   codeType,
	}

	err = a.db.Delete2FACode(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDelete2FACode, err)
	}

	return nil
}

// DeleteExpired2FACodes implements repository.AuthRepository.
func (a *authRepository) DeleteExpired2FACodes(ctx context.Context, currentTime int64) error {
	err := a.db.DeleteExpired2FACodes(ctx, currentTime)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeleteExpired2FACodes, err)
	}

	return nil
}

// GetUsersByOrganization implements repository.AuthRepository.
func (a *authRepository) GetUsersByOrganization(ctx context.Context, organizationId string) ([]*entity.Account, error) {
	convertId, err := utils.ConvertUUID(organizationId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.OrganizationIdInvalid, err)
	}

	results, err := a.db.GetUsersByOrganization(ctx, convertId)
	if err != nil {
		return nil, err
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

// GetUsersByRole implements repository.AuthRepository.
func (a *authRepository) GetUsersByRole(ctx context.Context, roleId string) ([]*entity.Account, error) {
	convertId, err := utils.ConvertUUID(roleId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.RoleIdInvalid, err)
	}

	results, err := a.db.GetUsersByRole(ctx, convertId)
	if err != nil {
		return nil, err
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

func NewAuthRepository(db *pgxpool.Pool) repository.AuthRepository {
	queries := authsqlc.New(db) // db is *pgxpool.Pool
	return &authRepository{db: queries}
}
