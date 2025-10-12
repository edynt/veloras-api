package service

import (
	"context"
	"math"
	"net/http"

	"github.com/edynt/chogiare/veloras-api/internal/review/application/service/dto"
	"github.com/edynt/chogiare/veloras-api/internal/review/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/review/domain/repository"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
)

type reviewService struct {
	reviewRepo repository.ReviewRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, req *dto.CreateReviewAppDTO, userID int32) (*dto.ReviewAppDTO, error) {
	review, err := s.reviewRepo.CreateReview(ctx, repository.CreateReviewParams{
		ProductID:  req.ProductID,
		UserID:     userID,
		OrderID:    req.OrderID,
		Rating:     req.Rating,
		Title:      req.Title,
		Comment:    req.Comment,
		Images:     req.Images,
		IsVerified: req.IsVerified,
	})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to create review", err)
	}

	return s.convertToReviewAppDTO(review), nil
}

func (s *reviewService) GetReview(ctx context.Context, id string) (*dto.ReviewAppDTO, error) {
	review, err := s.reviewRepo.GetReview(ctx, id)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Review not found", err)
	}

	return s.convertToReviewAppDTOWithDetails(review), nil
}

func (s *reviewService) ListReviews(ctx context.Context, page, pageSize int32) (*dto.ReviewListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	reviews, err := s.reviewRepo.ListReviews(ctx, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list reviews", err)
	}

	total := int64(len(reviews))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.ReviewListAppDTO{
		Reviews:    s.convertToReviewAppDTOsWithDetails(reviews),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *reviewService) ListReviewsByProduct(ctx context.Context, productID string, page, pageSize int32) (*dto.ReviewListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	reviews, err := s.reviewRepo.ListReviewsByProduct(ctx, productID, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list product reviews", err)
	}

	total := int64(len(reviews))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.ReviewListAppDTO{
		Reviews:    s.convertToReviewAppDTOsWithDetails(reviews),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *reviewService) ListReviewsByUser(ctx context.Context, userID int32, page, pageSize int32) (*dto.ReviewListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	reviews, err := s.reviewRepo.ListReviewsByUser(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list user reviews", err)
	}

	total := int64(len(reviews))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.ReviewListAppDTO{
		Reviews:    s.convertToReviewAppDTOsWithDetails(reviews),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *reviewService) UpdateReview(ctx context.Context, id string, req *dto.UpdateReviewAppDTO) (*dto.ReviewAppDTO, error) {
	review, err := s.reviewRepo.UpdateReview(ctx, repository.UpdateReviewParams{
		ID:         id,
		Rating:     req.Rating,
		Title:      req.Title,
		Comment:    req.Comment,
		Images:     req.Images,
		IsVerified: req.IsVerified,
	})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to update review", err)
	}

	return s.convertToReviewAppDTO(review), nil
}

func (s *reviewService) DeleteReview(ctx context.Context, id string) error {
	err := s.reviewRepo.DeleteReview(ctx, id)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to delete review", err)
	}
	return nil
}

func (s *reviewService) MarkReviewHelpful(ctx context.Context, reviewID string, userID int32) error {
	err := s.reviewRepo.MarkReviewHelpful(ctx, reviewID, userID)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to mark review as helpful", err)
	}
	return nil
}

func (s *reviewService) UnmarkReviewHelpful(ctx context.Context, reviewID string, userID int32) error {
	err := s.reviewRepo.UnmarkReviewHelpful(ctx, reviewID, userID)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to unmark review as helpful", err)
	}
	return nil
}

func (s *reviewService) GetReviewStats(ctx context.Context, productID string) (*dto.ReviewStatsAppDTO, error) {
	stats, err := s.reviewRepo.GetReviewStats(ctx, productID)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to get review stats", err)
	}

	return &dto.ReviewStatsAppDTO{
		TotalReviews:    stats.TotalReviews,
		AverageRating:   stats.AverageRating,
		RatingCounts:    stats.RatingCounts,
		VerifiedReviews: stats.VerifiedReviews,
	}, nil
}

func (s *reviewService) GetUserReviewStats(ctx context.Context, userID int32) (*dto.ReviewStatsAppDTO, error) {
	stats, err := s.reviewRepo.GetUserReviewStats(ctx, userID)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to get user review stats", err)
	}

	return &dto.ReviewStatsAppDTO{
		TotalReviews:    stats.TotalReviews,
		AverageRating:   stats.AverageRating,
		RatingCounts:    stats.RatingCounts,
		VerifiedReviews: stats.VerifiedReviews,
	}, nil
}

// Conversion methods
func (s *reviewService) convertToReviewAppDTO(review *entity.Review) *dto.ReviewAppDTO {
	return &dto.ReviewAppDTO{
		ID:         review.ID,
		ProductID:  review.ProductID,
		UserID:     review.UserID,
		OrderID:    review.OrderID,
		Rating:     review.Rating,
		Title:      review.Title,
		Comment:    review.Comment,
		Images:     review.Images,
		IsVerified: review.IsVerified,
		Helpful:    review.Helpful,
		CreatedAt:  review.CreatedAt,
		UpdatedAt:  review.UpdatedAt,
	}
}

func (s *reviewService) convertToReviewAppDTOWithDetails(review *entity.ReviewWithDetails) *dto.ReviewAppDTO {
	return &dto.ReviewAppDTO{
		ID:           review.ID,
		ProductID:    review.ProductID,
		UserID:       review.UserID,
		OrderID:      review.OrderID,
		Rating:       review.Rating,
		Title:        review.Title,
		Comment:      review.Comment,
		Images:       review.Images,
		IsVerified:   review.IsVerified,
		Helpful:      review.Helpful,
		UserName:     review.UserName,
		UserEmail:    review.UserEmail,
		UserAvatar:   review.UserAvatar,
		ProductName:  review.ProductName,
		ProductImage: review.ProductImage,
		CreatedAt:    review.CreatedAt,
		UpdatedAt:    review.UpdatedAt,
	}
}

func (s *reviewService) convertToReviewAppDTOsWithDetails(reviews []entity.ReviewWithDetails) []dto.ReviewAppDTO {
	result := make([]dto.ReviewAppDTO, len(reviews))
	for i, review := range reviews {
		result[i] = *s.convertToReviewAppDTOWithDetails(&review)
	}
	return result
}
