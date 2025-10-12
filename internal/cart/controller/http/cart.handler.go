package http

import (
	"net/http"

	"github.com/edynt/chogiare/veloras-api/internal/cart/application/service"
	appDto "github.com/edynt/chogiare/veloras-api/internal/cart/application/service/dto"
	"github.com/edynt/chogiare/veloras-api/internal/cart/controller/dto"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/edynt/chogiare/veloras-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type CartHandler struct {
	service service.CartService
}

func NewCartHandler(service service.CartService) *CartHandler {
	return &CartHandler{service: service}
}

// GetCart
// @Summary Get user's cart
// @Description Get the current user's shopping cart with all items
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.CartResponse "User's cart"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /cart [get]
func (h *CartHandler) GetCart(ctx *gin.Context) (res interface{}, err error) {
	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	cart, err := h.service.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	return h.convertToCartResponse(cart), nil
}

// AddCartItem
// @Summary Add item to cart
// @Description Add a product to the user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body dto.AddCartItemRequest true "Cart item data"
// @Success 201 {object} dto.CartItemResponse "Item added to cart"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /cart/items [post]
func (h *CartHandler) AddCartItem(ctx *gin.Context) (res interface{}, err error) {
	var req dto.AddCartItemRequest
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
	appReq := &appDto.AddCartItemAppDTO{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	item, err := h.service.AddCartItem(ctx, userID, appReq)
	if err != nil {
		return nil, err
	}

	return h.convertToCartItemResponse(item), nil
}

// UpdateCartItemQuantity
// @Summary Update cart item quantity
// @Description Update the quantity of an item in the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Cart item ID"
// @Param item body dto.UpdateCartItemQuantityRequest true "Quantity update data"
// @Success 200 {object} dto.CartItemResponse "Cart item updated"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Cart item not found"
// @Router /cart/items/{id} [put]
func (h *CartHandler) UpdateCartItemQuantity(ctx *gin.Context) (res interface{}, err error) {
	itemID := ctx.Param("id")
	if itemID == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Item ID is required")
	}

	var req dto.UpdateCartItemQuantityRequest
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
	appReq := &appDto.UpdateCartItemQuantityAppDTO{
		Quantity: req.Quantity,
	}

	item, err := h.service.UpdateCartItemQuantity(ctx, userID, itemID, appReq)
	if err != nil {
		return nil, err
	}

	return h.convertToCartItemResponse(item), nil
}

// RemoveCartItem
// @Summary Remove item from cart
// @Description Remove an item from the user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Cart item ID"
// @Success 204 "Item removed from cart"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Cart item not found"
// @Router /cart/items/{id} [delete]
func (h *CartHandler) RemoveCartItem(ctx *gin.Context) (res interface{}, err error) {
	itemID := ctx.Param("id")
	if itemID == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Item ID is required")
	}

	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	err = h.service.RemoveCartItem(ctx, userID, itemID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ClearCart
// @Summary Clear cart
// @Description Remove all items from the user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 204 "Cart cleared"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /cart/clear [post]
func (h *CartHandler) ClearCart(ctx *gin.Context) (res interface{}, err error) {
	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	err = h.service.ClearCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetCartStats
// @Summary Get cart statistics
// @Description Get statistics about the user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.CartStatsResponse "Cart statistics"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /cart/stats [get]
func (h *CartHandler) GetCartStats(ctx *gin.Context) (res interface{}, err error) {
	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	stats, err := h.service.GetCartStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return h.convertToCartStatsResponse(stats), nil
}

// Conversion methods
func (h *CartHandler) convertToCartResponse(cart *appDto.CartAppDTO) *dto.CartResponse {
	result := &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}

	// Convert cart items
	result.Items = make([]dto.CartItemResponse, len(cart.Items))
	for i, item := range cart.Items {
		result.Items[i] = dto.CartItemResponse{
			ID:            item.ID,
			CartID:        item.CartID,
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			Price:         item.Price,
			ProductName:   item.ProductName,
			ProductImage:  item.ProductImage,
			ProductPrice:  item.ProductPrice,
			ProductStock:  item.ProductStock,
			ProductStatus: item.ProductStatus,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
		}
	}

	return result
}

func (h *CartHandler) convertToCartItemResponse(item *appDto.CartItemAppDTO) *dto.CartItemResponse {
	return &dto.CartItemResponse{
		ID:            item.ID,
		CartID:        item.CartID,
		ProductID:     item.ProductID,
		Quantity:      item.Quantity,
		Price:         item.Price,
		ProductName:   item.ProductName,
		ProductImage:  item.ProductImage,
		ProductPrice:  item.ProductPrice,
		ProductStock:  item.ProductStock,
		ProductStatus: item.ProductStatus,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}
}

func (h *CartHandler) convertToCartStatsResponse(stats *appDto.CartStatsAppDTO) *dto.CartStatsResponse {
	return &dto.CartStatsResponse{
		TotalItems:     stats.TotalItems,
		TotalValue:     stats.TotalValue,
		UniqueProducts: stats.UniqueProducts,
	}
}
