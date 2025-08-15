package service

import (
	"context"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
)

type AuthService interface {
	// User Management
	CreateUser(ctx context.Context, accountAppDTO appDto.AccountAppDTO) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*appDto.UserResponse, error)
	GetUserById(ctx context.Context, userId string) (*appDto.UserResponse, error)
	UpdateUserProfile(ctx context.Context, userId string, firstName, lastName, phoneNumber, language string) error
	UpdateUserPassword(ctx context.Context, userId string, currentPassword, newPassword string) error
	
	// Authentication
	VerifyUser(ctx context.Context, verificationEmailAppDTO appDto.EmailVerification) (bool, error)
	LoginUser(ctx context.Context, loginReq appDto.LoginRequest) (*appDto.LoginResponse, error)
	LogoutUser(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*appDto.RefreshTokenResponse, error)
	
	// Password Management
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email, resetToken, newPassword string) error
	ChangePassword(ctx context.Context, userId, currentPassword, newPassword string) error
	
	// 2FA Management
	Setup2FA(ctx context.Context, userId string) (*appDto.TwoFASetupResponse, error)
	Verify2FA(ctx context.Context, userId, code, codeType string) (*appDto.TwoFAVerifyResponse, error)
	Enable2FA(ctx context.Context, userId, secret, code string) error
	Disable2FA(ctx context.Context, userId, code string) error
	Get2FAStatus(ctx context.Context, userId string) (*appDto.TwoFAStatusResponse, error)
	
	// Session Management
	GetUserSessions(ctx context.Context, userId string) ([]*appDto.SessionInfo, error)
	DeleteSession(ctx context.Context, sessionId int32) error
	DeleteAllUserSessions(ctx context.Context, userId string) error
	
	// Account Security
	IncrementFailedLoginAttempts(ctx context.Context, userId string) error
	ResetFailedLoginAttempts(ctx context.Context, userId string) error
	LockUserAccount(ctx context.Context, userId string) error
	CheckAccountLocked(ctx context.Context, userId string) (bool, error)
}
