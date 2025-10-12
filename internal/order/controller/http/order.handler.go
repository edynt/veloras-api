package http

import (
	"net/http"
	"strconv"

	"github.com/edynt/chogiare/veloras-api/internal/order/application/service"
	appDto "github.com/edynt/chogiare/veloras-api/internal/order/application/service/dto"
	ctlDto "github.com/edynt/chogiare/veloras-api/internal/order/controller/dto"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/edynt/chogiare/veloras-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateOrder
// @Summary Create a new order
// @Description Create a new order with items
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body ctlDto.CreateOrderRequest true "Order data"
// @Success 201 {object} ctlDto.OrderResponse "Order created successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", err.Error())
	}

	// Validate the request
	validation, exists := ctx.Get("validation")
	if !exists {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Validation middleware not found", "")
	}
	if apiErr := utils.ValidateStruct(req, validation.(*validator.Validate)); apiErr != nil {
		return nil, apiErr
	}

	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	// Convert to application DTO
	appReq := &appDto.CreateOrderAppDTO{
		StoreID:         req.StoreID,
		PaymentMethod:   req.PaymentMethod,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		Notes:           req.Notes,
		Items:           make([]appDto.CreateOrderItemAppDTO, len(req.Items)),
	}

	for i, item := range req.Items {
		appReq.Items[i] = appDto.CreateOrderItemAppDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	order, err := h.service.CreateOrder(ctx, appReq, userID)
	if err != nil {
		return nil, err
	}

	return h.convertToOrderResponse(order), nil
}

// GetOrder
// @Summary Get order by ID
// @Description Get order details by ID
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} ctlDto.OrderResponse "Order details"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Order not found"
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Order ID is required")
	}

	order, err := h.service.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	return h.convertToOrderResponse(order), nil
}

// ListOrders
// @Summary List all orders
// @Description Get a paginated list of all orders
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} ctlDto.OrderListResponse "List of orders"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /orders [get]
func (h *OrderHandler) ListOrders(ctx *gin.Context) (res interface{}, err error) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	orders, err := h.service.ListOrders(ctx, int32(page), int32(pageSize))
	if err != nil {
		return nil, err
	}

	return h.convertToOrderListResponse(orders), nil
}

// ListUserOrders
// @Summary List user's orders
// @Description Get a paginated list of orders for the current user
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} ctlDto.OrderListResponse "List of user orders"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /orders/my [get]
func (h *OrderHandler) ListUserOrders(ctx *gin.Context) (res interface{}, err error) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	orders, err := h.service.ListOrdersByUser(ctx, userID, int32(page), int32(pageSize))
	if err != nil {
		return nil, err
	}

	return h.convertToOrderListResponse(orders), nil
}

// ListStoreOrders
// @Summary List store's orders
// @Description Get a paginated list of orders for a specific store
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param store_id path string true "Store ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} ctlDto.OrderListResponse "List of store orders"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /orders/store/{store_id} [get]
func (h *OrderHandler) ListStoreOrders(ctx *gin.Context) (res interface{}, err error) {
	storeID := ctx.Param("store_id")
	if storeID == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Store ID is required")
	}

	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	orders, err := h.service.ListOrdersByStore(ctx, storeID, int32(page), int32(pageSize))
	if err != nil {
		return nil, err
	}

	return h.convertToOrderListResponse(orders), nil
}

// UpdateOrder
// @Summary Update order
// @Description Update order details
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Param order body ctlDto.UpdateOrderRequest true "Order update data"
// @Success 200 {object} ctlDto.OrderResponse "Order updated successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Order not found"
// @Router /orders/{id} [put]
func (h *OrderHandler) UpdateOrder(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Order ID is required")
	}

	var req ctlDto.UpdateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", err.Error())
	}

	// Convert to application DTO
	appReq := &appDto.UpdateOrderAppDTO{
		Status:          req.Status,
		PaymentStatus:   req.PaymentStatus,
		PaymentMethod:   req.PaymentMethod,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		Notes:           req.Notes,
	}

	order, err := h.service.UpdateOrder(ctx, id, appReq)
	if err != nil {
		return nil, err
	}

	return h.convertToOrderResponse(order), nil
}

// UpdateOrderStatus
// @Summary Update order status
// @Description Update the status of an order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Param status body string true "New status"
// @Success 200 {object} ctlDto.OrderResponse "Order status updated successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Order not found"
// @Router /orders/{id}/status [patch]
func (h *OrderHandler) UpdateOrderStatus(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Order ID is required")
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", err.Error())
	}

	order, err := h.service.UpdateOrderStatus(ctx, id, req.Status)
	if err != nil {
		return nil, err
	}

	return h.convertToOrderResponse(order), nil
}

// UpdateOrderPaymentStatus
// @Summary Update order payment status
// @Description Update the payment status of an order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Param status body string true "New payment status"
// @Success 200 {object} ctlDto.OrderResponse "Order payment status updated successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Order not found"
// @Router /orders/{id}/payment-status [patch]
func (h *OrderHandler) UpdateOrderPaymentStatus(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Order ID is required")
	}

	var req struct {
		PaymentStatus string `json:"paymentStatus" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", err.Error())
	}

	order, err := h.service.UpdateOrderPaymentStatus(ctx, id, req.PaymentStatus)
	if err != nil {
		return nil, err
	}

	return h.convertToOrderResponse(order), nil
}

// DeleteOrder
// @Summary Delete order
// @Description Delete an order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 204 "Order deleted successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Order not found"
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Order ID is required")
	}

	err = h.service.DeleteOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetOrderStats
// @Summary Get order statistics
// @Description Get overall order statistics
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ctlDto.OrderStatsResponse "Order statistics"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /orders/stats [get]
func (h *OrderHandler) GetOrderStats(ctx *gin.Context) (res interface{}, err error) {
	stats, err := h.service.GetOrderStats(ctx)
	if err != nil {
		return nil, err
	}

	return h.convertToOrderStatsResponse(stats), nil
}

// GetStoreOrderStats
// @Summary Get store order statistics
// @Description Get order statistics for a specific store
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param store_id path string true "Store ID"
// @Success 200 {object} ctlDto.OrderStatsResponse "Store order statistics"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /orders/stats/store/{store_id} [get]
func (h *OrderHandler) GetStoreOrderStats(ctx *gin.Context) (res interface{}, err error) {
	storeID := ctx.Param("store_id")
	if storeID == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Store ID is required")
	}

	stats, err := h.service.GetOrderStatsByStore(ctx, storeID)
	if err != nil {
		return nil, err
	}

	return h.convertToOrderStatsResponse(stats), nil
}

// GetUserOrderStats
// @Summary Get user order statistics
// @Description Get order statistics for the current user
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ctlDto.OrderStatsResponse "User order statistics"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /orders/stats/my [get]
func (h *OrderHandler) GetUserOrderStats(ctx *gin.Context) (res interface{}, err error) {
	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	stats, err := h.service.GetOrderStatsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return h.convertToOrderStatsResponse(stats), nil
}

// Conversion methods
func (h *OrderHandler) convertToOrderResponse(order *appDto.OrderAppDTO) *ctlDto.OrderResponse {
	result := &ctlDto.OrderResponse{
		ID:              order.ID,
		UserID:          order.UserID,
		StoreID:         order.StoreID,
		Status:          order.Status,
		PaymentStatus:   order.PaymentStatus,
		PaymentMethod:   order.PaymentMethod,
		Subtotal:        order.Subtotal,
		Tax:             order.Tax,
		Shipping:        order.Shipping,
		Discount:        order.Discount,
		Total:           order.Total,
		Currency:        order.Currency,
		ShippingAddress: order.ShippingAddress,
		BillingAddress:  order.BillingAddress,
		Notes:           order.Notes,
		StoreName:       order.StoreName,
		StoreLogo:       order.StoreLogo,
		UserEmail:       order.UserEmail,
		UserName:        order.UserName,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}

	// Convert order items
	result.Items = make([]ctlDto.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		result.Items[i] = ctlDto.OrderItemResponse{
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

func (h *OrderHandler) convertToOrderListResponse(orders *appDto.OrderListAppDTO) *ctlDto.OrderListResponse {
	result := &ctlDto.OrderListResponse{
		Total:      orders.Total,
		Page:       orders.Page,
		PageSize:   orders.PageSize,
		TotalPages: orders.TotalPages,
	}

	result.Orders = make([]ctlDto.OrderResponse, len(orders.Orders))
	for i, order := range orders.Orders {
		result.Orders[i] = *h.convertToOrderResponse(&order)
	}

	return result
}

func (h *OrderHandler) convertToOrderStatsResponse(stats *appDto.OrderStatsAppDTO) *ctlDto.OrderStatsResponse {
	return &ctlDto.OrderStatsResponse{
		TotalOrders:       stats.TotalOrders,
		PendingOrders:     stats.PendingOrders,
		ProcessingOrders:  stats.ProcessingOrders,
		ShippedOrders:     stats.ShippedOrders,
		DeliveredOrders:   stats.DeliveredOrders,
		CancelledOrders:   stats.CancelledOrders,
		TotalRevenue:      stats.TotalRevenue,
		AverageOrderValue: stats.AverageOrderValue,
	}
}
