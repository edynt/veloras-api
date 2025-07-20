package entity

import "time"

type Account struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"fisrt_name"`
	LastName    string `json:"last_name"`
}

type CreateVerificationCode struct {
	UserID    string    `json:"user_id"`
	Code      int       `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}
