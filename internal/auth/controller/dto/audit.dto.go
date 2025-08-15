package dto

type GetAuditLogsReq struct {
	UserID       string `json:"user_id"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	Action       string `json:"action"`
	StartDate    int64  `json:"start_date"`
	EndDate      int64  `json:"end_date"`
	Limit        int    `json:"limit"`
}

type CreateAuditLogReq struct {
	UserID       string                 `json:"user_id" binding:"required"`
	Action       string                 `json:"action" binding:"required"`
	ResourceType string                 `json:"resource_type" binding:"required"`
	ResourceID   string                 `json:"resource_id"`
	Details      map[string]interface{} `json:"details"`
	IPAddress    string                 `json:"ip_address"`
	UserAgent    string                 `json:"user_agent"`
}
