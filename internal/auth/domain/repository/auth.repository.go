package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
)

type AuthRepository interface {
	// User Management
	CreateUser(ctx context.Context, account *entity.Account) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.Account, error)
	GetUserByUsername(ctx context.Context, userName string) (*entity.Account, error)
	GetUserById(ctx context.Context, userId string) (*entity.Account, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	UpdateUserProfile(ctx context.Context, userId string, firstName, lastName, phoneNumber, language string) error
	UpdateUserPassword(ctx context.Context, userId string, newPassword string) error
	UpdateUser2FASecret(ctx context.Context, userId string, secret string) error
	UpdateUser2FAEnabled(ctx context.Context, userId string, enabled bool) error
	UpdateLastLogin(ctx context.Context, userId string) error

	// Account Status & Security
	UpdateUserStatus(ctx context.Context, userId string, status int) error
	ActiveUser(ctx context.Context, userId string) error
	IncrementFailedLoginAttempts(ctx context.Context, userId string) error
	ResetFailedLoginAttempts(ctx context.Context, userId string) error
	LockUserAccount(ctx context.Context, userId string, status int, lockedUntil int64) error

	// Email Verification
	CreateVerificationCode(ctx context.Context, userVerification *entity.EmailVerification) error
	GetVerificationCode(ctx context.Context, userId string, code int) (*entity.EmailVerification, error)
	DeleteVerificationCode(ctx context.Context, userId string, code int) error

	// Password Reset
	CreatePasswordReset(ctx context.Context, userID string, resetToken string, expiresAt int64) error
	GetPasswordReset(ctx context.Context, userID string, resetToken string) (*entity.PasswordReset, error)
	DeletePasswordReset(ctx context.Context, userID string, resetToken string) error

	// Session Management
	SaveToken(ctx context.Context, token *entity.Session) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*entity.Session, error)
	GetSessionsByUser(ctx context.Context, userId string) ([]*entity.Session, error)
	DeleteSessionByRefreshToken(ctx context.Context, refreshToken string) error
	DeleteAllUserSessions(ctx context.Context, userId string) error
	DeleteExpiredSessions(ctx context.Context, currentTime int64) error
	UpdateSessionExpiry(ctx context.Context, sessionId int32, expiresAt int64) error

	// 2FA Management
	Create2FACode(ctx context.Context, userId string, code string, expiresAt int64, codeType string) error
	Get2FACode(ctx context.Context, userId string, code string, codeType string) (*entity.TwoFACode, error)
	Delete2FACode(ctx context.Context, userId string, codeType string) error
	DeleteExpired2FACodes(ctx context.Context, currentTime int64) error

	// Organization Management
	GetUsersByOrganization(ctx context.Context, organizationId string) ([]*entity.Account, error)
	GetUsersByRole(ctx context.Context, roleId string) ([]*entity.Account, error)
}
