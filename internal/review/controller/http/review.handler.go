package http

import (
	"net/http"
	"strconv"

	"github.com/edynt/chogiare/veloras-api/internal/review/application/service"
	appDto "github.com/edynt/chogiare/veloras-api/internal/review/application/service/dto"
	"github.com/edynt/chogiare/veloras-api/internal/review/controller/dto"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/edynt/chogiare/veloras-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type ReviewHandler struct {
	service service.ReviewService
}

func NewReviewHandler(service service.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

// CreateReview
// @Summary Create a new review
// @Description Create a new product review
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param review body dto.CreateReviewRequest true "Review data"
// @Success 201 {object} dto.ReviewResponse "Review created successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /reviews [post]
func (h *ReviewHandler) CreateReview(ctx *gin.Context) (res interface{}, err error) {
	var req dto.CreateReviewRequest
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
	appReq := &appDto.CreateReviewAppDTO{
		ProductID:  req.ProductID,
		SellerID:   req.SellerID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		Images:     req.Images,
		IsVerified: req.IsVerified,
	}

	review, err := h.service.CreateReview(ctx, appReq, userID)
	if err != nil {
		return nil, err
	}

	return h.convertToReviewResponse(review), nil
}

// GetReview
// @Summary Get review by ID
// @Description Get review details by ID
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Review ID"
// @Success 200 {object} dto.ReviewResponse "Review details"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Review not found"
// @Router /reviews/{id} [get]
func (h *ReviewHandler) GetReview(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Review ID is required")
	}

	review, err := h.service.GetReview(ctx, id)
	if err != nil {
		return nil, err
	}

	return h.convertToReviewResponse(review), nil
}

// ListReviews
// @Summary List all reviews
// @Description Get a paginated list of all reviews
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} dto.ReviewListResponse "List of reviews"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /reviews [get]
func (h *ReviewHandler) ListReviews(ctx *gin.Context) (res interface{}, err error) {
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

	reviews, err := h.service.ListReviews(ctx, int32(page), int32(pageSize))
	if err != nil {
		return nil, err
	}

	return h.convertToReviewListResponse(reviews), nil
}

// ListReviewsByProduct
// @Summary List reviews for a product
// @Description Get a paginated list of reviews for a specific product
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product_id path string true "Product ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} dto.ReviewListResponse "List of product reviews"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /reviews/product/{product_id} [get]
func (h *ReviewHandler) ListReviewsByProduct(ctx *gin.Context) (res interface{}, err error) {
	productID := ctx.Param("product_id")
	if productID == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Product ID is required")
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

	reviews, err := h.service.ListReviewsByProduct(ctx, productID, int32(page), int32(pageSize))
	if err != nil {
		return nil, err
	}

	return h.convertToReviewListResponse(reviews), nil
}

// ListUserReviews
// @Summary List user's reviews
// @Description Get a paginated list of reviews by the current user
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} dto.ReviewListResponse "List of user reviews"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /reviews/my [get]
func (h *ReviewHandler) ListUserReviews(ctx *gin.Context) (res interface{}, err error) {
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

	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	reviews, err := h.service.ListReviewsByUser(ctx, userID, int32(page), int32(pageSize))
	if err != nil {
		return nil, err
	}

	return h.convertToReviewListResponse(reviews), nil
}

// UpdateReview
// @Summary Update review
// @Description Update review details
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Review ID"
// @Param review body dto.UpdateReviewRequest true "Review update data"
// @Success 200 {object} dto.ReviewResponse "Review updated successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Review not found"
// @Router /reviews/{id} [put]
func (h *ReviewHandler) UpdateReview(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Review ID is required")
	}

	var req dto.UpdateReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", err.Error())
	}

	// Convert to application DTO
	appReq := &appDto.UpdateReviewAppDTO{
		Rating:     req.Rating,
		Comment:    req.Comment,
		Images:     req.Images,
		IsVerified: req.IsVerified,
	}

	review, err := h.service.UpdateReview(ctx, id, appReq)
	if err != nil {
		return nil, err
	}

	return h.convertToReviewResponse(review), nil
}

// DeleteReview
// @Summary Delete review
// @Description Delete a review
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Review ID"
// @Success 204 "Review deleted successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Review not found"
// @Router /reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Review ID is required")
	}

	err = h.service.DeleteReview(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// MarkReviewHelpful
// @Summary Mark review as helpful
// @Description Mark a review as helpful
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Review ID"
// @Success 200 "Review marked as helpful"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Review not found"
// @Router /reviews/{id}/helpful [post]
func (h *ReviewHandler) MarkReviewHelpful(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Review ID is required")
	}

	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	err = h.service.MarkReviewHelpful(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// UnmarkReviewHelpful
// @Summary Unmark review as helpful
// @Description Unmark a review as helpful
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Review ID"
// @Success 200 "Review unmarked as helpful"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Review not found"
// @Router /reviews/{id}/helpful [delete]
func (h *ReviewHandler) UnmarkReviewHelpful(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Review ID is required")
	}

	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	err = h.service.UnmarkReviewHelpful(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetReviewStats
// @Summary Get review statistics
// @Description Get review statistics for a product
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product_id path string true "Product ID"
// @Success 200 {object} dto.ReviewStatsResponse "Review statistics"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /reviews/stats/product/{product_id} [get]
func (h *ReviewHandler) GetReviewStats(ctx *gin.Context) (res interface{}, err error) {
	productID := ctx.Param("product_id")
	if productID == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Product ID is required")
	}

	stats, err := h.service.GetReviewStats(ctx, productID)
	if err != nil {
		return nil, err
	}

	return h.convertToReviewStatsResponse(stats), nil
}

// GetUserReviewStats
// @Summary Get user review statistics
// @Description Get review statistics for the current user
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ReviewStatsResponse "User review statistics"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /reviews/stats/my [get]
func (h *ReviewHandler) GetUserReviewStats(ctx *gin.Context) (res interface{}, err error) {
	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	stats, err := h.service.GetUserReviewStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return h.convertToReviewStatsResponse(stats), nil
}

// Conversion methods
func (h *ReviewHandler) convertToReviewResponse(review *appDto.ReviewAppDTO) *dto.ReviewResponse {
	return &dto.ReviewResponse{
		ID:           review.ID,
		ProductID:    review.ProductID,
		BuyerID:      review.BuyerID,
		SellerID:     review.SellerID,
		Rating:       review.Rating,
		Comment:      review.Comment,
		Images:       review.Images,
		IsVerified:   review.IsVerified,
		UserName:     review.UserName,
		UserEmail:    review.UserEmail,
		UserAvatar:   review.UserAvatar,
		ProductName:  review.ProductName,
		ProductImage: review.ProductImage,
		CreatedAt:    review.CreatedAt,
		UpdatedAt:    review.UpdatedAt,
	}
}

func (h *ReviewHandler) convertToReviewListResponse(reviews *appDto.ReviewListAppDTO) *dto.ReviewListResponse {
	result := &dto.ReviewListResponse{
		Total:      reviews.Total,
		Page:       reviews.Page,
		PageSize:   reviews.PageSize,
		TotalPages: reviews.TotalPages,
	}

	result.Reviews = make([]dto.ReviewResponse, len(reviews.Reviews))
	for i, review := range reviews.Reviews {
		result.Reviews[i] = *h.convertToReviewResponse(&review)
	}

	return result
}

func (h *ReviewHandler) convertToReviewStatsResponse(stats *appDto.ReviewStatsAppDTO) *dto.ReviewStatsResponse {
	return &dto.ReviewStatsResponse{
		TotalReviews:    stats.TotalReviews,
		AverageRating:   stats.AverageRating,
		RatingCounts:    stats.RatingCounts,
		VerifiedReviews: stats.VerifiedReviews,
	}
}
