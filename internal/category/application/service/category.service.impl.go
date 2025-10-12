package service

import (
	"context"
	"math"
	"net/http"

	"github.com/edynt/chogiare/veloras-api/internal/category/application/service/dto"
	"github.com/edynt/chogiare/veloras-api/internal/category/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/category/domain/repository"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
)

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, req dto.CreateCategoryAppDTO) (*dto.CategoryAppDTO, error) {
	category, err := s.categoryRepo.CreateCategory(ctx, repository.CreateCategoryParams{
		Name:         req.Name,
		Slug:         req.Slug,
		Description:  req.Description,
		Image:        req.Image,
		ParentID:     req.ParentID,
		ProductCount: 0,
		IsActive:     req.IsActive,
	})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to create category", err)
	}

	return s.convertToCategoryAppDTO(category), nil
}

func (s *categoryService) GetCategory(ctx context.Context, id string) (*dto.CategoryAppDTO, error) {
	category, err := s.categoryRepo.GetCategory(ctx, id)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Category not found", err)
	}

	return s.convertToCategoryAppDTO(category), nil
}

func (s *categoryService) GetCategoryBySlug(ctx context.Context, slug string) (*dto.CategoryAppDTO, error) {
	category, err := s.categoryRepo.GetCategoryBySlug(ctx, slug)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Category not found", err)
	}

	return s.convertToCategoryAppDTO(category), nil
}

func (s *categoryService) ListCategories(ctx context.Context) ([]dto.CategoryAppDTO, error) {
	categories, err := s.categoryRepo.ListCategories(ctx)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list categories", err)
	}

	return s.convertToCategoryAppDTOs(categories), nil
}

func (s *categoryService) ListCategoriesWithPagination(ctx context.Context, page, pageSize int32) (*dto.CategoryListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	categories, err := s.categoryRepo.ListCategoriesWithPagination(ctx, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list categories", err)
	}

	// For now, we'll return the categories without total count
	// In a real implementation, you'd want to add a count query
	total := int64(len(categories))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.CategoryListAppDTO{
		Categories: s.convertToCategoryAppDTOs(categories),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *categoryService) ListSubCategories(ctx context.Context, parentID string) ([]dto.CategoryAppDTO, error) {
	categories, err := s.categoryRepo.ListSubCategories(ctx, parentID)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list subcategories", err)
	}

	return s.convertToCategoryAppDTOs(categories), nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id string, req dto.UpdateCategoryAppDTO) (*dto.CategoryAppDTO, error) {
	// First, get the existing category
	existingCategory, err := s.categoryRepo.GetCategory(ctx, id)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Category not found", err)
	}

	// Build update parameters
	updateParams := repository.UpdateCategoryParams{
		ID: id,
	}

	// Only update fields that are provided
	if req.Name != nil {
		updateParams.Name = *req.Name
	} else {
		updateParams.Name = existingCategory.Name
	}

	if req.Slug != nil {
		updateParams.Slug = *req.Slug
	} else {
		updateParams.Slug = existingCategory.Slug
	}

	if req.Description != nil {
		updateParams.Description = req.Description
	} else {
		updateParams.Description = existingCategory.Description
	}

	if req.Image != nil {
		updateParams.Image = req.Image
	} else {
		updateParams.Image = existingCategory.Image
	}

	if req.ParentID != nil {
		updateParams.ParentID = req.ParentID
	} else {
		updateParams.ParentID = existingCategory.ParentID
	}

	if req.IsActive != nil {
		updateParams.IsActive = *req.IsActive
	} else {
		updateParams.IsActive = existingCategory.IsActive
	}

	updateParams.ProductCount = existingCategory.ProductCount

	category, err := s.categoryRepo.UpdateCategory(ctx, updateParams)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to update category", err)
	}

	return s.convertToCategoryAppDTO(category), nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id string) error {
	err := s.categoryRepo.DeleteCategory(ctx, id)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to delete category", err)
	}

	return nil
}

func (s *categoryService) GetCategoryStats(ctx context.Context) (*dto.CategoryStatsAppDTO, error) {
	stats, err := s.categoryRepo.GetCategoryStats(ctx)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to get category stats", err)
	}

	return &dto.CategoryStatsAppDTO{
		TotalCategories:  stats.TotalCategories,
		ParentCategories: stats.ParentCategories,
		SubCategories:    stats.SubCategories,
		TotalProducts:    stats.TotalProducts,
	}, nil
}

// Helper methods for conversion
func (s *categoryService) convertToCategoryAppDTO(category *entity.Category) *dto.CategoryAppDTO {
	return &dto.CategoryAppDTO{
		ID:           category.ID,
		Name:         category.Name,
		Slug:         category.Slug,
		Description:  category.Description,
		Image:        category.Image,
		ParentID:     category.ParentID,
		ProductCount: category.ProductCount,
		IsActive:     category.IsActive,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}
}

func (s *categoryService) convertToCategoryAppDTOs(categories []entity.Category) []dto.CategoryAppDTO {
	result := make([]dto.CategoryAppDTO, len(categories))
	for i, category := range categories {
		result[i] = *s.convertToCategoryAppDTO(&category)
	}
	return result
}
