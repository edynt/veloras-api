package repository

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/store/domain/model/entity"
)

type StoreRepository interface {
	// Store CRUD operations
	CreateStore(ctx context.Context, req CreateStoreParams) (*entity.Store, error)
	GetStore(ctx context.Context, id string) (*entity.StoreWithDetails, error)
	GetStoreByUserID(ctx context.Context, userID int32) (*entity.StoreWithDetails, error)
	ListStores(ctx context.Context, limit, offset int32) ([]entity.StoreWithDetails, error)
	SearchStores(ctx context.Context, query string, limit, offset int32) ([]entity.StoreWithDetails, error)
	UpdateStore(ctx context.Context, req UpdateStoreParams) (*entity.Store, error)
	UpdateStoreRating(ctx context.Context, id string, rating float64, reviewCount int32) error
	UpdateStoreProductCount(ctx context.Context, id string, increment int32) error
	UpdateStoreFollowerCount(ctx context.Context, id string, increment int32) error
	DeleteStore(ctx context.Context, id string) error

	// Statistics
	GetStoreStats(ctx context.Context) (*entity.StoreStats, error)
}

type CreateStoreParams struct {
	UserID      int32
	Name        string
	Description *string
	Logo        *string
	Banner      *string
	Website     *string
	Phone       *string
	Email       *string
	Address     *string
	City        *string
	State       *string
	Country     *string
	PostalCode  *string
	IsVerified  bool
	IsActive    bool
}

type UpdateStoreParams struct {
	ID          string
	Name        string
	Description *string
	Logo        *string
	Banner      *string
	Website     *string
	Phone       *string
	Email       *string
	Address     *string
	City        *string
	State       *string
	Country     *string
	PostalCode  *string
	IsVerified  bool
	IsActive    bool
}
