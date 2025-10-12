package repository

import (
	"context"
	"math/big"

	"github.com/edynt/chogiare/veloras-api/internal/cart/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/cart/domain/repository"
	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type cartRepository struct {
	db      *pgxpool.Pool
	queries *gen.Queries
}

func NewCartRepository(db *pgxpool.Pool) repository.CartRepository {
	return &cartRepository{
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

func (r *cartRepository) GetOrCreateCart(ctx context.Context, userID int32) (*entity.Cart, error) {
	_, err := r.queries.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCCart(userID)
	return &result, nil
}

func (r *cartRepository) GetCart(ctx context.Context, userID int32) (*entity.CartWithDetails, error) {
	_, err := r.queries.GetCartWithDetails(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get cart items separately since the query returns aggregated data
	items, err := r.queries.ListCartItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCCartWithDetails(userID, items)
	return &result, nil
}

func (r *cartRepository) ClearCart(ctx context.Context, userID int32) error {
	return r.queries.ClearCart(ctx, userID)
}

func (r *cartRepository) DeleteCart(ctx context.Context, userID int32) error {
	return r.queries.DeleteCart(ctx, userID)
}

func (r *cartRepository) AddCartItem(ctx context.Context, userID int32, req repository.AddCartItemParams) (*entity.CartItem, error) {
	item, err := r.queries.AddCartItem(ctx, gen.AddCartItemParams{
		UserID:    userID,
		ProductID: stringToUUID(req.ProductID),
		Quantity:  req.Quantity,
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCCartItem(item)
	return &result, nil
}

func (r *cartRepository) UpdateCartItemQuantity(ctx context.Context, userID int32, req repository.UpdateCartItemQuantityParams) (*entity.CartItem, error) {
	item, err := r.queries.UpdateCartItemQuantity(ctx, gen.UpdateCartItemQuantityParams{
		UserID:    userID,
		ProductID: stringToUUID(req.ItemID),
		Quantity:  req.Quantity,
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCCartItem(item)
	return &result, nil
}

func (r *cartRepository) RemoveCartItem(ctx context.Context, userID int32, itemID string) error {
	return r.queries.RemoveCartItem(ctx, gen.RemoveCartItemParams{
		UserID:    userID,
		ProductID: stringToUUID(itemID),
	})
}

func (r *cartRepository) GetCartItems(ctx context.Context, userID int32) ([]entity.CartItemWithDetails, error) {
	items, err := r.queries.GetCartItemsWithDetails(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]entity.CartItemWithDetails, len(items))
	for i, item := range items {
		result[i] = entity.FromSQLCCartItemWithDetailsFromGet(item, userID)
	}

	return result, nil
}

func (r *cartRepository) GetCartItem(ctx context.Context, userID int32, itemID string) (*entity.CartItemWithDetails, error) {
	item, err := r.queries.GetCartItemWithDetails(ctx, stringToUUID(itemID))
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCCartItemWithDetailsFromGetSingle(item, userID)
	return &result, nil
}

func (r *cartRepository) GetCartStats(ctx context.Context, userID int32) (*entity.CartStats, error) {
	stats, err := r.queries.GetCartStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalValue := 0.0
	if stats.TotalValue != nil {
		if val, ok := stats.TotalValue.(float64); ok {
			totalValue = val
		}
	}

	result := &entity.CartStats{
		TotalItems:     stats.TotalItems,
		TotalValue:     totalValue,
		UniqueProducts: stats.UniqueProducts,
	}

	return result, nil
}
