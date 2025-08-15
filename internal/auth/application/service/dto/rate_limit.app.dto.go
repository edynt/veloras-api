package dto

type RateLimitResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Action      string `json:"action"`
	IPAddress   string `json:"ip_address"`
	Attempts    int    `json:"attempts"`
	WindowStart int64  `json:"window_start"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type CheckRateLimitRequest struct {
	UserID    string `json:"user_id"`
	Action    string `json:"action"`
	IPAddress string `json:"ip_address"`
}

type RateLimitConfig struct {
	MaxAttempts  int   `json:"max_attempts"`
	WindowSize   int64 `json:"window_size"`
	BlockDuration int64 `json:"block_duration"`
}
