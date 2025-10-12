package service

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/category/application/service/dto"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req dto.CreateCategoryAppDTO) (*dto.CategoryAppDTO, error)
	GetCategory(ctx context.Context, id string) (*dto.CategoryAppDTO, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*dto.CategoryAppDTO, error)
	ListCategories(ctx context.Context) ([]dto.CategoryAppDTO, error)
	ListCategoriesWithPagination(ctx context.Context, page, pageSize int32) (*dto.CategoryListAppDTO, error)
	ListSubCategories(ctx context.Context, parentID string) ([]dto.CategoryAppDTO, error)
	UpdateCategory(ctx context.Context, id string, req dto.UpdateCategoryAppDTO) (*dto.CategoryAppDTO, error)
	DeleteCategory(ctx context.Context, id string) error
	GetCategoryStats(ctx context.Context) (*dto.CategoryStatsAppDTO, error)
}
