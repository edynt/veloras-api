package entity

import "time"

type User struct {
	ID          int32
	Email       string
	Username    string
	IsVerified  bool
	PhoneNumber string
	FirstName   string
	LastName    string
	Status      int32
	Language    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserList struct {
	Users      []User
	Total      int64
	Page       int32
	PageSize   int32
	TotalPages int32
}
