package repository

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/category/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/category/domain/repository"
	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type categoryRepository struct {
	db      *pgxpool.Pool
	queries *gen.Queries
}

func NewCategoryRepository(db *pgxpool.Pool) repository.CategoryRepository {
	return &categoryRepository{
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

func stringPtrToText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func stringPtrToUUID(s *string) pgtype.UUID {
	if s == nil {
		return pgtype.UUID{Valid: false}
	}
	return stringToUUID(*s)
}

func int32ToInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{Int32: i, Valid: true}
}

func boolToPgBool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

func (r *categoryRepository) CreateCategory(ctx context.Context, req repository.CreateCategoryParams) (*entity.Category, error) {
	category, err := r.queries.CreateCategory(ctx, gen.CreateCategoryParams{
		Name:         req.Name,
		Slug:         req.Slug,
		Description:  stringPtrToText(req.Description),
		Image:        stringPtrToText(req.Image),
		ParentID:     stringPtrToUUID(req.ParentID),
		ProductCount: int32ToInt4(req.ProductCount),
		IsActive:     boolToPgBool(req.IsActive),
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCCategory(category)
	return &result, nil
}

func (r *categoryRepository) GetCategory(ctx context.Context, id string) (*entity.Category, error) {
	category, err := r.queries.GetCategory(ctx, stringToUUID(id))
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCCategory(category)
	return &result, nil
}

func (r *categoryRepository) GetCategoryBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	category, err := r.queries.GetCategoryBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCCategory(category)
	return &result, nil
}

func (r *categoryRepository) ListCategories(ctx context.Context) ([]entity.Category, error) {
	categories, err := r.queries.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]entity.Category, len(categories))
	for i, c := range categories {
		result[i] = entity.FromSQLCCategory(c)
	}

	return result, nil
}

func (r *categoryRepository) ListCategoriesWithPagination(ctx context.Context, limit, offset int32) ([]entity.Category, error) {
	categories, err := r.queries.ListCategoriesWithPagination(ctx, gen.ListCategoriesWithPaginationParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.Category, len(categories))
	for i, c := range categories {
		result[i] = entity.FromSQLCCategory(c)
	}

	return result, nil
}

func (r *categoryRepository) ListSubCategories(ctx context.Context, parentID string) ([]entity.Category, error) {
	categories, err := r.queries.ListSubCategories(ctx, stringToUUID(parentID))
	if err != nil {
		return nil, err
	}

	result := make([]entity.Category, len(categories))
	for i, c := range categories {
		result[i] = entity.FromSQLCCategory(c)
	}

	return result, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, req repository.UpdateCategoryParams) (*entity.Category, error) {
	category, err := r.queries.UpdateCategory(ctx, gen.UpdateCategoryParams{
		ID:           stringToUUID(req.ID),
		Name:         req.Name,
		Slug:         req.Slug,
		Description:  stringPtrToText(req.Description),
		Image:        stringPtrToText(req.Image),
		ParentID:     stringPtrToUUID(req.ParentID),
		ProductCount: int32ToInt4(req.ProductCount),
		IsActive:     boolToPgBool(req.IsActive),
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCCategory(category)
	return &result, nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id string) error {
	return r.queries.DeleteCategory(ctx, stringToUUID(id))
}

func (r *categoryRepository) UpdateCategoryProductCount(ctx context.Context, id string, increment int32) error {
	return r.queries.UpdateCategoryProductCount(ctx, gen.UpdateCategoryProductCountParams{
		ID:           stringToUUID(id),
		ProductCount: int32ToInt4(increment),
	})
}

func (r *categoryRepository) GetCategoryStats(ctx context.Context) (*entity.CategoryStats, error) {
	stats, err := r.queries.GetCategoryStats(ctx)
	if err != nil {
		return nil, err
	}

	result := &entity.CategoryStats{
		TotalCategories:  stats.TotalCategories,
		ParentCategories: stats.ParentCategories,
		SubCategories:    stats.SubCategories,
		TotalProducts:    stats.TotalProducts,
	}

	return result, nil
}
