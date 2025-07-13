package dto

import "github.com/edynnt/veloras-api/internal/auth/domain/model/entity"

type AccountAppDTO struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Language    string `json:"language"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"fisrt_name"`
	LastName    string `json:"last_name"`
}

func NewAccountFromDTO(dto *AccountAppDTO) *entity.Account {
	return &entity.Account{
		Username:    dto.Username,
		Email:       dto.Email,
		Password:    dto.Password,
		PhoneNumber: dto.PhoneNumber,
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
	}
}
