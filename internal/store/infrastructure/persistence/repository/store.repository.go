package repository

import (
	"context"
	"math/big"

	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/edynt/chogiare/veloras-api/internal/store/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/store/domain/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type storeRepository struct {
	db      *pgxpool.Pool
	queries *gen.Queries
}

func NewStoreRepository(db *pgxpool.Pool) repository.StoreRepository {
	return &storeRepository{
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

func boolToPgBool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

func int32ToInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{Int32: i, Valid: true}
}

func (r *storeRepository) CreateStore(ctx context.Context, req repository.CreateStoreParams) (*entity.Store, error) {
	store, err := r.queries.CreateStore(ctx, gen.CreateStoreParams{
		UserID:      req.UserID,
		Name:        req.Name,
		Description: stringPtrToText(req.Description),
		Logo:        stringPtrToText(req.Logo),
		Banner:      stringPtrToText(req.Banner),
		Address:     stringPtrToText(req.Address),
		Phone:       stringPtrToText(req.Phone),
		Email:       stringPtrToText(req.Email),
		IsVerified:  boolToPgBool(req.IsVerified),
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCStore(store)
	return &result, nil
}

func (r *storeRepository) GetStore(ctx context.Context, id string) (*entity.StoreWithDetails, error) {
	store, err := r.queries.GetStore(ctx, stringToUUID(id))
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCStoreWithDetails(store)
	return &result, nil
}

func (r *storeRepository) GetStoreByUserID(ctx context.Context, userID int32) (*entity.StoreWithDetails, error) {
	store, err := r.queries.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCStoreWithDetailsFromGetByUserID(store)
	return &result, nil
}

func (r *storeRepository) ListStores(ctx context.Context, limit, offset int32) ([]entity.StoreWithDetails, error) {
	stores, err := r.queries.ListStores(ctx, gen.ListStoresParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.StoreWithDetails, len(stores))
	for i, s := range stores {
		result[i] = entity.FromSQLCStoreWithDetailsFromList(s)
	}

	return result, nil
}

func (r *storeRepository) SearchStores(ctx context.Context, query string, limit, offset int32) ([]entity.StoreWithDetails, error) {
	stores, err := r.queries.SearchStores(ctx, gen.SearchStoresParams{
		Column1: stringPtrToText(&query),
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.StoreWithDetails, len(stores))
	for i, s := range stores {
		result[i] = entity.FromSQLCStoreWithDetailsFromSearch(s)
	}

	return result, nil
}

func (r *storeRepository) UpdateStore(ctx context.Context, req repository.UpdateStoreParams) (*entity.Store, error) {
	store, err := r.queries.UpdateStore(ctx, gen.UpdateStoreParams{
		ID:          stringToUUID(req.ID),
		Name:        req.Name,
		Description: stringPtrToText(req.Description),
		Logo:        stringPtrToText(req.Logo),
		Banner:      stringPtrToText(req.Banner),
		Address:     stringPtrToText(req.Address),
		Phone:       stringPtrToText(req.Phone),
		Email:       stringPtrToText(req.Email),
		IsVerified:  boolToPgBool(req.IsVerified),
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCStore(store)
	return &result, nil
}

func (r *storeRepository) UpdateStoreRating(ctx context.Context, id string, rating float64, reviewCount int32) error {
	return r.queries.UpdateStoreRating(ctx, gen.UpdateStoreRatingParams{
		ID:          stringToUUID(id),
		Rating:      float64ToNumeric(rating),
		ReviewCount: int32ToInt4(reviewCount),
	})
}

func (r *storeRepository) UpdateStoreProductCount(ctx context.Context, id string, increment int32) error {
	// Note: product_count column does not exist in current schema
	// This is a placeholder implementation
	return r.queries.UpdateStoreProductCount(ctx)
}

func (r *storeRepository) UpdateStoreFollowerCount(ctx context.Context, id string, increment int32) error {
	// Note: follower_count column does not exist in current schema
	// This is a placeholder implementation
	return r.queries.UpdateStoreFollowerCount(ctx)
}

func (r *storeRepository) DeleteStore(ctx context.Context, id string) error {
	return r.queries.DeleteStore(ctx, stringToUUID(id))
}

func (r *storeRepository) GetStoreStats(ctx context.Context) (*entity.StoreStats, error) {
	stats, err := r.queries.GetStoreStats(ctx)
	if err != nil {
		return nil, err
	}

	result := &entity.StoreStats{
		TotalStores:    stats.TotalStores,
		ActiveStores:   0, // ActiveStores not available in current schema
		VerifiedStores: stats.VerifiedStores,
		AverageRating:  stats.AverageRating,
		TotalProducts:  0, // TotalProducts not available in current schema
		TotalFollowers: 0, // TotalFollowers not available in current schema
	}

	return result, nil
}
