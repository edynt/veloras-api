package entity

import (
	"time"

	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/jackc/pgx/v5/pgtype"
)

type Product struct {
	ID            string
	Title         string
	Description   string
	Price         float64
	OriginalPrice *float64
	CategoryID    string
	Images        []string
	Condition     string
	Tags          []string
	Location      string
	Stock         int32
	SellerID      int32
	Status        string
	Badges        []string
	Rating        float64
	ReviewCount   int32
	ViewCount     int32
	IsFeatured    bool
	IsPromoted    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ProductWithDetails struct {
	Product
	CategoryName *string
	CategorySlug *string
	SellerName   *string
	SellerEmail  *string
	StoreName    *string
	StoreLogo    *string
	StoreRating  *float64
}

type ProductList struct {
	Products   []ProductWithDetails
	Total      int64
	Page       int32
	PageSize   int32
	TotalPages int32
}

type ProductStats struct {
	TotalProducts    int64
	ActiveProducts   int64
	DraftProducts    int64
	SoldProducts     int64
	FeaturedProducts int64
	PromotedProducts int64
	AverageRating    *float64
	TotalViews       *int64
}

// Convert from SQLC model to domain entity
func FromSQLCProduct(p gen.Product) Product {
	return Product{
		ID:            p.ID.String(),
		Title:         p.Title,
		Description:   p.Description,
		Price:         convertNumericToFloat64(p.Price),
		OriginalPrice: convertNumericPtrToFloat64Ptr(p.OriginalPrice),
		CategoryID:    p.CategoryID.String(),
		Images:        p.Images,
		Condition:     p.Condition,
		Tags:          p.Tags,
		Location:      p.Location,
		Stock:         p.Stock,
		SellerID:      p.SellerID,
		Status:        p.Status,
		Badges:        p.Badges,
		Rating:        convertNumericToFloat64(p.Rating),
		ReviewCount:   p.ReviewCount.Int32,
		ViewCount:     p.ViewCount.Int32,
		IsFeatured:    p.IsFeatured.Bool,
		IsPromoted:    p.IsPromoted.Bool,
		CreatedAt:     p.CreatedAt.Time,
		UpdatedAt:     p.UpdatedAt.Time,
	}
}

// Convert from SQLC model with details to domain entity
func FromSQLCProductWithDetails(p gen.GetProductRow) ProductWithDetails {
	return ProductWithDetails{
		Product: Product{
			ID:            p.ID.String(),
			Title:         p.Title,
			Description:   p.Description,
			Price:         convertNumericToFloat64(p.Price),
			OriginalPrice: convertNumericPtrToFloat64Ptr(p.OriginalPrice),
			CategoryID:    p.CategoryID.String(),
			Images:        p.Images,
			Condition:     p.Condition,
			Tags:          p.Tags,
			Location:      p.Location,
			Stock:         p.Stock,
			SellerID:      p.SellerID,
			Status:        p.Status,
			Badges:        p.Badges,
			Rating:        convertNumericToFloat64(p.Rating),
			ReviewCount:   p.ReviewCount.Int32,
			ViewCount:     p.ViewCount.Int32,
			IsFeatured:    p.IsFeatured.Bool,
			IsPromoted:    p.IsPromoted.Bool,
			CreatedAt:     p.CreatedAt.Time,
			UpdatedAt:     p.UpdatedAt.Time,
		},
		CategoryName: convertTextPtr(p.CategoryName),
		CategorySlug: convertTextPtr(p.CategorySlug),
		SellerName:   convertInterfacePtr(p.SellerName),
		SellerEmail:  convertTextPtr(p.SellerEmail),
		StoreName:    convertTextPtr(p.StoreName),
		StoreLogo:    convertTextPtr(p.StoreLogo),
		StoreRating:  convertNumericPtrToFloat64Ptr(p.StoreRating),
	}
}

// Helper function to convert pgtype.Numeric to float64
func convertNumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	return float64(n.Int.Int64())
}

// Helper function to convert *pgtype.Numeric to *float64
func convertNumericPtrToFloat64Ptr(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}
	val := float64(n.Int.Int64())
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
