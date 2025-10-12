package entity

import (
	"time"

	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/jackc/pgx/v5/pgtype"
)

type Review struct {
	ID         string
	ProductID  string
	BuyerID    int32
	SellerID   int32
	Rating     int32
	Comment    *string
	Images     []string
	IsVerified bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ReviewWithDetails struct {
	Review
	UserName     *string
	UserEmail    *string
	UserAvatar   *string
	ProductName  *string
	ProductImage *string
}

type ReviewStats struct {
	TotalReviews    int64
	AverageRating   float64
	RatingCounts    map[int32]int64
	VerifiedReviews int64
}

// Convert from SQLC model to domain entity
func FromSQLCReview(r gen.Review) Review {
	return Review{
		ID:         r.ID.String(),
		ProductID:  r.ProductID.String(),
		BuyerID:    r.BuyerID,
		SellerID:   r.SellerID,
		Rating:     r.Rating,
		Comment:    convertTextPtr(r.Comment),
		Images:     r.Images,
		IsVerified: r.IsVerified.Bool,
		CreatedAt:  r.CreatedAt.Time,
		UpdatedAt:  r.CreatedAt.Time, // UpdatedAt not in generated model, using CreatedAt
	}
}

// Convert from SQLC model with details to domain entity
func FromSQLCReviewWithDetails(r gen.GetReviewRow) ReviewWithDetails {
	return ReviewWithDetails{
		Review: Review{
			ID:         r.ID.String(),
			ProductID:  r.ProductID.String(),
			BuyerID:    r.BuyerID,
			SellerID:   r.SellerID,
			Rating:     r.Rating,
			Comment:    convertTextPtr(r.Comment),
			Images:     r.Images,
			IsVerified: r.IsVerified.Bool,
			CreatedAt:  r.CreatedAt.Time,
			UpdatedAt:  r.CreatedAt.Time, // UpdatedAt not in generated model, using CreatedAt
		},
		UserName:     nil, // BuyerName not available in current schema // Using BuyerName from generated model
		UserEmail:    nil, // UserEmail not in generated model
		UserAvatar:   nil, // UserAvatar not in generated model
		ProductName:  nil, // ProductTitle not available in current schema   // Using ProductTitle from generated model
		ProductImage: nil, // ProductImage not in generated model
	}
}

// Convert from SQLC ListReviewsRow to domain entity
func FromSQLCReviewWithDetailsFromList(r gen.ListReviewsRow) ReviewWithDetails {
	return ReviewWithDetails{
		Review: Review{
			ID:         r.ID.String(),
			ProductID:  r.ProductID.String(),
			BuyerID:    r.BuyerID,
			SellerID:   r.SellerID,
			Rating:     r.Rating,
			Comment:    convertTextPtr(r.Comment),
			Images:     r.Images,
			IsVerified: r.IsVerified.Bool,
			CreatedAt:  r.CreatedAt.Time,
			UpdatedAt:  r.CreatedAt.Time,
		},
		UserName:     nil, // BuyerName not available in current schema
		UserEmail:    nil,
		UserAvatar:   nil,
		ProductName:  nil, // ProductTitle not available in current schema
		ProductImage: nil, // ProductImage not available in current schema
	}
}

// Convert from SQLC ListReviewsByProductRow to domain entity
func FromSQLCReviewWithDetailsFromListByProduct(r gen.ListReviewsByProductRow) ReviewWithDetails {
	return ReviewWithDetails{
		Review: Review{
			ID:         r.ID.String(),
			ProductID:  r.ProductID.String(),
			BuyerID:    r.BuyerID,
			SellerID:   r.SellerID,
			Rating:     r.Rating,
			Comment:    convertTextPtr(r.Comment),
			Images:     r.Images,
			IsVerified: r.IsVerified.Bool,
			CreatedAt:  r.CreatedAt.Time,
			UpdatedAt:  r.CreatedAt.Time,
		},
		UserName:     nil, // BuyerName not available in current schema
		UserEmail:    nil,
		UserAvatar:   nil,
		ProductName:  nil, // ProductTitle not available in current schema
		ProductImage: nil, // ProductImage not available in current schema
	}
}

// Convert from SQLC ListReviewsByUserRow to domain entity
func FromSQLCReviewWithDetailsFromListByUser(r gen.ListReviewsByUserRow) ReviewWithDetails {
	return ReviewWithDetails{
		Review: Review{
			ID:         r.ID.String(),
			ProductID:  r.ProductID.String(),
			BuyerID:    r.BuyerID,
			SellerID:   r.SellerID,
			Rating:     r.Rating,
			Comment:    convertTextPtr(r.Comment),
			Images:     r.Images,
			IsVerified: r.IsVerified.Bool,
			CreatedAt:  r.CreatedAt.Time,
			UpdatedAt:  r.CreatedAt.Time,
		},
		UserName:     nil, // BuyerName not available in current schema
		UserEmail:    nil,
		UserAvatar:   nil,
		ProductName:  nil, // ProductTitle not available in current schema
		ProductImage: nil, // ProductImage not available in current schema
	}
}

// Helper function to convert pgtype.UUID to *string
func convertUUIDPtr(u pgtype.UUID) *string {
	if !u.Valid {
		return nil
	}
	val := u.String()
	return &val
}

// Helper function to convert pgtype.Text to *string
func convertTextPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

// Helper function to convert interface{} to *string
func convertInterfacePtr(i interface{}) *string {
	if i == nil {
		return nil
	}
	if str, ok := i.(string); ok {
		return &str
	}
	return nil
}
