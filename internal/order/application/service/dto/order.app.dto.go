package dto

import "time"

type OrderAppDTO struct {
	ID              string            `json:"id"`
	UserID          int32             `json:"userId"`
	StoreID         string            `json:"storeId"`
	Status          string            `json:"status"`
	PaymentStatus   string            `json:"paymentStatus"`
	PaymentMethod   string            `json:"paymentMethod"`
	Subtotal        float64           `json:"subtotal"`
	Tax             float64           `json:"tax"`
	Shipping        float64           `json:"shipping"`
	Discount        float64           `json:"discount"`
	Total           float64           `json:"total"`
	Currency        string            `json:"currency"`
	ShippingAddress string            `json:"shippingAddress"`
	BillingAddress  string            `json:"billingAddress"`
	Notes           *string           `json:"notes,omitempty"`
	StoreName       *string           `json:"storeName,omitempty"`
	StoreLogo       *string           `json:"storeLogo,omitempty"`
	UserEmail       *string           `json:"userEmail,omitempty"`
	UserName        *string           `json:"userName,omitempty"`
	Items           []OrderItemAppDTO `json:"items,omitempty"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
}

type OrderItemAppDTO struct {
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

type OrderListAppDTO struct {
	Orders     []OrderAppDTO `json:"orders"`
	Total      int64         `json:"total"`
	Page       int32         `json:"page"`
	PageSize   int32         `json:"pageSize"`
	TotalPages int32         `json:"totalPages"`
}

type CreateOrderAppDTO struct {
	StoreID         string                  `json:"storeId" binding:"required"`
	PaymentMethod   string                  `json:"paymentMethod" binding:"required"`
	ShippingAddress string                  `json:"shippingAddress" binding:"required"`
	BillingAddress  string                  `json:"billingAddress" binding:"required"`
	Notes           *string                 `json:"notes,omitempty"`
	Items           []CreateOrderItemAppDTO `json:"items" binding:"required,min=1"`
}

type CreateOrderItemAppDTO struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int32  `json:"quantity" binding:"required,min=1"`
}

type UpdateOrderAppDTO struct {
	Status          *string `json:"status,omitempty"`
	PaymentStatus   *string `json:"paymentStatus,omitempty"`
	PaymentMethod   *string `json:"paymentMethod,omitempty"`
	ShippingAddress *string `json:"shippingAddress,omitempty"`
	BillingAddress  *string `json:"billingAddress,omitempty"`
	Notes           *string `json:"notes,omitempty"`
}

type OrderStatsAppDTO struct {
	TotalOrders       int64   `json:"totalOrders"`
	PendingOrders     int64   `json:"pendingOrders"`
	ProcessingOrders  int64   `json:"processingOrders"`
	ShippedOrders     int64   `json:"shippedOrders"`
	DeliveredOrders   int64   `json:"deliveredOrders"`
	CancelledOrders   int64   `json:"cancelledOrders"`
	TotalRevenue      float64 `json:"totalRevenue"`
	AverageOrderValue float64 `json:"averageOrderValue"`
}
