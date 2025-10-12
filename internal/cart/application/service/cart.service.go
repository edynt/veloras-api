package service

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/cart/application/service/dto"
)

type CartService interface {
	// Cart operations
	GetCart(ctx context.Context, userID int32) (*dto.CartAppDTO, error)
	ClearCart(ctx context.Context, userID int32) error

	// Cart item operations
	AddCartItem(ctx context.Context, userID int32, req *dto.AddCartItemAppDTO) (*dto.CartItemAppDTO, error)
	UpdateCartItemQuantity(ctx context.Context, itemID string, req *dto.UpdateCartItemQuantityAppDTO) (*dto.CartItemAppDTO, error)
	RemoveCartItem(ctx context.Context, itemID string) error

	// Statistics
	GetCartStats(ctx context.Context, userID int32) (*dto.CartStatsAppDTO, error)
}
