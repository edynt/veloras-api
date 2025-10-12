package service

import (
	"context"
	"errors"
	"math"
	"net/http"

	"github.com/edynt/chogiare/veloras-api/internal/product/application/service/dto"
	"github.com/edynt/chogiare/veloras-api/internal/product/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/product/domain/repository"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
)

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req dto.CreateProductAppDTO, sellerID int32) (*dto.ProductAppDTO, error) {
	// Set default status if not provided
	status := "draft"
	if req.Status != "" {
		status = req.Status
	}

	product, err := s.productRepo.CreateProduct(ctx, repository.CreateProductParams{
		Title:         req.Title,
		Description:   req.Description,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		CategoryID:    req.CategoryID,
		Images:        req.Images,
		Condition:     req.Condition,
		Tags:          req.Tags,
		Location:      req.Location,
		Stock:         req.Stock,
		SellerID:      sellerID,
		Status:        status,
		Badges:        req.Badges,
		IsFeatured:    req.IsFeatured,
		IsPromoted:    req.IsPromoted,
	})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to create product", err)
	}

	return s.convertToProductAppDTO(product), nil
}

func (s *productService) GetProduct(ctx context.Context, id string) (*dto.ProductAppDTO, error) {
	product, err := s.productRepo.GetProduct(ctx, id)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Product not found", err)
	}

	return s.convertToProductAppDTOWithDetails(product), nil
}

func (s *productService) ListProducts(ctx context.Context, page, pageSize int32) (*dto.ProductListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	products, err := s.productRepo.ListProducts(ctx, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list products", err)
	}

	// For now, we'll return the products without total count
	// In a real implementation, you'd want to add a count query
	total := int64(len(products))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.ProductListAppDTO{
		Products:   s.convertToProductAppDTOsWithDetails(products),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *productService) ListProductsByCategory(ctx context.Context, categoryID string, page, pageSize int32) (*dto.ProductListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	products, err := s.productRepo.ListProductsByCategory(ctx, categoryID, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list products by category", err)
	}

	total := int64(len(products))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.ProductListAppDTO{
		Products:   s.convertToProductAppDTOsWithDetails(products),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *productService) ListProductsBySeller(ctx context.Context, sellerID int32, page, pageSize int32) (*dto.ProductListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	products, err := s.productRepo.ListProductsBySeller(ctx, sellerID, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list products by seller", err)
	}

	total := int64(len(products))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.ProductListAppDTO{
		Products:   s.convertToProductAppDTOsWithDetails(products),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *productService) SearchProducts(ctx context.Context, req dto.SearchProductsAppDTO) (*dto.ProductListAppDTO, error) {
	page := int32(1)
	if req.Page != nil {
		page = *req.Page
	}
	pageSize := int32(20)
	if req.Limit != nil {
		pageSize = *req.Limit
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	query := ""
	if req.Query != nil {
		query = *req.Query
	}

	products, err := s.productRepo.SearchProducts(ctx, query, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to search products", err)
	}

	total := int64(len(products))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.ProductListAppDTO{
		Products:   s.convertToProductAppDTOsWithDetails(products),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *productService) GetFeaturedProducts(ctx context.Context, limit int32) ([]dto.ProductAppDTO, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	products, err := s.productRepo.GetFeaturedProducts(ctx, limit)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to get featured products", err)
	}

	return s.convertToProductAppDTOsWithDetails(products), nil
}

func (s *productService) GetPromotedProducts(ctx context.Context, limit int32) ([]dto.ProductAppDTO, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	products, err := s.productRepo.GetPromotedProducts(ctx, limit)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to get promoted products", err)
	}

	return s.convertToProductAppDTOsWithDetails(products), nil
}

func (s *productService) UpdateProduct(ctx context.Context, id string, req dto.UpdateProductAppDTO, sellerID int32) (*dto.ProductAppDTO, error) {
	// First, get the existing product to check ownership
	existingProduct, err := s.productRepo.GetProduct(ctx, id)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Product not found", err)
	}

	// Check if the seller owns this product
	if existingProduct.SellerID != sellerID {
		return nil, response.NewAPIError(http.StatusForbidden, "You don't have permission to update this product", errors.New("unauthorized"))
	}

	// Build update parameters
	updateParams := repository.UpdateProductParams{
		ID: id,
	}

	// Only update fields that are provided
	if req.Title != nil {
		updateParams.Title = *req.Title
	} else {
		updateParams.Title = existingProduct.Title
	}

	if req.Description != nil {
		updateParams.Description = *req.Description
	} else {
		updateParams.Description = existingProduct.Description
	}

	if req.Price != nil {
		updateParams.Price = *req.Price
	} else {
		updateParams.Price = existingProduct.Price
	}

	if req.OriginalPrice != nil {
		updateParams.OriginalPrice = req.OriginalPrice
	} else {
		updateParams.OriginalPrice = existingProduct.OriginalPrice
	}

	if req.CategoryID != nil {
		updateParams.CategoryID = *req.CategoryID
	} else {
		updateParams.CategoryID = existingProduct.CategoryID
	}

	if req.Images != nil {
		updateParams.Images = req.Images
	} else {
		updateParams.Images = existingProduct.Images
	}

	if req.Condition != nil {
		updateParams.Condition = *req.Condition
	} else {
		updateParams.Condition = existingProduct.Condition
	}

	if req.Tags != nil {
		updateParams.Tags = req.Tags
	} else {
		updateParams.Tags = existingProduct.Tags
	}

	if req.Location != nil {
		updateParams.Location = *req.Location
	} else {
		updateParams.Location = existingProduct.Location
	}

	if req.Stock != nil {
		updateParams.Stock = *req.Stock
	} else {
		updateParams.Stock = existingProduct.Stock
	}

	if req.Status != nil {
		updateParams.Status = *req.Status
	} else {
		updateParams.Status = existingProduct.Status
	}

	if req.Badges != nil {
		updateParams.Badges = req.Badges
	} else {
		updateParams.Badges = existingProduct.Badges
	}

	if req.IsFeatured != nil {
		updateParams.IsFeatured = *req.IsFeatured
	} else {
		updateParams.IsFeatured = existingProduct.IsFeatured
	}

	if req.IsPromoted != nil {
		updateParams.IsPromoted = *req.IsPromoted
	} else {
		updateParams.IsPromoted = existingProduct.IsPromoted
	}

	product, err := s.productRepo.UpdateProduct(ctx, updateParams)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to update product", err)
	}

	return s.convertToProductAppDTO(product), nil
}

func (s *productService) UpdateProductStatus(ctx context.Context, id string, status string, sellerID int32) (*dto.ProductAppDTO, error) {
	// First, get the existing product to check ownership
	existingProduct, err := s.productRepo.GetProduct(ctx, id)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Product not found", err)
	}

	// Check if the seller owns this product
	if existingProduct.SellerID != sellerID {
		return nil, response.NewAPIError(http.StatusForbidden, "You don't have permission to update this product", errors.New("unauthorized"))
	}

	product, err := s.productRepo.UpdateProductStatus(ctx, id, status)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to update product status", err)
	}

	return s.convertToProductAppDTO(product), nil
}

func (s *productService) UpdateProductStock(ctx context.Context, id string, stock int32, sellerID int32) (*dto.ProductAppDTO, error) {
	// First, get the existing product to check ownership
	existingProduct, err := s.productRepo.GetProduct(ctx, id)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Product not found", err)
	}

	// Check if the seller owns this product
	if existingProduct.SellerID != sellerID {
		return nil, response.NewAPIError(http.StatusForbidden, "You don't have permission to update this product", errors.New("unauthorized"))
	}

	product, err := s.productRepo.UpdateProductStock(ctx, id, stock)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to update product stock", err)
	}

	return s.convertToProductAppDTO(product), nil
}

func (s *productService) IncrementProductViews(ctx context.Context, id string) error {
	err := s.productRepo.IncrementProductViews(ctx, id)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to increment product views", err)
	}
	return nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string, sellerID int32) error {
	// First, get the existing product to check ownership
	existingProduct, err := s.productRepo.GetProduct(ctx, id)
	if err != nil {
		return response.NewAPIError(http.StatusNotFound, "Product not found", err)
	}

	// Check if the seller owns this product
	if existingProduct.SellerID != sellerID {
		return response.NewAPIError(http.StatusForbidden, "You don't have permission to delete this product", errors.New("unauthorized"))
	}

	err = s.productRepo.DeleteProduct(ctx, id)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to delete product", err)
	}

	return nil
}

func (s *productService) GetProductStats(ctx context.Context) (*dto.ProductStatsAppDTO, error) {
	stats, err := s.productRepo.GetProductStats(ctx)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to get product stats", err)
	}

	return &dto.ProductStatsAppDTO{
		TotalProducts:    stats.TotalProducts,
		ActiveProducts:   stats.ActiveProducts,
		DraftProducts:    stats.DraftProducts,
		SoldProducts:     stats.SoldProducts,
		FeaturedProducts: stats.FeaturedProducts,
		PromotedProducts: stats.PromotedProducts,
		AverageRating:    stats.AverageRating,
		TotalViews:       stats.TotalViews,
	}, nil
}

func (s *productService) GetMyProducts(ctx context.Context, sellerID int32, page, pageSize int32) (*dto.ProductListAppDTO, error) {
	return s.ListProductsBySeller(ctx, sellerID, page, pageSize)
}

func (s *productService) BulkUpdateProducts(ctx context.Context, ids []string, req dto.UpdateProductAppDTO, sellerID int32) ([]dto.ProductAppDTO, error) {
	var results []dto.ProductAppDTO

	for _, id := range ids {
		product, err := s.UpdateProduct(ctx, id, req, sellerID)
		if err != nil {
			return nil, err
		}
		results = append(results, *product)
	}

	return results, nil
}

// Helper methods for conversion
func (s *productService) convertToProductAppDTO(product *entity.Product) *dto.ProductAppDTO {
	return &dto.ProductAppDTO{
		ID:            product.ID,
		Title:         product.Title,
		Description:   product.Description,
		Price:         product.Price,
		OriginalPrice: product.OriginalPrice,
		CategoryID:    product.CategoryID,
		Images:        product.Images,
		Condition:     product.Condition,
		Tags:          product.Tags,
		Location:      product.Location,
		Stock:         product.Stock,
		SellerID:      product.SellerID,
		Status:        product.Status,
		Badges:        product.Badges,
		Rating:        product.Rating,
		ReviewCount:   product.ReviewCount,
		ViewCount:     product.ViewCount,
		IsFeatured:    product.IsFeatured,
		IsPromoted:    product.IsPromoted,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
	}
}

func (s *productService) convertToProductAppDTOWithDetails(product *entity.ProductWithDetails) *dto.ProductAppDTO {
	result := &dto.ProductAppDTO{
		ID:            product.ID,
		Title:         product.Title,
		Description:   product.Description,
		Price:         product.Price,
		OriginalPrice: product.OriginalPrice,
		CategoryID:    product.CategoryID,
		Images:        product.Images,
		Condition:     product.Condition,
		Tags:          product.Tags,
		Location:      product.Location,
		Stock:         product.Stock,
		SellerID:      product.SellerID,
		Status:        product.Status,
		Badges:        product.Badges,
		Rating:        product.Rating,
		ReviewCount:   product.ReviewCount,
		ViewCount:     product.ViewCount,
		IsFeatured:    product.IsFeatured,
		IsPromoted:    product.IsPromoted,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
	}

	// Add category information if available
	if product.CategoryName != nil && product.CategorySlug != nil {
		result.Category = &dto.CategoryAppDTO{
			ID:   product.CategoryID,
			Name: *product.CategoryName,
			Slug: *product.CategorySlug,
		}
	}

	// Add seller information if available
	if product.SellerName != nil {
		result.Seller = &dto.SellerAppDTO{
			ID:   product.SellerID,
			Name: *product.SellerName,
		}
		if product.SellerEmail != nil {
			result.Seller.Email = *product.SellerEmail
		}
	}

	// Add store information if available
	if product.StoreName != nil {
		result.Store = &dto.StoreAppDTO{
			Name: *product.StoreName,
		}
		if product.StoreLogo != nil {
			result.Store.Logo = *product.StoreLogo
		}
		if product.StoreRating != nil {
			result.Store.Rating = *product.StoreRating
		}
	}

	return result
}

func (s *productService) convertToProductAppDTOsWithDetails(products []entity.ProductWithDetails) []dto.ProductAppDTO {
	result := make([]dto.ProductAppDTO, len(products))
	for i, product := range products {
		result[i] = *s.convertToProductAppDTOWithDetails(&product)
	}
	return result
}
