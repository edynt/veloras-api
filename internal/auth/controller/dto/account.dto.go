package dto

type UserRegisterReq struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Password    string `json:"password" binding:"required,min=8"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address"`
	Country     string `json:"country"`
	Language    string `json:"language" binding:"required"`
}

type UserLoginReq struct {
	Username   string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	TwoFACode string `json:"two_fa_code,omitempty"`
}

type ChangePasswordReq struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

type ForgotPasswordReq struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordReq struct {
	Email       string `json:"email" binding:"required,email"`
	ResetToken  string `json:"reset_token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type LogoutReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TwoFASetupReq struct {
	UserID string `json:"user_id" binding:"required"`
}

type TwoFAVerifyReq struct {
	UserID string `json:"user_id" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Type   string `json:"type" binding:"required"`
}

type TwoFAEnableReq struct {
	UserID string `json:"user_id" binding:"required"`
	Secret string `json:"secret" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

type TwoFADisableReq struct {
	UserID string `json:"user_id" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

type TwoFAStatusReq struct {
	UserID string `json:"user_id" binding:"required"`
}
