package repository

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/review/domain/model/entity"
)

type ReviewRepository interface {
	// Review CRUD operations
	CreateReview(ctx context.Context, req CreateReviewParams) (*entity.Review, error)
	GetReview(ctx context.Context, id string) (*entity.ReviewWithDetails, error)
	ListReviews(ctx context.Context, limit, offset int32) ([]entity.ReviewWithDetails, error)
	ListReviewsByProduct(ctx context.Context, productID string, limit, offset int32) ([]entity.ReviewWithDetails, error)
	ListReviewsByUser(ctx context.Context, userID int32, limit, offset int32) ([]entity.ReviewWithDetails, error)
	UpdateReview(ctx context.Context, req UpdateReviewParams) (*entity.Review, error)
	DeleteReview(ctx context.Context, id string) error

	// Review interactions
	MarkReviewHelpful(ctx context.Context, reviewID string, userID int32) error
	UnmarkReviewHelpful(ctx context.Context, reviewID string, userID int32) error

	// Statistics
	GetReviewStats(ctx context.Context, productID string) (*entity.ReviewStats, error)
	GetUserReviewStats(ctx context.Context, userID int32) (*entity.ReviewStats, error)
}

type CreateReviewParams struct {
	ProductID  string
	UserID     int32
	OrderID    *string
	Rating     int32
	Title      *string
	Comment    *string
	Images     []string
	IsVerified bool
}

type UpdateReviewParams struct {
	ID         string
	Rating     int32
	Title      *string
	Comment    *string
	Images     []string
	IsVerified bool
}
