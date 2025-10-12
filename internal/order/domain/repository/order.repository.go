package repository

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/order/domain/model/entity"
)

type OrderRepository interface {
	// Order CRUD operations
	CreateOrder(ctx context.Context, req CreateOrderParams) (*entity.Order, error)
	GetOrder(ctx context.Context, id string) (*entity.OrderWithDetails, error)
	ListOrders(ctx context.Context, limit, offset int32) ([]entity.OrderWithDetails, error)
	ListOrdersByUser(ctx context.Context, userID int32, limit, offset int32) ([]entity.OrderWithDetails, error)
	ListOrdersByStore(ctx context.Context, storeID string, limit, offset int32) ([]entity.OrderWithDetails, error)
	UpdateOrder(ctx context.Context, req UpdateOrderParams) (*entity.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status string) (*entity.Order, error)
	UpdateOrderPaymentStatus(ctx context.Context, id string, paymentStatus string) (*entity.Order, error)
	DeleteOrder(ctx context.Context, id string) error

	// Order item operations
	CreateOrderItem(ctx context.Context, req CreateOrderItemParams) (*entity.OrderItem, error)
	GetOrderItems(ctx context.Context, orderID string) ([]entity.OrderItem, error)
	UpdateOrderItem(ctx context.Context, req UpdateOrderItemParams) (*entity.OrderItem, error)
	DeleteOrderItem(ctx context.Context, id string) error

	// Statistics
	GetOrderStats(ctx context.Context) (*entity.OrderStats, error)
	GetOrderStatsByStore(ctx context.Context, storeID string) (*entity.OrderStats, error)
	GetOrderStatsByUser(ctx context.Context, userID int32) (*entity.OrderStats, error)
}

type CreateOrderParams struct {
	UserID          int32
	StoreID         string
	Status          string
	PaymentStatus   string
	PaymentMethod   string
	Subtotal        float64
	Tax             float64
	Shipping        float64
	Discount        float64
	Total           float64
	Currency        string
	ShippingAddress string
	BillingAddress  string
	Notes           *string
}

type UpdateOrderParams struct {
	ID              string
	Status          string
	PaymentStatus   string
	PaymentMethod   string
	Subtotal        float64
	Tax             float64
	Shipping        float64
	Discount        float64
	Total           float64
	Currency        string
	ShippingAddress string
	BillingAddress  string
	Notes           *string
}

type CreateOrderItemParams struct {
	OrderID      string
	ProductID    string
	ProductName  string
	ProductImage string
	Price        float64
	Quantity     int32
	Subtotal     float64
}

type UpdateOrderItemParams struct {
	ID           string
	ProductID    string
	ProductName  string
	ProductImage string
	Price        float64
	Quantity     int32
	Subtotal     float64
}
