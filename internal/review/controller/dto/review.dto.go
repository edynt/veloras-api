package dto

import "time"

type CreateReviewRequest struct {
	ProductID  string   `json:"productId" binding:"required"`
	SellerID   int32    `json:"sellerId" binding:"required"`
	Rating     int32    `json:"rating" binding:"required,min=1,max=5"`
	Comment    *string  `json:"comment,omitempty"`
	Images     []string `json:"images,omitempty"`
	IsVerified bool     `json:"isVerified"`
}

type UpdateReviewRequest struct {
	Rating     int32    `json:"rating" binding:"required,min=1,max=5"`
	Comment    *string  `json:"comment,omitempty"`
	Images     []string `json:"images,omitempty"`
	IsVerified bool     `json:"isVerified"`
}

type ReviewResponse struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"productId"`
	BuyerID      int32     `json:"buyerId"`
	SellerID     int32     `json:"sellerId"`
	Rating       int32     `json:"rating"`
	Comment      *string   `json:"comment,omitempty"`
	Images       []string  `json:"images"`
	IsVerified   bool      `json:"isVerified"`
	UserName     *string   `json:"userName,omitempty"`
	UserEmail    *string   `json:"userEmail,omitempty"`
	UserAvatar   *string   `json:"userAvatar,omitempty"`
	ProductName  *string   `json:"productName,omitempty"`
	ProductImage *string   `json:"productImage,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type ReviewListResponse struct {
	Reviews    []ReviewResponse `json:"reviews"`
	Total      int64            `json:"total"`
	Page       int32            `json:"page"`
	PageSize   int32            `json:"pageSize"`
	TotalPages int32            `json:"totalPages"`
}

type ReviewStatsResponse struct {
	TotalReviews    int64           `json:"totalReviews"`
	AverageRating   float64         `json:"averageRating"`
	RatingCounts    map[int32]int64 `json:"ratingCounts"`
	VerifiedReviews int64           `json:"verifiedReviews"`
}
