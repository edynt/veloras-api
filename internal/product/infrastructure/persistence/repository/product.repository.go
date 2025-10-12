package repository

import (
	"context"
	"math/big"

	"github.com/edynt/chogiare/veloras-api/internal/product/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/product/domain/repository"
	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepository struct {
	db      *pgxpool.Pool
	queries *gen.Queries
}

func NewProductRepository(db *pgxpool.Pool) repository.ProductRepository {
	return &productRepository{
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

func float64ToNumeric(f float64) pgtype.Numeric {
	return pgtype.Numeric{Int: big.NewInt(int64(f * 100)), Valid: true, Exp: -2} // Assuming 2 decimal places
}

func float64PtrToNumeric(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{Valid: false}
	}
	return float64ToNumeric(*f)
}

func boolToPgBool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

func (r *productRepository) CreateProduct(ctx context.Context, req repository.CreateProductParams) (*entity.Product, error) {
	product, err := r.queries.CreateProduct(ctx, gen.CreateProductParams{
		Title:         req.Title,
		Description:   req.Description,
		Price:         float64ToNumeric(req.Price),
		OriginalPrice: float64PtrToNumeric(req.OriginalPrice),
		CategoryID:    stringToUUID(req.CategoryID),
		Images:        req.Images,
		Condition:     req.Condition,
		Tags:          req.Tags,
		Location:      req.Location,
		Stock:         req.Stock,
		SellerID:      req.SellerID,
		Status:        req.Status,
		Badges:        req.Badges,
		IsFeatured:    boolToPgBool(req.IsFeatured),
		IsPromoted:    boolToPgBool(req.IsPromoted),
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCProduct(product)
	return &result, nil
}

func (r *productRepository) GetProduct(ctx context.Context, id string) (*entity.ProductWithDetails, error) {
	product, err := r.queries.GetProduct(ctx, stringToUUID(id))
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCProductWithDetails(product)
	return &result, nil
}

func (r *productRepository) ListProducts(ctx context.Context, limit, offset int32) ([]entity.ProductWithDetails, error) {
	products, err := r.queries.ListProducts(ctx, gen.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.ProductWithDetails, len(products))
	for i, p := range products {
		result[i] = entity.FromSQLCProductWithDetails(gen.GetProductRow{
			ID:            p.ID,
			Title:         p.Title,
			Description:   p.Description,
			Price:         p.Price,
			OriginalPrice: p.OriginalPrice,
			CategoryID:    p.CategoryID,
			Images:        p.Images,
			Condition:     p.Condition,
			Tags:          p.Tags,
			Location:      p.Location,
			Stock:         p.Stock,
			SellerID:      p.SellerID,
			Status:        p.Status,
			Badges:        p.Badges,
			Rating:        p.Rating,
			ReviewCount:   p.ReviewCount,
			ViewCount:     p.ViewCount,
			IsFeatured:    p.IsFeatured,
			IsPromoted:    p.IsPromoted,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			CategoryName:  p.CategoryName,
			CategorySlug:  p.CategorySlug,
			SellerName:    p.SellerName,
			StoreName:     p.StoreName,
			StoreLogo:     p.StoreLogo,
			StoreRating:   p.StoreRating,
		})
	}

	return result, nil
}

func (r *productRepository) ListProductsByCategory(ctx context.Context, categoryID string, limit, offset int32) ([]entity.ProductWithDetails, error) {
	products, err := r.queries.ListProductsByCategory(ctx, gen.ListProductsByCategoryParams{
		CategoryID: stringToUUID(categoryID),
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.ProductWithDetails, len(products))
	for i, p := range products {
		result[i] = entity.FromSQLCProductWithDetails(gen.GetProductRow{
			ID:            p.ID,
			Title:         p.Title,
			Description:   p.Description,
			Price:         p.Price,
			OriginalPrice: p.OriginalPrice,
			CategoryID:    p.CategoryID,
			Images:        p.Images,
			Condition:     p.Condition,
			Tags:          p.Tags,
			Location:      p.Location,
			Stock:         p.Stock,
			SellerID:      p.SellerID,
			Status:        p.Status,
			Badges:        p.Badges,
			Rating:        p.Rating,
			ReviewCount:   p.ReviewCount,
			ViewCount:     p.ViewCount,
			IsFeatured:    p.IsFeatured,
			IsPromoted:    p.IsPromoted,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			CategoryName:  p.CategoryName,
			CategorySlug:  p.CategorySlug,
			SellerName:    p.SellerName,
			StoreName:     p.StoreName,
			StoreLogo:     p.StoreLogo,
			StoreRating:   p.StoreRating,
		})
	}

	return result, nil
}

func (r *productRepository) ListProductsBySeller(ctx context.Context, sellerID int32, limit, offset int32) ([]entity.ProductWithDetails, error) {
	products, err := r.queries.ListProductsBySeller(ctx, gen.ListProductsBySellerParams{
		SellerID: sellerID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.ProductWithDetails, len(products))
	for i, p := range products {
		result[i] = entity.FromSQLCProductWithDetails(gen.GetProductRow{
			ID:            p.ID,
			Title:         p.Title,
			Description:   p.Description,
			Price:         p.Price,
			OriginalPrice: p.OriginalPrice,
			CategoryID:    p.CategoryID,
			Images:        p.Images,
			Condition:     p.Condition,
			Tags:          p.Tags,
			Location:      p.Location,
			Stock:         p.Stock,
			SellerID:      p.SellerID,
			Status:        p.Status,
			Badges:        p.Badges,
			Rating:        p.Rating,
			ReviewCount:   p.ReviewCount,
			ViewCount:     p.ViewCount,
			IsFeatured:    p.IsFeatured,
			IsPromoted:    p.IsPromoted,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			CategoryName:  p.CategoryName,
			CategorySlug:  p.CategorySlug,
			SellerName:    p.SellerName,
			StoreName:     p.StoreName,
			StoreLogo:     p.StoreLogo,
			StoreRating:   p.StoreRating,
		})
	}

	return result, nil
}

func (r *productRepository) SearchProducts(ctx context.Context, query string, limit, offset int32) ([]entity.ProductWithDetails, error) {
	products, err := r.queries.SearchProducts(ctx, gen.SearchProductsParams{
		PlaintoTsquery: query,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.ProductWithDetails, len(products))
	for i, p := range products {
		result[i] = entity.FromSQLCProductWithDetails(gen.GetProductRow{
			ID:            p.ID,
			Title:         p.Title,
			Description:   p.Description,
			Price:         p.Price,
			OriginalPrice: p.OriginalPrice,
			CategoryID:    p.CategoryID,
			Images:        p.Images,
			Condition:     p.Condition,
			Tags:          p.Tags,
			Location:      p.Location,
			Stock:         p.Stock,
			SellerID:      p.SellerID,
			Status:        p.Status,
			Badges:        p.Badges,
			Rating:        p.Rating,
			ReviewCount:   p.ReviewCount,
			ViewCount:     p.ViewCount,
			IsFeatured:    p.IsFeatured,
			IsPromoted:    p.IsPromoted,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			CategoryName:  p.CategoryName,
			CategorySlug:  p.CategorySlug,
			SellerName:    p.SellerName,
			StoreName:     p.StoreName,
			StoreLogo:     p.StoreLogo,
			StoreRating:   p.StoreRating,
		})
	}

	return result, nil
}

func (r *productRepository) GetFeaturedProducts(ctx context.Context, limit int32) ([]entity.ProductWithDetails, error) {
	products, err := r.queries.GetFeaturedProducts(ctx, limit)
	if err != nil {
		return nil, err
	}

	result := make([]entity.ProductWithDetails, len(products))
	for i, p := range products {
		result[i] = entity.FromSQLCProductWithDetails(gen.GetProductRow{
			ID:            p.ID,
			Title:         p.Title,
			Description:   p.Description,
			Price:         p.Price,
			OriginalPrice: p.OriginalPrice,
			CategoryID:    p.CategoryID,
			Images:        p.Images,
			Condition:     p.Condition,
			Tags:          p.Tags,
			Location:      p.Location,
			Stock:         p.Stock,
			SellerID:      p.SellerID,
			Status:        p.Status,
			Badges:        p.Badges,
			Rating:        p.Rating,
			ReviewCount:   p.ReviewCount,
			ViewCount:     p.ViewCount,
			IsFeatured:    p.IsFeatured,
			IsPromoted:    p.IsPromoted,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			CategoryName:  p.CategoryName,
			CategorySlug:  p.CategorySlug,
			SellerName:    p.SellerName,
			StoreName:     p.StoreName,
			StoreLogo:     p.StoreLogo,
			StoreRating:   p.StoreRating,
		})
	}

	return result, nil
}

func (r *productRepository) GetPromotedProducts(ctx context.Context, limit int32) ([]entity.ProductWithDetails, error) {
	products, err := r.queries.GetPromotedProducts(ctx, limit)
	if err != nil {
		return nil, err
	}

	result := make([]entity.ProductWithDetails, len(products))
	for i, p := range products {
		result[i] = entity.FromSQLCProductWithDetails(gen.GetProductRow{
			ID:            p.ID,
			Title:         p.Title,
			Description:   p.Description,
			Price:         p.Price,
			OriginalPrice: p.OriginalPrice,
			CategoryID:    p.CategoryID,
			Images:        p.Images,
			Condition:     p.Condition,
			Tags:          p.Tags,
			Location:      p.Location,
			Stock:         p.Stock,
			SellerID:      p.SellerID,
			Status:        p.Status,
			Badges:        p.Badges,
			Rating:        p.Rating,
			ReviewCount:   p.ReviewCount,
			ViewCount:     p.ViewCount,
			IsFeatured:    p.IsFeatured,
			IsPromoted:    p.IsPromoted,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			CategoryName:  p.CategoryName,
			CategorySlug:  p.CategorySlug,
			SellerName:    p.SellerName,
			StoreName:     p.StoreName,
			StoreLogo:     p.StoreLogo,
			StoreRating:   p.StoreRating,
		})
	}

	return result, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, req repository.UpdateProductParams) (*entity.Product, error) {
	product, err := r.queries.UpdateProduct(ctx, gen.UpdateProductParams{
		ID:            stringToUUID(req.ID),
		Title:         req.Title,
		Description:   req.Description,
		Price:         float64ToNumeric(req.Price),
		OriginalPrice: float64PtrToNumeric(req.OriginalPrice),
		CategoryID:    stringToUUID(req.CategoryID),
		Images:        req.Images,
		Condition:     req.Condition,
		Tags:          req.Tags,
		Location:      req.Location,
		Stock:         req.Stock,
		Status:        req.Status,
		Badges:        req.Badges,
		IsFeatured:    boolToPgBool(req.IsFeatured),
		IsPromoted:    boolToPgBool(req.IsPromoted),
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCProduct(product)
	return &result, nil
}

func (r *productRepository) UpdateProductStatus(ctx context.Context, id string, status string) (*entity.Product, error) {
	product, err := r.queries.UpdateProductStatus(ctx, gen.UpdateProductStatusParams{
		ID:     stringToUUID(id),
		Status: status,
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCProduct(product)
	return &result, nil
}

func (r *productRepository) UpdateProductStock(ctx context.Context, id string, stock int32) (*entity.Product, error) {
	product, err := r.queries.UpdateProductStock(ctx, gen.UpdateProductStockParams{
		ID:    stringToUUID(id),
		Stock: stock,
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCProduct(product)
	return &result, nil
}

func (r *productRepository) IncrementProductViews(ctx context.Context, id string) error {
	return r.queries.IncrementProductViews(ctx, stringToUUID(id))
}

func (r *productRepository) UpdateProductRating(ctx context.Context, id string, rating float64, reviewCount int32) error {
	return r.queries.UpdateProductRating(ctx, gen.UpdateProductRatingParams{
		ID:          stringToUUID(id),
		Rating:      float64ToNumeric(rating),
		ReviewCount: pgtype.Int4{Int32: reviewCount, Valid: true},
	})
}

func (r *productRepository) DeleteProduct(ctx context.Context, id string) error {
	return r.queries.DeleteProduct(ctx, stringToUUID(id))
}

func (r *productRepository) GetProductStats(ctx context.Context) (*entity.ProductStats, error) {
	stats, err := r.queries.GetProductStats(ctx)
	if err != nil {
		return nil, err
	}

	result := &entity.ProductStats{
		TotalProducts:    stats.TotalProducts,
		ActiveProducts:   stats.ActiveProducts,
		DraftProducts:    stats.DraftProducts,
		SoldProducts:     stats.SoldProducts,
		FeaturedProducts: stats.FeaturedProducts,
		PromotedProducts: stats.PromotedProducts,
		AverageRating:    &stats.AverageRating,
		TotalViews:       &stats.TotalViews,
	}

	return result, nil
}

func (r *productRepository) GetProductsByFilters(ctx context.Context, req repository.GetProductsByFiltersParams) ([]entity.ProductWithDetails, error) {
	// Convert parameters to the expected types
	var categoryID pgtype.UUID
	if req.CategoryID != nil {
		categoryID = stringToUUID(*req.CategoryID)
	}

	var minPrice pgtype.Numeric
	if req.MinPrice != nil {
		minPrice = float64ToNumeric(*req.MinPrice)
	}

	var maxPrice pgtype.Numeric
	if req.MaxPrice != nil {
		maxPrice = float64ToNumeric(*req.MaxPrice)
	}

	var condition string
	if req.Condition != nil {
		condition = *req.Condition
	}

	var location string
	if req.Location != nil {
		location = *req.Location
	}

	var sellerID int32
	if req.SellerID != nil {
		sellerID = *req.SellerID
	}

	var rating pgtype.Numeric
	if req.Rating != nil {
		rating = float64ToNumeric(*req.Rating)
	}

	var sortBy interface{}
	if req.SortBy != nil {
		sortBy = *req.SortBy
	}

	products, err := r.queries.GetProductsByFilters(ctx, gen.GetProductsByFiltersParams{
		Column1: categoryID,
		Column2: minPrice,
		Column3: maxPrice,
		Column4: condition,
		Column5: location,
		Column6: sellerID,
		Column7: rating,
		Column8: sortBy,
		Limit:   req.Limit,
		Offset:  req.Offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.ProductWithDetails, len(products))
	for i, p := range products {
		result[i] = entity.FromSQLCProductWithDetails(gen.GetProductRow{
			ID:            p.ID,
			Title:         p.Title,
			Description:   p.Description,
			Price:         p.Price,
			OriginalPrice: p.OriginalPrice,
			CategoryID:    p.CategoryID,
			Images:        p.Images,
			Condition:     p.Condition,
			Tags:          p.Tags,
			Location:      p.Location,
			Stock:         p.Stock,
			SellerID:      p.SellerID,
			Status:        p.Status,
			Badges:        p.Badges,
			Rating:        p.Rating,
			ReviewCount:   p.ReviewCount,
			ViewCount:     p.ViewCount,
			IsFeatured:    p.IsFeatured,
			IsPromoted:    p.IsPromoted,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			CategoryName:  p.CategoryName,
			CategorySlug:  p.CategorySlug,
			SellerName:    p.SellerName,
			StoreName:     p.StoreName,
			StoreLogo:     p.StoreLogo,
			StoreRating:   p.StoreRating,
		})
	}

	return result, nil
}
