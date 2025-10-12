package service

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/store/application/service/dto"
)

type StoreService interface {
	// Store operations
	CreateStore(ctx context.Context, req *dto.CreateStoreAppDTO, userID int32) (*dto.StoreAppDTO, error)
	GetStore(ctx context.Context, id string) (*dto.StoreAppDTO, error)
	GetStoreByUserID(ctx context.Context, userID int32) (*dto.StoreAppDTO, error)
	ListStores(ctx context.Context, page, pageSize int32) (*dto.StoreListAppDTO, error)
	SearchStores(ctx context.Context, query string, page, pageSize int32) (*dto.StoreListAppDTO, error)
	UpdateStore(ctx context.Context, id string, req *dto.UpdateStoreAppDTO) (*dto.StoreAppDTO, error)
	DeleteStore(ctx context.Context, id string) error

	// Store management
	UpdateStoreRating(ctx context.Context, id string, rating float64, reviewCount int32) error
	UpdateStoreProductCount(ctx context.Context, id string, increment int32) error
	UpdateStoreFollowerCount(ctx context.Context, id string, increment int32) error

	// Statistics
	GetStoreStats(ctx context.Context) (*dto.StoreStatsAppDTO, error)
}
