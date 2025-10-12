package repository

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/product/domain/model/entity"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, req CreateProductParams) (*entity.Product, error)
	GetProduct(ctx context.Context, id string) (*entity.ProductWithDetails, error)
	ListProducts(ctx context.Context, limit, offset int32) ([]entity.ProductWithDetails, error)
	ListProductsByCategory(ctx context.Context, categoryID string, limit, offset int32) ([]entity.ProductWithDetails, error)
	ListProductsBySeller(ctx context.Context, sellerID int32, limit, offset int32) ([]entity.ProductWithDetails, error)
	SearchProducts(ctx context.Context, query string, limit, offset int32) ([]entity.ProductWithDetails, error)
	GetFeaturedProducts(ctx context.Context, limit int32) ([]entity.ProductWithDetails, error)
	GetPromotedProducts(ctx context.Context, limit int32) ([]entity.ProductWithDetails, error)
	UpdateProduct(ctx context.Context, req UpdateProductParams) (*entity.Product, error)
	UpdateProductStatus(ctx context.Context, id string, status string) (*entity.Product, error)
	UpdateProductStock(ctx context.Context, id string, stock int32) (*entity.Product, error)
	IncrementProductViews(ctx context.Context, id string) error
	UpdateProductRating(ctx context.Context, id string, rating float64, reviewCount int32) error
	DeleteProduct(ctx context.Context, id string) error
	GetProductStats(ctx context.Context) (*entity.ProductStats, error)
	GetProductsByFilters(ctx context.Context, req GetProductsByFiltersParams) ([]entity.ProductWithDetails, error)
}

type CreateProductParams struct {
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
	IsFeatured    bool
	IsPromoted    bool
}

type UpdateProductParams struct {
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
	Status        string
	Badges        []string
	IsFeatured    bool
	IsPromoted    bool
}

type GetProductsByFiltersParams struct {
	CategoryID *string
	MinPrice   *float64
	MaxPrice   *float64
	Condition  *string
	Location   *string
	SellerID   *int32
	Rating     *float64
	SortBy     *string
	Limit      int32
	Offset     int32
}
