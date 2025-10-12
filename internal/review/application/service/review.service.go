package service

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/review/application/service/dto"
)

type ReviewService interface {
	// Review operations
	CreateReview(ctx context.Context, req *dto.CreateReviewAppDTO, userID int32) (*dto.ReviewAppDTO, error)
	GetReview(ctx context.Context, id string) (*dto.ReviewAppDTO, error)
	ListReviews(ctx context.Context, page, pageSize int32) (*dto.ReviewListAppDTO, error)
	ListReviewsByProduct(ctx context.Context, productID string, page, pageSize int32) (*dto.ReviewListAppDTO, error)
	ListReviewsByUser(ctx context.Context, userID int32, page, pageSize int32) (*dto.ReviewListAppDTO, error)
	UpdateReview(ctx context.Context, id string, req *dto.UpdateReviewAppDTO) (*dto.ReviewAppDTO, error)
	DeleteReview(ctx context.Context, id string) error

	// Review interactions
	MarkReviewHelpful(ctx context.Context, reviewID string, userID int32) error
	UnmarkReviewHelpful(ctx context.Context, reviewID string, userID int32) error

	// Statistics
	GetReviewStats(ctx context.Context, productID string) (*dto.ReviewStatsAppDTO, error)
	GetUserReviewStats(ctx context.Context, userID int32) (*dto.ReviewStatsAppDTO, error)
}
