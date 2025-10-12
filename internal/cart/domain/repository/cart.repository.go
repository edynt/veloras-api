package repository

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/cart/domain/model/entity"
)

type CartRepository interface {
	// Cart operations
	GetOrCreateCart(ctx context.Context, userID int32) (*entity.Cart, error)
	GetCart(ctx context.Context, userID int32) (*entity.CartWithDetails, error)
	ClearCart(ctx context.Context, userID int32) error
	DeleteCart(ctx context.Context, userID int32) error

	// Cart item operations
	AddCartItem(ctx context.Context, userID int32, req AddCartItemParams) (*entity.CartItem, error)
	UpdateCartItemQuantity(ctx context.Context, userID int32, req UpdateCartItemQuantityParams) (*entity.CartItem, error)
	RemoveCartItem(ctx context.Context, userID int32, itemID string) error
	GetCartItems(ctx context.Context, userID int32) ([]entity.CartItemWithDetails, error)
	GetCartItem(ctx context.Context, userID int32, itemID string) (*entity.CartItemWithDetails, error)

	// Statistics
	GetCartStats(ctx context.Context, userID int32) (*entity.CartStats, error)
}

type AddCartItemParams struct {
	ProductID string
	Quantity  int32
}

type UpdateCartItemQuantityParams struct {
	ItemID   string
	Quantity int32
}
