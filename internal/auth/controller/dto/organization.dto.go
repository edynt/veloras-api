package dto

type CreateOrganizationReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id"`
}

type UpdateOrganizationReq struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type AssignUserToOrganizationReq struct {
	UserID         string `json:"user_id" binding:"required"`
	OrganizationID string `json:"organization_id" binding:"required"`
}

type GetOrganizationUsersReq struct {
	OrganizationID string `json:"organization_id" binding:"required"`
}
