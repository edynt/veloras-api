package http

import (
	"net/http"
	"strconv"

	"github.com/edynt/chogiare/veloras-api/internal/category/application/service"
	appDto "github.com/edynt/chogiare/veloras-api/internal/category/application/service/dto"
	ctlDto "github.com/edynt/chogiare/veloras-api/internal/category/controller/dto"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/edynt/chogiare/veloras-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory creates a new category
// @Summary Create a new category
// @Description Create a new product category
// @Tags categories
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoryRequest true "Category creation data"
// @Success 201 {object} response.Response{data=dto.CategoryResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req ctlDto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewAPIError(http.StatusBadRequest, "Invalid request data", err.Error()))
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req, nil); err != nil {
		c.JSON(http.StatusBadRequest, response.NewAPIError(http.StatusBadRequest, "Validation failed", err.Error()))
		return
	}

	// Convert request to service DTO
	serviceReq := appDto.CreateCategoryAppDTO{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Image:       req.Image,
		ParentID:    req.ParentID,
		IsActive:    req.IsActive,
	}

	category, err := h.categoryService.CreateCategory(c.Request.Context(), serviceReq)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create category", err)
		return
	}

	response.SuccessResponse(c, h.convertToCategoryResponse(category))
}

// GetCategory gets a category by ID
// @Summary Get category by ID
// @Description Get a category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} response.Response{data=dto.CategoryResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "Category ID is required", nil)
		return
	}

	category, err := h.categoryService.GetCategory(c.Request.Context(), id)
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "Category not found", err)
		return
	}

	response.SuccessResponse(c, h.convertToCategoryResponse(category))
}

// GetCategoryBySlug gets a category by slug
// @Summary Get category by slug
// @Description Get a category by its slug
// @Tags categories
// @Accept json
// @Produce json
// @Param slug path string true "Category slug"
// @Success 200 {object} response.Response{data=dto.CategoryResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories/slug/{slug} [get]
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "Category slug is required", nil)
		return
	}

	category, err := h.categoryService.GetCategoryBySlug(c.Request.Context(), slug)
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "Category not found", err)
		return
	}

	response.SuccessResponse(c, h.convertToCategoryResponse(category))
}

// ListCategories lists all categories
// @Summary List all categories
// @Description Get a list of all active categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.CategoryResponse}
// @Failure 500 {object} response.Response
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.categoryService.ListCategories(c.Request.Context())
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to list categories", err)
		return
	}

	response.SuccessResponse(c, h.convertToCategoryResponses(categories))
}

// ListCategoriesWithPagination lists categories with pagination
// @Summary List categories with pagination
// @Description Get a paginated list of categories
// @Tags categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} response.Response{data=dto.CategoryListResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories/paginated [get]
func (h *CategoryHandler) ListCategoriesWithPagination(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 32)

	categories, err := h.categoryService.ListCategoriesWithPagination(c.Request.Context(), int32(page), int32(limit))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to list categories", err)
		return
	}

	response.SuccessResponse(c, h.convertToCategoryListResponse(categories))
}

// ListSubCategories lists subcategories
// @Summary List subcategories
// @Description Get a list of subcategories for a parent category
// @Tags categories
// @Accept json
// @Produce json
// @Param parentId path string true "Parent category ID"
// @Success 200 {object} response.Response{data=[]dto.CategoryResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories/{parentId}/subcategories [get]
func (h *CategoryHandler) ListSubCategories(c *gin.Context) {
	parentID := c.Param("parentId")
	if parentID == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "Parent category ID is required", nil)
		return
	}

	categories, err := h.categoryService.ListSubCategories(c.Request.Context(), parentID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to list subcategories", err)
		return
	}

	response.SuccessResponse(c, h.convertToCategoryResponses(categories))
}

// UpdateCategory updates a category
// @Summary Update a category
// @Description Update a category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Category update data"
// @Success 200 {object} response.Response{data=dto.CategoryResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "Category ID is required", nil)
		return
	}

	var req ctlDto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// Validate request
	validation, exists := c.Get("validation")
	if !exists {
		response.ErrorResponse(c, http.StatusInternalServerError, "Validation middleware not found", nil)
		return
	}
	if apiErr := utils.ValidateStruct(req, validation.(*validator.Validate)); apiErr != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Validation failed", apiErr)
		return
	}

	// Convert request to service DTO
	serviceReq := appDto.UpdateCategoryAppDTO{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Image:       req.Image,
		ParentID:    req.ParentID,
		IsActive:    req.IsActive,
	}

	category, err := h.categoryService.UpdateCategory(c.Request.Context(), id, serviceReq)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update category", err)
		return
	}

	response.SuccessResponse(c, h.convertToCategoryResponse(category))
}

// DeleteCategory deletes a category
// @Summary Delete a category
// @Description Delete a category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "Category ID is required", nil)
		return
	}

	err := h.categoryService.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete category", err)
		return
	}

	response.SuccessResponse(c, nil)
}

// GetCategoryStats gets category statistics
// @Summary Get category statistics
// @Description Get overall category statistics
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=dto.CategoryStatsResponse}
// @Failure 500 {object} response.Response
// @Router /categories/stats [get]
func (h *CategoryHandler) GetCategoryStats(c *gin.Context) {
	stats, err := h.categoryService.GetCategoryStats(c.Request.Context())
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to get category stats", err)
		return
	}

	response.SuccessResponse(c, h.convertToCategoryStatsResponse(stats))
}

// Helper methods for conversion
func (h *CategoryHandler) convertToCategoryResponse(category *appDto.CategoryAppDTO) *ctlDto.CategoryResponse {
	return &ctlDto.CategoryResponse{
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

func (h *CategoryHandler) convertToCategoryResponses(categories []appDto.CategoryAppDTO) []ctlDto.CategoryResponse {
	result := make([]ctlDto.CategoryResponse, len(categories))
	for i, category := range categories {
		result[i] = *h.convertToCategoryResponse(&category)
	}
	return result
}

func (h *CategoryHandler) convertToCategoryListResponse(categories *appDto.CategoryListAppDTO) *ctlDto.CategoryListResponse {
	return &ctlDto.CategoryListResponse{
		Categories: h.convertToCategoryResponses(categories.Categories),
		Total:      categories.Total,
		Page:       categories.Page,
		PageSize:   categories.PageSize,
		TotalPages: categories.TotalPages,
	}
}

func (h *CategoryHandler) convertToCategoryStatsResponse(stats *appDto.CategoryStatsAppDTO) *ctlDto.CategoryStatsResponse {
	return &ctlDto.CategoryStatsResponse{
		TotalCategories:  stats.TotalCategories,
		ParentCategories: stats.ParentCategories,
		SubCategories:    stats.SubCategories,
		TotalProducts:    stats.TotalProducts,
	}
}
