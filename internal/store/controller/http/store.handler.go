package http

import (
	"net/http"
	"strconv"

	"github.com/edynt/chogiare/veloras-api/internal/store/application/service"
	appDto "github.com/edynt/chogiare/veloras-api/internal/store/application/service/dto"
	"github.com/edynt/chogiare/veloras-api/internal/store/controller/dto"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/edynt/chogiare/veloras-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type StoreHandler struct {
	service service.StoreService
}

func NewStoreHandler(service service.StoreService) *StoreHandler {
	return &StoreHandler{service: service}
}

// CreateStore
// @Summary Create a new store
// @Description Create a new store for the current user
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param store body dto.CreateStoreRequest true "Store data"
// @Success 201 {object} dto.StoreResponse "Store created successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /stores [post]
func (h *StoreHandler) CreateStore(ctx *gin.Context) (res interface{}, err error) {
	var req dto.CreateStoreRequest
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
	appReq := &appDto.CreateStoreAppDTO{
		Name:        req.Name,
		Description: req.Description,
		Logo:        req.Logo,
		Banner:      req.Banner,
		Website:     req.Website,
		Phone:       req.Phone,
		Email:       req.Email,
		Address:     req.Address,
		City:        req.City,
		State:       req.State,
		Country:     req.Country,
		PostalCode:  req.PostalCode,
		IsVerified:  req.IsVerified,
		IsActive:    req.IsActive,
	}

	store, err := h.service.CreateStore(ctx, appReq, userID)
	if err != nil {
		return nil, err
	}

	return h.convertToStoreResponse(store), nil
}

// GetStore
// @Summary Get store by ID
// @Description Get store details by ID
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Store ID"
// @Success 200 {object} dto.StoreResponse "Store details"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Store not found"
// @Router /stores/{id} [get]
func (h *StoreHandler) GetStore(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Store ID is required")
	}

	store, err := h.service.GetStore(ctx, id)
	if err != nil {
		return nil, err
	}

	return h.convertToStoreResponse(store), nil
}

// GetMyStore
// @Summary Get current user's store
// @Description Get the store for the current user
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.StoreResponse "User's store"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Store not found"
// @Router /stores/my [get]
func (h *StoreHandler) GetMyStore(ctx *gin.Context) (res interface{}, err error) {
	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	store, err := h.service.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return h.convertToStoreResponse(store), nil
}

// ListStores
// @Summary List all stores
// @Description Get a paginated list of all stores
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} dto.StoreListResponse "List of stores"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /stores [get]
func (h *StoreHandler) ListStores(ctx *gin.Context) (res interface{}, err error) {
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

	stores, err := h.service.ListStores(ctx, int32(page), int32(pageSize))
	if err != nil {
		return nil, err
	}

	return h.convertToStoreListResponse(stores), nil
}

// SearchStores
// @Summary Search stores
// @Description Search stores by name or description
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} dto.StoreListResponse "List of matching stores"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /stores/search [get]
func (h *StoreHandler) SearchStores(ctx *gin.Context) (res interface{}, err error) {
	query := ctx.Query("q")
	if query == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Search query is required")
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

	stores, err := h.service.SearchStores(ctx, query, int32(page), int32(pageSize))
	if err != nil {
		return nil, err
	}

	return h.convertToStoreListResponse(stores), nil
}

// UpdateStore
// @Summary Update store
// @Description Update store details
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Store ID"
// @Param store body dto.UpdateStoreRequest true "Store update data"
// @Success 200 {object} dto.StoreResponse "Store updated successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Store not found"
// @Router /stores/{id} [put]
func (h *StoreHandler) UpdateStore(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Store ID is required")
	}

	var req dto.UpdateStoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", err.Error())
	}

	// Convert to application DTO
	appReq := &appDto.UpdateStoreAppDTO{
		Name:        req.Name,
		Description: req.Description,
		Logo:        req.Logo,
		Banner:      req.Banner,
		Website:     req.Website,
		Phone:       req.Phone,
		Email:       req.Email,
		Address:     req.Address,
		City:        req.City,
		State:       req.State,
		Country:     req.Country,
		PostalCode:  req.PostalCode,
		IsVerified:  req.IsVerified,
		IsActive:    req.IsActive,
	}

	store, err := h.service.UpdateStore(ctx, id, appReq)
	if err != nil {
		return nil, err
	}

	return h.convertToStoreResponse(store), nil
}

// DeleteStore
// @Summary Delete store
// @Description Delete a store
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Store ID"
// @Success 204 "Store deleted successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Store not found"
// @Router /stores/{id} [delete]
func (h *StoreHandler) DeleteStore(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Store ID is required")
	}

	err = h.service.DeleteStore(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetStoreStats
// @Summary Get store statistics
// @Description Get overall store statistics
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.StoreStatsResponse "Store statistics"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /stores/stats [get]
func (h *StoreHandler) GetStoreStats(ctx *gin.Context) (res interface{}, err error) {
	stats, err := h.service.GetStoreStats(ctx)
	if err != nil {
		return nil, err
	}

	return h.convertToStoreStatsResponse(stats), nil
}

// Conversion methods
func (h *StoreHandler) convertToStoreResponse(store *appDto.StoreAppDTO) *dto.StoreResponse {
	return &dto.StoreResponse{
		ID:            store.ID,
		UserID:        store.UserID,
		Name:          store.Name,
		Description:   store.Description,
		Logo:          store.Logo,
		Banner:        store.Banner,
		Website:       store.Website,
		Phone:         store.Phone,
		Email:         store.Email,
		Address:       store.Address,
		City:          store.City,
		State:         store.State,
		Country:       store.Country,
		PostalCode:    store.PostalCode,
		Rating:        store.Rating,
		ReviewCount:   store.ReviewCount,
		ProductCount:  store.ProductCount,
		FollowerCount: store.FollowerCount,
		IsVerified:    store.IsVerified,
		IsActive:      store.IsActive,
		UserName:      store.UserName,
		UserEmail:     store.UserEmail,
		CreatedAt:     store.CreatedAt,
		UpdatedAt:     store.UpdatedAt,
	}
}

func (h *StoreHandler) convertToStoreListResponse(stores *appDto.StoreListAppDTO) *dto.StoreListResponse {
	result := &dto.StoreListResponse{
		Total:      stores.Total,
		Page:       stores.Page,
		PageSize:   stores.PageSize,
		TotalPages: stores.TotalPages,
	}

	result.Stores = make([]dto.StoreResponse, len(stores.Stores))
	for i, store := range stores.Stores {
		result.Stores[i] = *h.convertToStoreResponse(&store)
	}

	return result
}

func (h *StoreHandler) convertToStoreStatsResponse(stats *appDto.StoreStatsAppDTO) *dto.StoreStatsResponse {
	return &dto.StoreStatsResponse{
		TotalStores:    stats.TotalStores,
		ActiveStores:   stats.ActiveStores,
		VerifiedStores: stats.VerifiedStores,
		AverageRating:  stats.AverageRating,
		TotalProducts:  stats.TotalProducts,
		TotalFollowers: stats.TotalFollowers,
	}
}
