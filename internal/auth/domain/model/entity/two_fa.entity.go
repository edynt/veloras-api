package entity

type TwoFACode struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Code      string `json:"code"`
	ExpiresAt int64  `json:"expires_at"`
	Type      string `json:"type"`
	CreatedAt int64  `json:"created_at"`
}

type TwoFASetupRequest struct {
	UserID string `json:"user_id"`
}

type TwoFAVerifyRequest struct {
	UserID string `json:"user_id"`
	Code   string `json:"code"`
	Type   string `json:"type"`
}

type TwoFAEnableRequest struct {
	UserID string `json:"user_id"`
	Secret string `json:"secret"`
}
