package dto

type AuditLogResponse struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"user_id"`
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   string                 `json:"resource_id"`
	Details      map[string]interface{} `json:"details"`
	IPAddress    string                 `json:"ip_address"`
	UserAgent    string                 `json:"user_agent"`
	CreatedAt    int64                  `json:"created_at"`
}

type CreateAuditLogRequest struct {
	UserID       string                 `json:"user_id"`
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   string                 `json:"resource_id"`
	Details      map[string]interface{} `json:"details"`
	IPAddress    string                 `json:"ip_address"`
	UserAgent    string                 `json:"user_agent"`
}

type GetAuditLogsRequest struct {
	UserID       string `json:"user_id,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
	ResourceID   string `json:"resource_id,omitempty"`
	Action       string `json:"action,omitempty"`
	StartDate    int64  `json:"start_date,omitempty"`
	EndDate      int64  `json:"end_date,omitempty"`
	Limit        int    `json:"limit,omitempty"`
}
