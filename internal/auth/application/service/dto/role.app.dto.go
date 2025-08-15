package dto

type RoleAppDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleOutPut struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
}

type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AssignPermissionToRoleRequest struct {
	RoleID       string `json:"role_id"`
	PermissionID string `json:"permission_id"`
}

type AssignRoleToUserRequest struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
}

type GetUserRolesRequest struct {
	UserID string `json:"user_id"`
}

type GetRolePermissionsRequest struct {
	RoleID string `json:"role_id"`
}

type RoleWithPermissions struct {
	Role        RoleOutPut      `json:"role"`
	Permissions []PermissionOutPut `json:"permissions"`
}
