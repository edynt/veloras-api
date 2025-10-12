package repository

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/category/domain/model/entity"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, req CreateCategoryParams) (*entity.Category, error)
	GetCategory(ctx context.Context, id string) (*entity.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*entity.Category, error)
	ListCategories(ctx context.Context) ([]entity.Category, error)
	ListCategoriesWithPagination(ctx context.Context, limit, offset int32) ([]entity.Category, error)
	ListSubCategories(ctx context.Context, parentID string) ([]entity.Category, error)
	UpdateCategory(ctx context.Context, req UpdateCategoryParams) (*entity.Category, error)
	DeleteCategory(ctx context.Context, id string) error
	UpdateCategoryProductCount(ctx context.Context, id string, increment int32) error
	GetCategoryStats(ctx context.Context) (*entity.CategoryStats, error)
}

type CreateCategoryParams struct {
	Name         string
	Slug         string
	Description  *string
	Image        *string
	ParentID     *string
	ProductCount int32
	IsActive     bool
}

type UpdateCategoryParams struct {
	ID           string
	Name         string
	Slug         string
	Description  *string
	Image        *string
	ParentID     *string
	ProductCount int32
	IsActive     bool
}
