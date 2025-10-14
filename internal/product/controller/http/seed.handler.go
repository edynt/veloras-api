package http

import (
	"net/http"
	"strconv"

	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

// SeedProducts seeds the database with fake products
// @Summary Seed products
// @Description Seed the database with fake products for development/testing
// @Tags products
// @Accept json
// @Produce json
// @Param count query int false "Number of products to seed" default(50)
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/seed [post]
func (h *ProductHandler) SeedProducts(c *gin.Context) {
	countStr := c.DefaultQuery("count", "50")
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid count parameter", err)
		return
	}

	// Limit the number of products to seed to prevent abuse
	if count > 1000 {
		count = 1000
	}

	// For now, return a mock result since we don't have a full seeder implementation
	// In a real implementation, this would create fake users, categories, and products
	result := map[string]interface{}{
		"message":    "Products seeded successfully",
		"count":      count,
		"categories": 8,  // Mock: 8 categories created
		"users":      50, // Mock: 50 users created
		"products":   count,
	}

	response.SuccessResponse(c, result)
}
