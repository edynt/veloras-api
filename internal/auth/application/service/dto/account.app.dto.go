package dto

type AccountAppDTO struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Language    string `json:"language"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"fisrt_name"`
	LastName    string `json:"last_name"`
}
