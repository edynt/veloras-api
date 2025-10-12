package http

import (
	"net/http"
	"strconv"

	"github.com/edynt/chogiare/veloras-api/internal/product/application/service"
	appDto "github.com/edynt/chogiare/veloras-api/internal/product/application/service/dto"
	"github.com/edynt/chogiare/veloras-api/internal/product/controller/dto"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/edynt/chogiare/veloras-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct creates a new product
// @Summary Create a new product
// @Description Create a new product for the authenticated seller
// @Tags products
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Product creation data"
// @Success 201 {object} response.Response{data=dto.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products [post]
// @Security BearerAuth
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
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

	// Get seller ID from context (set by auth middleware)
	sellerID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	sellerIDInt, ok := sellerID.(int32)
	if !ok {
		response.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}

	// Convert request to service DTO
	serviceReq := appDto.CreateProductAppDTO{
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
		Status:        req.Status,
		Badges:        req.Badges,
		IsFeatured:    req.IsFeatured,
		IsPromoted:    req.IsPromoted,
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), serviceReq, sellerIDInt)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create product", err)
		return
	}

	response.SuccessResponse(c, h.convertToProductResponse(product))
}

// GetProduct gets a product by ID
// @Summary Get product by ID
// @Description Get a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response.Response{data=dto.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "Product not found", err)
		return
	}

	response.SuccessResponse(c, h.convertToProductResponse(product))
}

// ListProducts lists all products with pagination
// @Summary List all products
// @Description Get a paginated list of all active products
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} response.Response{data=dto.ProductListResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 32)

	products, err := h.productService.ListProducts(c.Request.Context(), int32(page), int32(limit))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to list products", err)
		return
	}

	response.SuccessResponse(c, h.convertToProductListResponse(products))
}

// ListProductsByCategory lists products by category
// @Summary List products by category
// @Description Get a paginated list of products in a specific category
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} response.Response{data=dto.ProductListResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories/{id}/products [get]
func (h *ProductHandler) ListProductsByCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "Category ID is required", nil)
		return
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 32)

	products, err := h.productService.ListProductsByCategory(c.Request.Context(), categoryID, int32(page), int32(limit))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to list products by category", err)
		return
	}

	response.SuccessResponse(c, h.convertToProductListResponse(products))
}

// SearchProducts searches products with filters
// @Summary Search products
// @Description Search products with various filters
// @Tags products
// @Accept json
// @Produce json
// @Param query query string false "Search query"
// @Param categoryId query string false "Category ID"
// @Param minPrice query number false "Minimum price"
// @Param maxPrice query number false "Maximum price"
// @Param condition query string false "Product condition"
// @Param location query string false "Location"
// @Param sellerId query int false "Seller ID"
// @Param rating query number false "Minimum rating"
// @Param sortBy query string false "Sort by field"
// @Param sortOrder query string false "Sort order" Enums(asc, desc)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} response.Response{data=dto.ProductListResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	req := dto.SearchProductsRequest{
		Query:      stringPtr(c.Query("query")),
		CategoryID: stringPtr(c.Query("categoryId")),
		Location:   stringPtr(c.Query("location")),
		Condition:  stringPtr(c.Query("condition")),
		SortBy:     stringPtr(c.Query("sortBy")),
		SortOrder:  stringPtr(c.Query("sortOrder")),
	}

	// Parse numeric parameters
	if minPriceStr := c.Query("minPrice"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			req.MinPrice = &minPrice
		}
	}
	if maxPriceStr := c.Query("maxPrice"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			req.MaxPrice = &maxPrice
		}
	}
	if ratingStr := c.Query("rating"); ratingStr != "" {
		if rating, err := strconv.ParseFloat(ratingStr, 64); err == nil {
			req.Rating = &rating
		}
	}
	if sellerIDStr := c.Query("sellerId"); sellerIDStr != "" {
		if sellerID, err := strconv.ParseInt(sellerIDStr, 10, 32); err == nil {
			sellerIDInt := int32(sellerID)
			req.SellerID = &sellerIDInt
		}
	}
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.ParseInt(pageStr, 10, 32); err == nil {
			pageInt32 := int32(page)
			req.Page = &pageInt32
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limitInt32 := int32(limit)
			req.Limit = &limitInt32
		}
	}

	// Convert to service DTO
	serviceReq := appDto.SearchProductsAppDTO{
		Query:      req.Query,
		CategoryID: req.CategoryID,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Condition:  req.Condition,
		Location:   req.Location,
		SellerID:   req.SellerID,
		Rating:     req.Rating,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		Page:       req.Page,
		Limit:      req.Limit,
	}

	products, err := h.productService.SearchProducts(c.Request.Context(), serviceReq)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to search products", err)
		return
	}

	response.SuccessResponse(c, h.convertToProductListResponse(products))
}

// GetFeaturedProducts gets featured products
// @Summary Get featured products
// @Description Get a list of featured products
// @Tags products
// @Accept json
// @Produce json
// @Param limit query int false "Number of products to return" default(10)
// @Success 200 {object} response.Response{data=[]dto.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/featured [get]
func (h *ProductHandler) GetFeaturedProducts(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)

	products, err := h.productService.GetFeaturedProducts(c.Request.Context(), int32(limit))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to get featured products", err)
		return
	}

	response.SuccessResponse(c, h.convertToProductResponses(products))
}

// GetPromotedProducts gets promoted products
// @Summary Get promoted products
// @Description Get a list of promoted products
// @Tags products
// @Accept json
// @Produce json
// @Param limit query int false "Number of products to return" default(10)
// @Success 200 {object} response.Response{data=[]dto.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/promoted [get]
func (h *ProductHandler) GetPromotedProducts(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)

	products, err := h.productService.GetPromotedProducts(c.Request.Context(), int32(limit))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to get promoted products", err)
		return
	}

	response.SuccessResponse(c, h.convertToProductResponses(products))
}

// UpdateProduct updates a product
// @Summary Update a product
// @Description Update a product (only by the seller who owns it)
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body dto.UpdateProductRequest true "Product update data"
// @Success 200 {object} response.Response{data=dto.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/{id} [put]
// @Security BearerAuth
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	var req dto.UpdateProductRequest
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

	// Get seller ID from context
	sellerID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	sellerIDInt, ok := sellerID.(int32)
	if !ok {
		response.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}

	// Convert request to service DTO
	serviceReq := appDto.UpdateProductAppDTO{
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
		Status:        req.Status,
		Badges:        req.Badges,
		IsFeatured:    req.IsFeatured,
		IsPromoted:    req.IsPromoted,
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), id, serviceReq, sellerIDInt)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update product", err)
		return
	}

	response.SuccessResponse(c, h.convertToProductResponse(product))
}

// UpdateProductStatus updates product status
// @Summary Update product status
// @Description Update the status of a product (only by the seller who owns it)
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body dto.UpdateProductStatusRequest true "Status update data"
// @Success 200 {object} response.Response{data=dto.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /seller/products/{id}/status [patch]
// @Security BearerAuth
func (h *ProductHandler) UpdateProductStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	var req dto.UpdateProductStatusRequest
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

	// Get seller ID from context
	sellerID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	sellerIDInt, ok := sellerID.(int32)
	if !ok {
		response.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}

	product, err := h.productService.UpdateProductStatus(c.Request.Context(), id, req.Status, sellerIDInt)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update product status", err)
		return
	}

	response.SuccessResponse(c, h.convertToProductResponse(product))
}

// UpdateProductStock updates product stock
// @Summary Update product stock
// @Description Update the stock quantity of a product (only by the seller who owns it)
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body dto.UpdateProductStockRequest true "Stock update data"
// @Success 200 {object} response.Response{data=dto.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /seller/products/{id}/stock [patch]
// @Security BearerAuth
func (h *ProductHandler) UpdateProductStock(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	var req dto.UpdateProductStockRequest
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

	// Get seller ID from context
	sellerID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	sellerIDInt, ok := sellerID.(int32)
	if !ok {
		response.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}

	product, err := h.productService.UpdateProductStock(c.Request.Context(), id, req.Stock, sellerIDInt)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update product stock", err)
		return
	}

	response.SuccessResponse(c, h.convertToProductResponse(product))
}

// IncrementProductViews increments product view count
// @Summary Increment product views
// @Description Increment the view count of a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/{id}/views [post]
func (h *ProductHandler) IncrementProductViews(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	err := h.productService.IncrementProductViews(c.Request.Context(), id)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to increment product views", err)
		return
	}

	response.SuccessResponse(c, nil)
}

// DeleteProduct deletes a product
// @Summary Delete a product
// @Description Delete a product (only by the seller who owns it)
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/{id} [delete]
// @Security BearerAuth
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	// Get seller ID from context
	sellerID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	sellerIDInt, ok := sellerID.(int32)
	if !ok {
		response.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}

	err := h.productService.DeleteProduct(c.Request.Context(), id, sellerIDInt)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete product", err)
		return
	}

	response.SuccessResponse(c, nil)
}

// GetProductStats gets product statistics
// @Summary Get product statistics
// @Description Get overall product statistics
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=dto.ProductStatsResponse}
// @Failure 500 {object} response.Response
// @Router /products/stats [get]
func (h *ProductHandler) GetProductStats(c *gin.Context) {
	stats, err := h.productService.GetProductStats(c.Request.Context())
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to get product stats", err)
		return
	}

	response.SuccessResponse(c, h.convertToProductStatsResponse(stats))
}

// GetMyProducts gets products for the authenticated seller
// @Summary Get my products
// @Description Get products for the authenticated seller
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} response.Response{data=dto.ProductListResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /seller/products [get]
// @Security BearerAuth
func (h *ProductHandler) GetMyProducts(c *gin.Context) {
	// Get seller ID from context
	sellerID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	sellerIDInt, ok := sellerID.(int32)
	if !ok {
		response.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 32)

	products, err := h.productService.GetMyProducts(c.Request.Context(), sellerIDInt, int32(page), int32(limit))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to get my products", err)
		return
	}

	response.SuccessResponse(c, h.convertToProductListResponse(products))
}

// BulkUpdateProducts updates multiple products
// @Summary Bulk update products
// @Description Update multiple products at once (only by the seller who owns them)
// @Tags products
// @Accept json
// @Produce json
// @Param request body dto.BulkUpdateProductsRequest true "Bulk update data"
// @Success 200 {object} response.Response{data=[]dto.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /seller/products/bulk [patch]
// @Security BearerAuth
func (h *ProductHandler) BulkUpdateProducts(c *gin.Context) {
	var req dto.BulkUpdateProductsRequest
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

	// Get seller ID from context
	sellerID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	sellerIDInt, ok := sellerID.(int32)
	if !ok {
		response.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}

	// Convert request to service DTO
	serviceReq := appDto.UpdateProductAppDTO{
		Title:         req.Update.Title,
		Description:   req.Update.Description,
		Price:         req.Update.Price,
		OriginalPrice: req.Update.OriginalPrice,
		CategoryID:    req.Update.CategoryID,
		Images:        req.Update.Images,
		Condition:     req.Update.Condition,
		Tags:          req.Update.Tags,
		Location:      req.Update.Location,
		Stock:         req.Update.Stock,
		Status:        req.Update.Status,
		Badges:        req.Update.Badges,
		IsFeatured:    req.Update.IsFeatured,
		IsPromoted:    req.Update.IsPromoted,
	}

	products, err := h.productService.BulkUpdateProducts(c.Request.Context(), req.IDs, serviceReq, sellerIDInt)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to bulk update products", err)
		return
	}

	response.SuccessResponse(c, h.convertToProductResponses(products))
}

// Helper methods for conversion
func (h *ProductHandler) convertToProductResponse(product *appDto.ProductAppDTO) *dto.ProductResponse {
	resp := &dto.ProductResponse{
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
	if product.Category != nil {
		resp.Category = &dto.CategoryResponse{
			ID:           product.Category.ID,
			Name:         product.Category.Name,
			Slug:         product.Category.Slug,
			Description:  product.Category.Description,
			Image:        product.Category.Image,
			ProductCount: product.Category.ProductCount,
			IsActive:     product.Category.IsActive,
		}
	}

	// Add seller information if available
	if product.Seller != nil {
		resp.Seller = &dto.SellerResponse{
			ID:    product.Seller.ID,
			Name:  product.Seller.Name,
			Email: product.Seller.Email,
		}
	}

	// Add store information if available
	if product.Store != nil {
		resp.Store = &dto.StoreResponse{
			ID:     product.Store.ID,
			Name:   product.Store.Name,
			Logo:   product.Store.Logo,
			Rating: product.Store.Rating,
		}
	}

	return resp
}

func (h *ProductHandler) convertToProductResponses(products []appDto.ProductAppDTO) []dto.ProductResponse {
	result := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		result[i] = *h.convertToProductResponse(&product)
	}
	return result
}

func (h *ProductHandler) convertToProductListResponse(products *appDto.ProductListAppDTO) *dto.ProductListResponse {
	return &dto.ProductListResponse{
		Products:   h.convertToProductResponses(products.Products),
		Total:      products.Total,
		Page:       products.Page,
		PageSize:   products.PageSize,
		TotalPages: products.TotalPages,
	}
}

// Helper function to convert string to *string
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (h *ProductHandler) convertToProductStatsResponse(stats *appDto.ProductStatsAppDTO) *dto.ProductStatsResponse {
	return &dto.ProductStatsResponse{
		TotalProducts:    stats.TotalProducts,
		ActiveProducts:   stats.ActiveProducts,
		DraftProducts:    stats.DraftProducts,
		SoldProducts:     stats.SoldProducts,
		FeaturedProducts: stats.FeaturedProducts,
		PromotedProducts: stats.PromotedProducts,
		AverageRating:    stats.AverageRating,
		TotalViews:       stats.TotalViews,
	}
}
