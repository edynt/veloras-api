package dto

type UserRegisterReq struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Fullname    string `json:"fullname"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Country     string `json:"country"`
	Language    string `json:"language"`
}
