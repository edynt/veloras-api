package dto

import "time"

// UserResponse represents the user information returned to the client
type UserResponse struct {
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

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int32          `json:"page"`
	PageSize   int32          `json:"page_size"`
	TotalPages int32          `json:"total_pages"`
}

// UpdateUserProfileRequest represents the request to update user profile
type UpdateUserProfileRequest struct {
	Username    string `json:"username" validate:"omitempty,min=3,max=50"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,min=10,max=20"`
	FirstName   string `json:"first_name" validate:"omitempty,min=1,max=50"`
	LastName    string `json:"last_name" validate:"omitempty,min=1,max=50"`
	Language    string `json:"language" validate:"omitempty,oneof=en vi"`
}
