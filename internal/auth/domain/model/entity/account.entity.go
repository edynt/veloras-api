package entity

type Account struct {
	ID                  string `json:"id"`
	Email               string `json:"email"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	PhoneNumber         string `json:"phone_number"`
	Language            string `json:"language"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	Status              int    `json:"status"`
	IsVerified          bool   `json:"is_verified"`
	TwoFASecret         string `json:"two_fa_secret"`
	TwoFAEnabled        bool   `json:"two_fa_enabled"`
	FailedLoginAttempts int    `json:"failed_login_attempts"`
	LockedUntil         int64  `json:"locked_until"`
	OrganizationID      string `json:"organization_id"`
	LastLoginAt         int64  `json:"last_login_at"`
}

type EmailVerification struct {
	UserID    string `json:"user_id"`
	Code      int    `json:"code"`
	ExpiresAt int64  `json:"expires_at"`
}

type UpdateUserStatus struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
}

type Session struct {
	ID           int32  `json:"id"`
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	CreatedAt    int64  `json:"created_at"`
	IsActive     bool   `json:"is_active"`
	LastActivity int64  `json:"last_activity"`
}

type PasswordReset struct {
	UserID     string `json:"user_id"`
	ResetToken string `json:"reset_token"`
	ExpiresAt  int64  `json:"expires_at"`
}

type ChangePasswordRequest struct {
	UserID          string `json:"user_id"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	ResetToken  string `json:"reset_token"`
	NewPassword string `json:"new_password"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
