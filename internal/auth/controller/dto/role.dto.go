package dto

type RoleReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateRoleReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateRoleReq struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type AssignPermissionToRoleReq struct {
	RoleID       string `json:"role_id" binding:"required"`
	PermissionID string `json:"permission_id" binding:"required"`
}

type AssignRoleToUserReq struct {
	UserID string `json:"user_id" binding:"required"`
	RoleID string `json:"role_id" binding:"required"`
}

type GetUserRolesReq struct {
	UserID string `json:"user_id" binding:"required"`
}

type GetRolePermissionsReq struct {
	RoleID string `json:"role_id" binding:"required"`
}
