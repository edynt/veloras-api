package dto

type AccountAppDTO struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Language    string `json:"language"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"fisrt_name"`
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
	FirstName      string `json:"fisrt_name"`
	LastName       string `json:"last_name"`
	Status         int    `json:"status"`
	Language       string `json:"language"`
	IsVerified     bool   `json:"is_verified"`
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	TokenExpiresAt int64  `json:"token_expires_at"`
}
