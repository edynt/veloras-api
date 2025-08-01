package entity

type Account struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Language    string `json:"language"`
	FirstName   string `json:"fisrt_name"`
	LastName    string `json:"last_name"`
	Status      int    `json:"status"`
	IsVerified  bool   `json:"is_verified"`
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
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}
