package repository

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/review/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/review/domain/repository"
	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type reviewRepository struct {
	db      *pgxpool.Pool
	queries *gen.Queries
}

func NewReviewRepository(db *pgxpool.Pool) repository.ReviewRepository {
	return &reviewRepository{
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

func stringPtrToUUID(s *string) pgtype.UUID {
	if s == nil {
		return pgtype.UUID{Valid: false}
	}
	return stringToUUID(*s)
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

func (r *reviewRepository) CreateReview(ctx context.Context, req repository.CreateReviewParams) (*entity.Review, error) {
	review, err := r.queries.CreateReview(ctx, gen.CreateReviewParams{
		ProductID:  stringToUUID(req.ProductID),
		BuyerID:    req.UserID,
		SellerID:   req.SellerID,
		Rating:     req.Rating,
		Comment:    stringPtrToText(req.Comment),
		Images:     req.Images,
		IsVerified: boolToPgBool(req.IsVerified),
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCReview(review)
	return &result, nil
}

func (r *reviewRepository) GetReview(ctx context.Context, id string) (*entity.ReviewWithDetails, error) {
	review, err := r.queries.GetReview(ctx, stringToUUID(id))
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCReviewWithDetails(review)
	return &result, nil
}

func (r *reviewRepository) ListReviews(ctx context.Context, limit, offset int32) ([]entity.ReviewWithDetails, error) {
	reviews, err := r.queries.ListReviews(ctx, gen.ListReviewsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.ReviewWithDetails, len(reviews))
	for i, r := range reviews {
		result[i] = entity.FromSQLCReviewWithDetailsFromList(r)
	}

	return result, nil
}

func (r *reviewRepository) ListReviewsByProduct(ctx context.Context, productID string, limit, offset int32) ([]entity.ReviewWithDetails, error) {
	reviews, err := r.queries.ListReviewsByProduct(ctx, gen.ListReviewsByProductParams{
		ProductID: stringToUUID(productID),
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.ReviewWithDetails, len(reviews))
	for i, r := range reviews {
		result[i] = entity.FromSQLCReviewWithDetailsFromListByProduct(r)
	}

	return result, nil
}

func (r *reviewRepository) ListReviewsByUser(ctx context.Context, userID int32, limit, offset int32) ([]entity.ReviewWithDetails, error) {
	reviews, err := r.queries.ListReviewsByUser(ctx, gen.ListReviewsByUserParams{
		BuyerID: userID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.ReviewWithDetails, len(reviews))
	for i, r := range reviews {
		result[i] = entity.FromSQLCReviewWithDetailsFromListByUser(r)
	}

	return result, nil
}

func (r *reviewRepository) UpdateReview(ctx context.Context, req repository.UpdateReviewParams) (*entity.Review, error) {
	review, err := r.queries.UpdateReview(ctx, gen.UpdateReviewParams{
		ID:         stringToUUID(req.ID),
		Rating:     req.Rating,
		Comment:    stringPtrToText(req.Comment),
		Images:     req.Images,
		IsVerified: boolToPgBool(req.IsVerified),
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCReview(review)
	return &result, nil
}

func (r *reviewRepository) DeleteReview(ctx context.Context, id string) error {
	return r.queries.DeleteReview(ctx, stringToUUID(id))
}

func (r *reviewRepository) MarkReviewHelpful(ctx context.Context, reviewID string, userID int32) error {
	// Note: review_helpful table does not exist in current schema
	// This is a placeholder implementation
	return r.queries.MarkReviewHelpful(ctx)
}

func (r *reviewRepository) UnmarkReviewHelpful(ctx context.Context, reviewID string, userID int32) error {
	// Note: review_helpful table does not exist in current schema
	// This is a placeholder implementation
	return r.queries.UnmarkReviewHelpful(ctx)
}

func (r *reviewRepository) GetReviewStats(ctx context.Context, productID string) (*entity.ReviewStats, error) {
	stats, err := r.queries.GetReviewStats(ctx)
	if err != nil {
		return nil, err
	}

	result := &entity.ReviewStats{
		TotalReviews:    stats.TotalReviews,
		AverageRating:   stats.AverageRating,
		RatingCounts:    make(map[int32]int64), // RatingCounts not available in current schema
		VerifiedReviews: 0,                     // VerifiedReviews not available in current schema
	}

	return result, nil
}

func (r *reviewRepository) GetUserReviewStats(ctx context.Context, userID int32) (*entity.ReviewStats, error) {
	stats, err := r.queries.GetUserReviewStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := &entity.ReviewStats{
		TotalReviews:    stats.TotalReviews,
		AverageRating:   stats.AverageRating,
		RatingCounts:    make(map[int32]int64), // RatingCounts not available in current schema
		VerifiedReviews: 0,                     // VerifiedReviews not available in current schema
	}

	return result, nil
}
