package dto

import "time"

type StoreAppDTO struct {
	ID            string    `json:"id"`
	UserID        int32     `json:"userId"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	Logo          *string   `json:"logo,omitempty"`
	Banner        *string   `json:"banner,omitempty"`
	Website       *string   `json:"website,omitempty"`
	Phone         *string   `json:"phone,omitempty"`
	Email         *string   `json:"email,omitempty"`
	Address       *string   `json:"address,omitempty"`
	City          *string   `json:"city,omitempty"`
	State         *string   `json:"state,omitempty"`
	Country       *string   `json:"country,omitempty"`
	PostalCode    *string   `json:"postalCode,omitempty"`
	Rating        float64   `json:"rating"`
	ReviewCount   int32     `json:"reviewCount"`
	ProductCount  int32     `json:"productCount"`
	FollowerCount int32     `json:"followerCount"`
	IsVerified    bool      `json:"isVerified"`
	IsActive      bool      `json:"isActive"`
	UserName      *string   `json:"userName,omitempty"`
	UserEmail     *string   `json:"userEmail,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type StoreListAppDTO struct {
	Stores     []StoreAppDTO `json:"stores"`
	Total      int64         `json:"total"`
	Page       int32         `json:"page"`
	PageSize   int32         `json:"pageSize"`
	TotalPages int32         `json:"totalPages"`
}

type CreateStoreAppDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	Logo        *string `json:"logo,omitempty"`
	Banner      *string `json:"banner,omitempty"`
	Website     *string `json:"website,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Email       *string `json:"email,omitempty"`
	Address     *string `json:"address,omitempty"`
	City        *string `json:"city,omitempty"`
	State       *string `json:"state,omitempty"`
	Country     *string `json:"country,omitempty"`
	PostalCode  *string `json:"postalCode,omitempty"`
	IsVerified  bool    `json:"isVerified"`
	IsActive    bool    `json:"isActive"`
}

type UpdateStoreAppDTO struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Logo        *string `json:"logo,omitempty"`
	Banner      *string `json:"banner,omitempty"`
	Website     *string `json:"website,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Email       *string `json:"email,omitempty"`
	Address     *string `json:"address,omitempty"`
	City        *string `json:"city,omitempty"`
	State       *string `json:"state,omitempty"`
	Country     *string `json:"country,omitempty"`
	PostalCode  *string `json:"postalCode,omitempty"`
	IsVerified  *bool   `json:"isVerified,omitempty"`
	IsActive    *bool   `json:"isActive,omitempty"`
}

type StoreStatsAppDTO struct {
	TotalStores    int64   `json:"totalStores"`
	ActiveStores   int64   `json:"activeStores"`
	VerifiedStores int64   `json:"verifiedStores"`
	AverageRating  float64 `json:"averageRating"`
	TotalProducts  int64   `json:"totalProducts"`
	TotalFollowers int64   `json:"totalFollowers"`
}
