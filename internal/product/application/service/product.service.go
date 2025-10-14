package service

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/product/application/service/dto"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req dto.CreateProductAppDTO, sellerID int32) (*dto.ProductAppDTO, error)
	GetProduct(ctx context.Context, id string) (*dto.ProductAppDTO, error)
	ListProducts(ctx context.Context, page, pageSize int32) (*dto.ProductListAppDTO, error)
	ListProductsByCategory(ctx context.Context, categoryID string, page, pageSize int32) (*dto.ProductListAppDTO, error)
	ListProductsBySeller(ctx context.Context, sellerID int32, page, pageSize int32) (*dto.ProductListAppDTO, error)
	SearchProducts(ctx context.Context, req dto.SearchProductsAppDTO) (*dto.ProductListAppDTO, error)
	GetFeaturedProducts(ctx context.Context, limit int32) ([]dto.ProductAppDTO, error)
	GetPromotedProducts(ctx context.Context, limit int32) ([]dto.ProductAppDTO, error)
	UpdateProduct(ctx context.Context, id string, req dto.UpdateProductAppDTO, sellerID int32) (*dto.ProductAppDTO, error)
	UpdateProductStatus(ctx context.Context, id string, status string, sellerID int32) (*dto.ProductAppDTO, error)
	UpdateProductStock(ctx context.Context, id string, stock int32, sellerID int32) (*dto.ProductAppDTO, error)
	IncrementProductViews(ctx context.Context, id string) error
	DeleteProduct(ctx context.Context, id string, sellerID int32) error
	GetProductStats(ctx context.Context) (*dto.ProductStatsAppDTO, error)
	GetMyProducts(ctx context.Context, sellerID int32, page, pageSize int32) (*dto.ProductListAppDTO, error)
	BulkUpdateProducts(ctx context.Context, ids []string, req dto.UpdateProductAppDTO, sellerID int32) ([]dto.ProductAppDTO, error)
	SeedProducts(ctx context.Context, count int32) (*dto.SeedResultAppDTO, error)
}
