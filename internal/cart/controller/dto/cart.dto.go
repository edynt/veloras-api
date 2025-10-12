package dto

import "time"

type AddCartItemRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int32  `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemQuantityRequest struct {
	Quantity int32 `json:"quantity" binding:"required,min=1"`
}

type CartResponse struct {
	ID        string             `json:"id"`
	UserID    int32              `json:"userId"`
	Items     []CartItemResponse `json:"items"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

type CartItemResponse struct {
	ID            string    `json:"id"`
	CartID        string    `json:"cartId"`
	ProductID     string    `json:"productId"`
	Quantity      int32     `json:"quantity"`
	Price         float64   `json:"price"`
	ProductName   string    `json:"productName"`
	ProductImage  string    `json:"productImage"`
	ProductPrice  float64   `json:"productPrice"`
	ProductStock  int32     `json:"productStock"`
	ProductStatus string    `json:"productStatus"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type CartStatsResponse struct {
	TotalItems     int64   `json:"totalItems"`
	TotalValue     float64 `json:"totalValue"`
	UniqueProducts int64   `json:"uniqueProducts"`
}
