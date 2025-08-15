package entity

type RateLimit struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Action      string `json:"action"`
	IPAddress   string `json:"ip_address"`
	Attempts    int    `json:"attempts"`
	WindowStart int64  `json:"window_start"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type RateLimitCheck struct {
	UserID    string `json:"user_id"`
	Action    string `json:"action"`
	IPAddress string `json:"ip_address"`
}
