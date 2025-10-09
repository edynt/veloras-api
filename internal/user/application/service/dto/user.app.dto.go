package dto

import "time"

type UserAppDTO struct {
	ID          int32     `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	IsVerified  bool      `json:"is_verified"`
	PhoneNumber string    `json:"phone_number"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Status      int32     `json:"status"`
	Language    string    `json:"language"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserListAppDTO struct {
	Users      []UserAppDTO `json:"users"`
	Total      int64        `json:"total"`
	Page       int32        `json:"page"`
	PageSize   int32        `json:"page_size"`
	TotalPages int32        `json:"total_pages"`
}

type UpdateUserProfileAppDTO struct {
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Language    string `json:"language"`
}
