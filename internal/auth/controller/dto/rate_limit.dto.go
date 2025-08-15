package dto

type CheckRateLimitReq struct {
	UserID    string `json:"user_id"`
	Action    string `json:"action" binding:"required"`
	IPAddress string `json:"ip_address" binding:"required"`
}

type RateLimitConfigReq struct {
	MaxAttempts  int   `json:"max_attempts" binding:"required,min=1"`
	WindowSize   int64 `json:"window_size" binding:"required,min=1"`
	BlockDuration int64 `json:"block_duration" binding:"required,min=1"`
}
