package dto

type PermissionReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
