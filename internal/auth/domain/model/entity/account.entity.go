package entity

type Account struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"fisrt_name"`
	LastName    string `json:"last_name"`
}
