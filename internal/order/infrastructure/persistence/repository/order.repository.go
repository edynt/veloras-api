package repository

import (
	"context"
	"fmt"
	"math/big"

	"github.com/edynt/chogiare/veloras-api/internal/order/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/order/domain/repository"
	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderRepository struct {
	db      *pgxpool.Pool
	queries *gen.Queries
}

func NewOrderRepository(db *pgxpool.Pool) repository.OrderRepository {
	return &orderRepository{
		db:      db,
		queries: gen.New(db),
	}
}

// Helper functions for type conversion
func stringToUUID(s string) pgtype.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: u, Valid: true}
}

func float64ToNumeric(f float64) pgtype.Numeric {
	return pgtype.Numeric{Int: big.NewInt(int64(f * 100)), Valid: true, Exp: -2} // Assuming 2 decimal places
}

func stringPtrToText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func (r *orderRepository) CreateOrder(ctx context.Context, req repository.CreateOrderParams) (*entity.Order, error) {
	order, err := r.queries.CreateOrder(ctx, gen.CreateOrderParams{
		BuyerID:           req.UserID,
		SellerID:          req.SellerID,
		ProductID:         stringToUUID(req.ProductID),
		Quantity:          req.Quantity,
		TotalAmount:       float64ToNumeric(req.Total),
		Status:            req.Status,
		PaymentMethod:     req.PaymentMethod,
		PaymentStatus:     req.PaymentStatus,
		ShippingAddressID: stringToUUID(req.ShippingAddressID),
		Notes:             stringPtrToText(req.Notes),
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCOrder(order)
	return &result, nil
}

func (r *orderRepository) GetOrder(ctx context.Context, id string) (*entity.OrderWithDetails, error) {
	order, err := r.queries.GetOrder(ctx, stringToUUID(id))
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCOrderWithDetails(order)
	return &result, nil
}

func (r *orderRepository) ListOrders(ctx context.Context, limit, offset int32) ([]entity.OrderWithDetails, error) {
	orders, err := r.queries.ListOrders(ctx, gen.ListOrdersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.OrderWithDetails, len(orders))
	for i, o := range orders {
		result[i] = entity.FromSQLCOrderWithDetailsFromList(o)
	}

	return result, nil
}

func (r *orderRepository) ListOrdersByUser(ctx context.Context, userID int32, limit, offset int32) ([]entity.OrderWithDetails, error) {
	orders, err := r.queries.ListOrdersByUser(ctx, gen.ListOrdersByUserParams{
		BuyerID: userID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.OrderWithDetails, len(orders))
	for i, o := range orders {
		result[i] = entity.FromSQLCOrderWithDetailsFromListByUser(o)
	}

	return result, nil
}

func (r *orderRepository) ListOrdersByStore(ctx context.Context, storeID string, limit, offset int32) ([]entity.OrderWithDetails, error) {
	orders, err := r.queries.ListOrdersByStore(ctx, gen.ListOrdersByStoreParams{
		ID:     stringToUUID(storeID),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.OrderWithDetails, len(orders))
	for i, o := range orders {
		result[i] = entity.FromSQLCOrderWithDetailsFromListByStore(o)
	}

	return result, nil
}

func (r *orderRepository) UpdateOrder(ctx context.Context, req repository.UpdateOrderParams) (*entity.Order, error) {
	order, err := r.queries.UpdateOrder(ctx, gen.UpdateOrderParams{
		ID:            stringToUUID(req.ID),
		Status:        req.Status,
		PaymentStatus: req.PaymentStatus,
		Notes:         stringPtrToText(req.Notes),
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCOrder(order)
	return &result, nil
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, id string, status string) (*entity.Order, error) {
	order, err := r.queries.UpdateOrderStatus(ctx, gen.UpdateOrderStatusParams{
		ID:     stringToUUID(id),
		Status: status,
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCOrder(order)
	return &result, nil
}

func (r *orderRepository) UpdateOrderPaymentStatus(ctx context.Context, id string, paymentStatus string) (*entity.Order, error) {
	order, err := r.queries.UpdateOrderPaymentStatus(ctx, gen.UpdateOrderPaymentStatusParams{
		ID:            stringToUUID(id),
		PaymentStatus: paymentStatus,
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCOrder(order)
	return &result, nil
}

func (r *orderRepository) DeleteOrder(ctx context.Context, id string) error {
	return r.queries.DeleteOrder(ctx, stringToUUID(id))
}

func (r *orderRepository) CreateOrderItem(ctx context.Context, req repository.CreateOrderItemParams) (*entity.OrderItem, error) {
	// Order items are not supported in current schema - orders contain product info directly
	return nil, fmt.Errorf("order items not supported in current schema")
}

func (r *orderRepository) GetOrderItems(ctx context.Context, orderID string) ([]entity.OrderItem, error) {
	// Order items are not supported in current schema - orders contain product info directly
	return []entity.OrderItem{}, nil
}

func (r *orderRepository) UpdateOrderItem(ctx context.Context, req repository.UpdateOrderItemParams) (*entity.OrderItem, error) {
	// Order items are not supported in current schema - orders contain product info directly
	return nil, fmt.Errorf("order items not supported in current schema")
}

func (r *orderRepository) DeleteOrderItem(ctx context.Context, id string) error {
	// Order items are not supported in current schema - orders contain product info directly
	return fmt.Errorf("order items not supported in current schema")
}

func (r *orderRepository) GetOrderStats(ctx context.Context) (*entity.OrderStats, error) {
	stats, err := r.queries.GetOrderStats(ctx)
	if err != nil {
		return nil, err
	}

	totalRevenue := 0.0
	if stats.TotalRevenue != 0 {
		totalRevenue = float64(stats.TotalRevenue)
	}

	result := &entity.OrderStats{
		TotalOrders:       stats.TotalOrders,
		PendingOrders:     stats.PendingOrders,
		ProcessingOrders:  0, // Not available in current schema
		ShippedOrders:     stats.ShippedOrders,
		DeliveredOrders:   stats.DeliveredOrders,
		CancelledOrders:   stats.CancelledOrders,
		TotalRevenue:      totalRevenue,
		AverageOrderValue: 0.0, // Not available in current schema
	}

	return result, nil
}

func (r *orderRepository) GetOrderStatsByStore(ctx context.Context, storeID string) (*entity.OrderStats, error) {
	// GetOrderStatsByStore not available in current schema
	return &entity.OrderStats{}, nil
}

func (r *orderRepository) GetOrderStatsByUser(ctx context.Context, userID int32) (*entity.OrderStats, error) {
	// GetOrderStatsByUser not available in current schema
	return &entity.OrderStats{}, nil
}
