package dto

type PermissionReq struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	ResourceType   string `json:"resource_type" binding:"required"`
	ResourceAction string `json:"resource_action" binding:"required"`
	OrganizationID string `json:"organization_id"`
}

type CreatePermissionReq struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	ResourceType   string `json:"resource_type" binding:"required"`
	ResourceAction string `json:"resource_action" binding:"required"`
	OrganizationID string `json:"organization_id"`
}

type UpdatePermissionReq struct {
	ID             string `json:"id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	ResourceType   string `json:"resource_type" binding:"required"`
	ResourceAction string `json:"resource_action" binding:"required"`
}

type CheckPermissionReq struct {
	UserID         string `json:"user_id" binding:"required"`
	ResourceType   string `json:"resource_type" binding:"required"`
	ResourceAction string `json:"resource_action" binding:"required"`
}

type GetPermissionsByResourceReq struct {
	ResourceType string `json:"resource_type" binding:"required"`
}

type GetPermissionsByOrganizationReq struct {
	OrganizationID string `json:"organization_id" binding:"required"`
}
