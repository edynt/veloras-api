package entity

type PasswordReset struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	ResetToken string `json:"reset_token"`
	ExpiresAt  int64  `json:"expires_at"`
}
