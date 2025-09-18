package dto

type RoleAppDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleOutPut struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RolePermissionAppDTO struct {
	RoleID       int `json:"role_id"`
	PermissionID int `json:"permission_id"`
}
