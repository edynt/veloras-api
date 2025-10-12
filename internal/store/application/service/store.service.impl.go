package service

import (
	"context"
	"math"
	"net/http"

	"github.com/edynt/chogiare/veloras-api/internal/store/application/service/dto"
	"github.com/edynt/chogiare/veloras-api/internal/store/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/store/domain/repository"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
)

type storeService struct {
	storeRepo repository.StoreRepository
}

func NewStoreService(storeRepo repository.StoreRepository) StoreService {
	return &storeService{
		storeRepo: storeRepo,
	}
}

func (s *storeService) CreateStore(ctx context.Context, req *dto.CreateStoreAppDTO, userID int32) (*dto.StoreAppDTO, error) {
	store, err := s.storeRepo.CreateStore(ctx, repository.CreateStoreParams{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Logo:        req.Logo,
		Banner:      req.Banner,
		Phone:       req.Phone,
		Email:       req.Email,
		Address:     req.Address,
		IsVerified:  req.IsVerified,
	})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to create store", err)
	}

	return s.convertToStoreAppDTO(store), nil
}

func (s *storeService) GetStore(ctx context.Context, id string) (*dto.StoreAppDTO, error) {
	store, err := s.storeRepo.GetStore(ctx, id)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Store not found", err)
	}

	return s.convertToStoreAppDTOWithDetails(store), nil
}

func (s *storeService) GetStoreByUserID(ctx context.Context, userID int32) (*dto.StoreAppDTO, error) {
	store, err := s.storeRepo.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Store not found", err)
	}

	return s.convertToStoreAppDTOWithDetails(store), nil
}

func (s *storeService) ListStores(ctx context.Context, page, pageSize int32) (*dto.StoreListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	stores, err := s.storeRepo.ListStores(ctx, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list stores", err)
	}

	total := int64(len(stores))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.StoreListAppDTO{
		Stores:     s.convertToStoreAppDTOsWithDetails(stores),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *storeService) SearchStores(ctx context.Context, query string, page, pageSize int32) (*dto.StoreListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	stores, err := s.storeRepo.SearchStores(ctx, query, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to search stores", err)
	}

	total := int64(len(stores))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.StoreListAppDTO{
		Stores:     s.convertToStoreAppDTOsWithDetails(stores),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *storeService) UpdateStore(ctx context.Context, id string, req *dto.UpdateStoreAppDTO) (*dto.StoreAppDTO, error) {
	// Get existing store to preserve unchanged fields
	existingStore, err := s.storeRepo.GetStore(ctx, id)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Store not found", err)
	}

	updateParams := repository.UpdateStoreParams{
		ID:          id,
		Name:        existingStore.Name,
		Description: existingStore.Description,
		Logo:        existingStore.Logo,
		Banner:      existingStore.Banner,
		Phone:       existingStore.Phone,
		Email:       existingStore.Email,
		Address:     existingStore.Address,
		IsVerified:  existingStore.IsVerified,
	}

	// Update only provided fields
	if req.Name != nil {
		updateParams.Name = *req.Name
	}
	if req.Description != nil {
		updateParams.Description = req.Description
	}
	if req.Logo != nil {
		updateParams.Logo = req.Logo
	}
	if req.Banner != nil {
		updateParams.Banner = req.Banner
	}
	if req.Phone != nil {
		updateParams.Phone = req.Phone
	}
	if req.Email != nil {
		updateParams.Email = req.Email
	}
	if req.Address != nil {
		updateParams.Address = req.Address
	}
	if req.IsVerified != nil {
		updateParams.IsVerified = *req.IsVerified
	}

	store, err := s.storeRepo.UpdateStore(ctx, updateParams)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to update store", err)
	}

	return s.convertToStoreAppDTO(store), nil
}

func (s *storeService) DeleteStore(ctx context.Context, id string) error {
	err := s.storeRepo.DeleteStore(ctx, id)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to delete store", err)
	}
	return nil
}

func (s *storeService) UpdateStoreRating(ctx context.Context, id string, rating float64, reviewCount int32) error {
	err := s.storeRepo.UpdateStoreRating(ctx, id, rating, reviewCount)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to update store rating", err)
	}
	return nil
}

func (s *storeService) UpdateStoreProductCount(ctx context.Context, id string, increment int32) error {
	err := s.storeRepo.UpdateStoreProductCount(ctx, id, increment)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to update store product count", err)
	}
	return nil
}

func (s *storeService) UpdateStoreFollowerCount(ctx context.Context, id string, increment int32) error {
	err := s.storeRepo.UpdateStoreFollowerCount(ctx, id, increment)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to update store follower count", err)
	}
	return nil
}

func (s *storeService) GetStoreStats(ctx context.Context) (*dto.StoreStatsAppDTO, error) {
	stats, err := s.storeRepo.GetStoreStats(ctx)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to get store stats", err)
	}

	return &dto.StoreStatsAppDTO{
		TotalStores:    stats.TotalStores,
		ActiveStores:   stats.ActiveStores,
		VerifiedStores: stats.VerifiedStores,
		AverageRating:  stats.AverageRating,
		TotalProducts:  stats.TotalProducts,
		TotalFollowers: stats.TotalFollowers,
	}, nil
}

// Conversion methods
func (s *storeService) convertToStoreAppDTO(store *entity.Store) *dto.StoreAppDTO {
	return &dto.StoreAppDTO{
		ID:          store.ID,
		UserID:      store.UserID,
		Name:        store.Name,
		Description: store.Description,
		Logo:        store.Logo,
		Banner:      store.Banner,
		Phone:       store.Phone,
		Email:       store.Email,
		Address:     store.Address,
		Rating:      store.Rating,
		ReviewCount: store.ReviewCount,
		IsVerified:  store.IsVerified,
		CreatedAt:   store.CreatedAt,
		UpdatedAt:   store.UpdatedAt,
	}
}

func (s *storeService) convertToStoreAppDTOWithDetails(store *entity.StoreWithDetails) *dto.StoreAppDTO {
	return &dto.StoreAppDTO{
		ID:          store.ID,
		UserID:      store.UserID,
		Name:        store.Name,
		Description: store.Description,
		Logo:        store.Logo,
		Banner:      store.Banner,
		Phone:       store.Phone,
		Email:       store.Email,
		Address:     store.Address,
		Rating:      store.Rating,
		ReviewCount: store.ReviewCount,
		IsVerified:  store.IsVerified,
		UserName:    store.UserName,
		UserEmail:   store.UserEmail,
		CreatedAt:   store.CreatedAt,
		UpdatedAt:   store.UpdatedAt,
	}
}

func (s *storeService) convertToStoreAppDTOsWithDetails(stores []entity.StoreWithDetails) []dto.StoreAppDTO {
	result := make([]dto.StoreAppDTO, len(stores))
	for i, store := range stores {
		result[i] = *s.convertToStoreAppDTOWithDetails(&store)
	}
	return result
}
