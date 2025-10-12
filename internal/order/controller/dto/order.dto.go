package dto

import "time"

type CreateOrderRequest struct {
	SellerID          int32   `json:"sellerId" binding:"required"`
	ProductID         string  `json:"productId" binding:"required"`
	Quantity          int32   `json:"quantity" binding:"required,min=1"`
	PaymentMethod     string  `json:"paymentMethod" binding:"required"`
	ShippingAddressID string  `json:"shippingAddressId" binding:"required"`
	Notes             *string `json:"notes,omitempty"`
}

type CreateOrderItemRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int32  `json:"quantity" binding:"required,min=1"`
}

type UpdateOrderRequest struct {
	Status        *string `json:"status,omitempty"`
	PaymentStatus *string `json:"paymentStatus,omitempty"`
	Notes         *string `json:"notes,omitempty"`
}

type OrderResponse struct {
	ID                string              `json:"id"`
	BuyerID           int32               `json:"buyerId"`
	SellerID          int32               `json:"sellerId"`
	ProductID         string              `json:"productId"`
	Quantity          int32               `json:"quantity"`
	TotalAmount       float64             `json:"totalAmount"`
	Status            string              `json:"status"`
	PaymentMethod     string              `json:"paymentMethod"`
	PaymentStatus     string              `json:"paymentStatus"`
	ShippingAddressID string              `json:"shippingAddressId"`
	Notes             *string             `json:"notes,omitempty"`
	StoreName         *string             `json:"storeName,omitempty"`
	StoreLogo         *string             `json:"storeLogo,omitempty"`
	UserEmail         *string             `json:"userEmail,omitempty"`
	UserName          *string             `json:"userName,omitempty"`
	Items             []OrderItemResponse `json:"items,omitempty"`
	CreatedAt         time.Time           `json:"createdAt"`
	UpdatedAt         time.Time           `json:"updatedAt"`
}

type OrderItemResponse struct {
	ID           string    `json:"id"`
	OrderID      string    `json:"orderId"`
	ProductID    string    `json:"productId"`
	ProductName  string    `json:"productName"`
	ProductImage string    `json:"productImage"`
	Price        float64   `json:"price"`
	Quantity     int32     `json:"quantity"`
	Subtotal     float64   `json:"subtotal"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type OrderListResponse struct {
	Orders     []OrderResponse `json:"orders"`
	Total      int64           `json:"total"`
	Page       int32           `json:"page"`
	PageSize   int32           `json:"pageSize"`
	TotalPages int32           `json:"totalPages"`
}

type OrderStatsResponse struct {
	TotalOrders       int64   `json:"totalOrders"`
	PendingOrders     int64   `json:"pendingOrders"`
	ProcessingOrders  int64   `json:"processingOrders"`
	ShippedOrders     int64   `json:"shippedOrders"`
	DeliveredOrders   int64   `json:"deliveredOrders"`
	CancelledOrders   int64   `json:"cancelledOrders"`
	TotalRevenue      float64 `json:"totalRevenue"`
	AverageOrderValue float64 `json:"averageOrderValue"`
}
