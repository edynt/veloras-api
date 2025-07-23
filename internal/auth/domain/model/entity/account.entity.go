package entity

type Account struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"fisrt_name"`
	LastName    string `json:"last_name"`
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
