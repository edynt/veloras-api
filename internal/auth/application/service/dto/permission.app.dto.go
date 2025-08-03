package dto

type PermissionAppDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionOutPut struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
