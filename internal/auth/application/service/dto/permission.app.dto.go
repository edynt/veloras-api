package dto

type PermissionAppDTO struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	ResourceType   string `json:"resource_type"`
	ResourceAction string `json:"resource_action"`
	OrganizationID string `json:"organization_id"`
}

type PermissionOutPut struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	ResourceType   string `json:"resource_type"`
	ResourceAction string `json:"resource_action"`
	OrganizationID string `json:"organization_id"`
	CreatedAt      int64  `json:"created_at"`
}

type CreatePermissionRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	ResourceType   string `json:"resource_type"`
	ResourceAction string `json:"resource_action"`
	OrganizationID string `json:"organization_id"`
}

type UpdatePermissionRequest struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	ResourceType   string `json:"resource_type"`
	ResourceAction string `json:"resource_action"`
}

type CheckPermissionRequest struct {
	UserID         string `json:"user_id"`
	ResourceType   string `json:"resource_type"`
	ResourceAction string `json:"resource_action"`
}

type CheckPermissionResponse struct {
	HasPermission bool   `json:"has_permission"`
	Message       string `json:"message"`
}

type GetPermissionsByResourceRequest struct {
	ResourceType string `json:"resource_type"`
}

type GetPermissionsByOrganizationRequest struct {
	OrganizationID string `json:"organization_id"`
}
