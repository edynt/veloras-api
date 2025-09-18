package dto

type RoleReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AssignPermissionReq struct {
	RoleID       int `json:"role_id"`
	PermissionID int `json:"permission_id"`
}
