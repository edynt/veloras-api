package service

import (
	"context"
	"math"
	"net/http"

	"github.com/edynt/chogiare/veloras-api/internal/order/application/service/dto"
	"github.com/edynt/chogiare/veloras-api/internal/order/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/order/domain/repository"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
)

type orderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req *dto.CreateOrderAppDTO, userID int32) (*dto.OrderAppDTO, error) {
	// Calculate totals (this would typically involve fetching product prices)
	// For now, we'll use a placeholder price
	subtotal := float64(req.Quantity) * 10.0 // Placeholder price

	tax := subtotal * 0.1 // 10% tax
	shipping := 5.0       // Fixed shipping cost
	discount := 0.0       // No discount for now
	total := subtotal + tax + shipping - discount

	order, err := s.orderRepo.CreateOrder(ctx, repository.CreateOrderParams{
		UserID:            userID,
		SellerID:          req.SellerID,
		ProductID:         req.ProductID,
		Quantity:          req.Quantity,
		Total:             total,
		Status:            "pending",
		PaymentMethod:     req.PaymentMethod,
		PaymentStatus:     "pending",
		ShippingAddressID: req.ShippingAddressID,
		Notes:             req.Notes,
	})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to create order", err)
	}

	// Order items are not supported in current schema - orders contain product info directly

	return s.convertToOrderAppDTO(order), nil
}

func (s *orderService) GetOrder(ctx context.Context, id string) (*dto.OrderAppDTO, error) {
	order, err := s.orderRepo.GetOrder(ctx, id)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Order not found", err)
	}

	return s.convertToOrderAppDTOWithDetails(order), nil
}

func (s *orderService) ListOrders(ctx context.Context, page, pageSize int32) (*dto.OrderListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	orders, err := s.orderRepo.ListOrders(ctx, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list orders", err)
	}

	total := int64(len(orders))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.OrderListAppDTO{
		Orders:     s.convertToOrderAppDTOsWithDetails(orders),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *orderService) ListOrdersByUser(ctx context.Context, userID int32, page, pageSize int32) (*dto.OrderListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	orders, err := s.orderRepo.ListOrdersByUser(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list user orders", err)
	}

	total := int64(len(orders))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.OrderListAppDTO{
		Orders:     s.convertToOrderAppDTOsWithDetails(orders),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *orderService) ListOrdersByStore(ctx context.Context, storeID string, page, pageSize int32) (*dto.OrderListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	orders, err := s.orderRepo.ListOrdersByStore(ctx, storeID, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list store orders", err)
	}

	total := int64(len(orders))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.OrderListAppDTO{
		Orders:     s.convertToOrderAppDTOsWithDetails(orders),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *orderService) UpdateOrder(ctx context.Context, id string, req *dto.UpdateOrderAppDTO) (*dto.OrderAppDTO, error) {
	// Get existing order to preserve unchanged fields
	existingOrder, err := s.orderRepo.GetOrder(ctx, id)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Order not found", err)
	}

	updateParams := repository.UpdateOrderParams{
		ID:            id,
		Status:        existingOrder.Status,
		PaymentStatus: existingOrder.PaymentStatus,
		Notes:         existingOrder.Notes,
	}

	// Update only provided fields
	if req.Status != nil {
		updateParams.Status = *req.Status
	}
	if req.PaymentStatus != nil {
		updateParams.PaymentStatus = *req.PaymentStatus
	}
	if req.Notes != nil {
		updateParams.Notes = req.Notes
	}

	order, err := s.orderRepo.UpdateOrder(ctx, updateParams)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to update order", err)
	}

	return s.convertToOrderAppDTO(order), nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, id string, status string) (*dto.OrderAppDTO, error) {
	order, err := s.orderRepo.UpdateOrderStatus(ctx, id, status)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to update order status", err)
	}

	return s.convertToOrderAppDTO(order), nil
}

func (s *orderService) UpdateOrderPaymentStatus(ctx context.Context, id string, paymentStatus string) (*dto.OrderAppDTO, error) {
	order, err := s.orderRepo.UpdateOrderPaymentStatus(ctx, id, paymentStatus)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to update payment status", err)
	}

	return s.convertToOrderAppDTO(order), nil
}

func (s *orderService) DeleteOrder(ctx context.Context, id string) error {
	err := s.orderRepo.DeleteOrder(ctx, id)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to delete order", err)
	}
	return nil
}

func (s *orderService) AddOrderItem(ctx context.Context, orderID string, req *dto.CreateOrderItemAppDTO) (*dto.OrderItemAppDTO, error) {
	item, err := s.orderRepo.CreateOrderItem(ctx, repository.CreateOrderItemParams{
		OrderID:      orderID,
		ProductID:    req.ProductID,
		ProductName:  "Product Name",      // This should come from product service
		ProductImage: "product-image.jpg", // This should come from product service
		Price:        10.0,                // Placeholder price
		Quantity:     req.Quantity,
		Subtotal:     float64(req.Quantity) * 10.0,
	})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to add order item", err)
	}

	return s.convertToOrderItemAppDTO(item), nil
}

func (s *orderService) UpdateOrderItem(ctx context.Context, itemID string, req *dto.CreateOrderItemAppDTO) (*dto.OrderItemAppDTO, error) {
	item, err := s.orderRepo.UpdateOrderItem(ctx, repository.UpdateOrderItemParams{
		ID:           itemID,
		ProductID:    req.ProductID,
		ProductName:  "Product Name",      // This should come from product service
		ProductImage: "product-image.jpg", // This should come from product service
		Price:        10.0,                // Placeholder price
		Quantity:     req.Quantity,
		Subtotal:     float64(req.Quantity) * 10.0,
	})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to update order item", err)
	}

	return s.convertToOrderItemAppDTO(item), nil
}

func (s *orderService) RemoveOrderItem(ctx context.Context, itemID string) error {
	err := s.orderRepo.DeleteOrderItem(ctx, itemID)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to remove order item", err)
	}
	return nil
}

func (s *orderService) GetOrderStats(ctx context.Context) (*dto.OrderStatsAppDTO, error) {
	stats, err := s.orderRepo.GetOrderStats(ctx)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to get order stats", err)
	}

	return &dto.OrderStatsAppDTO{
		TotalOrders:       stats.TotalOrders,
		PendingOrders:     stats.PendingOrders,
		ProcessingOrders:  stats.ProcessingOrders,
		ShippedOrders:     stats.ShippedOrders,
		DeliveredOrders:   stats.DeliveredOrders,
		CancelledOrders:   stats.CancelledOrders,
		TotalRevenue:      stats.TotalRevenue,
		AverageOrderValue: stats.AverageOrderValue,
	}, nil
}

func (s *orderService) GetOrderStatsByStore(ctx context.Context, storeID string) (*dto.OrderStatsAppDTO, error) {
	stats, err := s.orderRepo.GetOrderStatsByStore(ctx, storeID)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to get store order stats", err)
	}

	return &dto.OrderStatsAppDTO{
		TotalOrders:       stats.TotalOrders,
		PendingOrders:     stats.PendingOrders,
		ProcessingOrders:  stats.ProcessingOrders,
		ShippedOrders:     stats.ShippedOrders,
		DeliveredOrders:   stats.DeliveredOrders,
		CancelledOrders:   stats.CancelledOrders,
		TotalRevenue:      stats.TotalRevenue,
		AverageOrderValue: stats.AverageOrderValue,
	}, nil
}

func (s *orderService) GetOrderStatsByUser(ctx context.Context, userID int32) (*dto.OrderStatsAppDTO, error) {
	stats, err := s.orderRepo.GetOrderStatsByUser(ctx, userID)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to get user order stats", err)
	}

	return &dto.OrderStatsAppDTO{
		TotalOrders:       stats.TotalOrders,
		PendingOrders:     stats.PendingOrders,
		ProcessingOrders:  stats.ProcessingOrders,
		ShippedOrders:     stats.ShippedOrders,
		DeliveredOrders:   stats.DeliveredOrders,
		CancelledOrders:   stats.CancelledOrders,
		TotalRevenue:      stats.TotalRevenue,
		AverageOrderValue: stats.AverageOrderValue,
	}, nil
}

// Conversion methods
func (s *orderService) convertToOrderAppDTO(order *entity.Order) *dto.OrderAppDTO {
	return &dto.OrderAppDTO{
		ID:                order.ID,
		BuyerID:           order.BuyerID,
		SellerID:          order.SellerID,
		ProductID:         order.ProductID,
		Quantity:          order.Quantity,
		TotalAmount:       order.TotalAmount,
		Status:            order.Status,
		PaymentMethod:     order.PaymentMethod,
		PaymentStatus:     order.PaymentStatus,
		ShippingAddressID: order.ShippingAddressID,
		Notes:             order.Notes,
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
	}
}

func (s *orderService) convertToOrderAppDTOWithDetails(order *entity.OrderWithDetails) *dto.OrderAppDTO {
	result := &dto.OrderAppDTO{
		ID:                order.ID,
		BuyerID:           order.BuyerID,
		SellerID:          order.SellerID,
		ProductID:         order.ProductID,
		Quantity:          order.Quantity,
		TotalAmount:       order.TotalAmount,
		Status:            order.Status,
		PaymentMethod:     order.PaymentMethod,
		PaymentStatus:     order.PaymentStatus,
		ShippingAddressID: order.ShippingAddressID,
		Notes:             order.Notes,
		StoreName:         order.StoreName,
		StoreLogo:         order.StoreLogo,
		UserEmail:         order.UserEmail,
		UserName:          order.UserName,
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
	}

	// Convert order items
	result.Items = make([]dto.OrderItemAppDTO, len(order.Items))
	for i, item := range order.Items {
		result.Items[i] = dto.OrderItemAppDTO{
			ID:           item.ID,
			OrderID:      item.OrderID,
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			ProductImage: item.ProductImage,
			Price:        item.Price,
			Quantity:     item.Quantity,
			Subtotal:     item.Subtotal,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
	}

	return result
}

func (s *orderService) convertToOrderAppDTOsWithDetails(orders []entity.OrderWithDetails) []dto.OrderAppDTO {
	result := make([]dto.OrderAppDTO, len(orders))
	for i, order := range orders {
		result[i] = *s.convertToOrderAppDTOWithDetails(&order)
	}
	return result
}

func (s *orderService) convertToOrderItemAppDTO(item *entity.OrderItem) *dto.OrderItemAppDTO {
	return &dto.OrderItemAppDTO{
		ID:           item.ID,
		OrderID:      item.OrderID,
		ProductID:    item.ProductID,
		ProductName:  item.ProductName,
		ProductImage: item.ProductImage,
		Price:        item.Price,
		Quantity:     item.Quantity,
		Subtotal:     item.Subtotal,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	}
}
