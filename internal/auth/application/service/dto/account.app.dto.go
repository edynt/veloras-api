package dto

type AccountAppDTO struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Language    string `json:"language"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

type EmailVerification struct {
	UserID string `json:"user_id"`
	Code   int    `json:"code"`
}

type UserOutPut struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Status         int    `json:"status"`
	Language       string `json:"language"`
	IsVerified     bool   `json:"is_verified"`
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	TokenExpiresAt int64  `json:"token_expires_at"`
}

type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	TwoFACode string `json:"two_fa_code,omitempty"`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	IsVerified     bool   `json:"is_verified"`
	TwoFAEnabled   bool   `json:"two_fa_enabled"`
	OrganizationID string `json:"organization_id"`
}

type TwoFASetupResponse struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code"`
}

type TwoFAVerifyResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type TwoFAStatusResponse struct {
	Enabled bool   `json:"enabled"`
	Secret  string `json:"secret,omitempty"`
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

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type SessionInfo struct {
	ID           int32  `json:"id"`
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	CreatedAt    int64  `json:"created_at"`
	LastActivity int64  `json:"last_activity"`
	IsActive     bool   `json:"is_active"`
}

type PasswordResetRequest struct {
	Email string `json:"email"`
}

type PasswordReset struct {
	Email       string `json:"email"`
	ResetToken  string `json:"reset_token"`
	NewPassword string `json:"new_password"`
}
