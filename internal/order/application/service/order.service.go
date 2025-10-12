package service

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/order/application/service/dto"
)

type OrderService interface {
	// Order operations
	CreateOrder(ctx context.Context, req *dto.CreateOrderAppDTO, userID int32) (*dto.OrderAppDTO, error)
	GetOrder(ctx context.Context, id string) (*dto.OrderAppDTO, error)
	ListOrders(ctx context.Context, page, pageSize int32) (*dto.OrderListAppDTO, error)
	ListOrdersByUser(ctx context.Context, userID int32, page, pageSize int32) (*dto.OrderListAppDTO, error)
	ListOrdersByStore(ctx context.Context, storeID string, page, pageSize int32) (*dto.OrderListAppDTO, error)
	UpdateOrder(ctx context.Context, id string, req *dto.UpdateOrderAppDTO) (*dto.OrderAppDTO, error)
	UpdateOrderStatus(ctx context.Context, id string, status string) (*dto.OrderAppDTO, error)
	UpdateOrderPaymentStatus(ctx context.Context, id string, paymentStatus string) (*dto.OrderAppDTO, error)
	DeleteOrder(ctx context.Context, id string) error

	// Order item operations
	AddOrderItem(ctx context.Context, orderID string, req *dto.CreateOrderItemAppDTO) (*dto.OrderItemAppDTO, error)
	UpdateOrderItem(ctx context.Context, itemID string, req *dto.CreateOrderItemAppDTO) (*dto.OrderItemAppDTO, error)
	RemoveOrderItem(ctx context.Context, itemID string) error

	// Statistics
	GetOrderStats(ctx context.Context) (*dto.OrderStatsAppDTO, error)
	GetOrderStatsByStore(ctx context.Context, storeID string) (*dto.OrderStatsAppDTO, error)
	GetOrderStatsByUser(ctx context.Context, userID int32) (*dto.OrderStatsAppDTO, error)
}
