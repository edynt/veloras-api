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
}
