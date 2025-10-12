package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(router *gin.RouterGroup, categoryHandler *CategoryHandler) {
	categories := router.Group("/categories")
	{
		categories.GET("", categoryHandler.ListCategories)
		categories.GET("/paginated", categoryHandler.ListCategoriesWithPagination)
		categories.GET("/stats", categoryHandler.GetCategoryStats)
		categories.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
		categories.GET("/:id", categoryHandler.GetCategory)
		categories.GET("/:parentId/subcategories", categoryHandler.ListSubCategories)
		categories.POST("", categoryHandler.CreateCategory)
		categories.PUT("/:id", categoryHandler.UpdateCategory)
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
	}
}
