package service

import (
	"context"
	"net/http"

	"github.com/edynt/chogiare/veloras-api/internal/cart/application/service/dto"
	"github.com/edynt/chogiare/veloras-api/internal/cart/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/cart/domain/repository"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
)

type cartService struct {
	cartRepo repository.CartRepository
}

func NewCartService(cartRepo repository.CartRepository) CartService {
	return &cartService{
		cartRepo: cartRepo,
	}
}

func (s *cartService) GetCart(ctx context.Context, userID int32) (*dto.CartAppDTO, error) {
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to get cart", err)
	}

	return s.convertToCartAppDTO(cart), nil
}

func (s *cartService) ClearCart(ctx context.Context, userID int32) error {
	err := s.cartRepo.ClearCart(ctx, userID)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to clear cart", err)
	}
	return nil
}

func (s *cartService) AddCartItem(ctx context.Context, userID int32, req *dto.AddCartItemAppDTO) (*dto.CartItemAppDTO, error) {
	// Get or create cart for user
	_, err := s.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to get cart", err)
	}

	// Add item to cart
	item, err := s.cartRepo.AddCartItem(ctx, userID, repository.AddCartItemParams{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to add item to cart", err)
	}

	return s.convertToCartItemAppDTO(item), nil
}

func (s *cartService) UpdateCartItemQuantity(ctx context.Context, userID int32, itemID string, req *dto.UpdateCartItemQuantityAppDTO) (*dto.CartItemAppDTO, error) {
	item, err := s.cartRepo.UpdateCartItemQuantity(ctx, userID, repository.UpdateCartItemQuantityParams{
		ItemID:   itemID,
		Quantity: req.Quantity,
	})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to update cart item quantity", err)
	}

	return s.convertToCartItemAppDTO(item), nil
}

func (s *cartService) RemoveCartItem(ctx context.Context, userID int32, itemID string) error {
	err := s.cartRepo.RemoveCartItem(ctx, userID, itemID)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to remove cart item", err)
	}
	return nil
}

func (s *cartService) GetCartStats(ctx context.Context, userID int32) (*dto.CartStatsAppDTO, error) {
	stats, err := s.cartRepo.GetCartStats(ctx, userID)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to get cart stats", err)
	}

	return &dto.CartStatsAppDTO{
		TotalItems:     stats.TotalItems,
		TotalValue:     stats.TotalValue,
		UniqueProducts: stats.UniqueProducts,
	}, nil
}

// Conversion methods
func (s *cartService) convertToCartAppDTO(cart *entity.CartWithDetails) *dto.CartAppDTO {
	result := &dto.CartAppDTO{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}

	// Convert cart items
	result.Items = make([]dto.CartItemAppDTO, len(cart.Items))
	for i, item := range cart.Items {
		result.Items[i] = dto.CartItemAppDTO{
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

func (s *cartService) convertToCartItemAppDTO(item *entity.CartItem) *dto.CartItemAppDTO {
	return &dto.CartItemAppDTO{
		ID:        item.ID,
		CartID:    item.CartID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		Price:     item.Price,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
