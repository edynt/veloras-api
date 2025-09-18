package entity

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RolePermission struct {
	RoleID       int `json:"role_id"`
	PermissionID int `json:"permission_id"`
}
