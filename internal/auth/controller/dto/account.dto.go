package dto

type UserRegisterReq struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Country     string `json:"country"`
	Language    string `json:"language"`
}

type UserLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ChangePasswordReq struct {
	CurrentPassword string `json:"current_password" validate:"required,min=6"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,eqfield=NewPassword"`
}
