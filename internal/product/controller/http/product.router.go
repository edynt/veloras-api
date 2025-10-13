package http

import (
	"github.com/edynt/chogiare/veloras-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(router *gin.RouterGroup, productHandler *ProductHandler) {
	// Public routes
	products := router.Group("/products")
	{
		products.GET("", productHandler.ListProducts)
		products.GET("/search", productHandler.SearchProducts)
		products.GET("/featured", productHandler.GetFeaturedProducts)
		products.GET("/promoted", productHandler.GetPromotedProducts)
		products.GET("/stats", productHandler.GetProductStats)
		products.GET("/:id", productHandler.GetProduct)
		products.POST("/:id/views", productHandler.IncrementProductViews)
	}

	// Protected routes (require authentication)
	productsProtected := router.Group("/products")
	productsProtected.Use(middleware.AuthenMiddleware())
	{
		productsProtected.POST("", productHandler.CreateProduct)
		productsProtected.PUT("/:id", productHandler.UpdateProduct)
		productsProtected.DELETE("/:id", productHandler.DeleteProduct)
	}

	// Seller-specific routes (require seller role)
	seller := router.Group("/seller")
	seller.Use(middleware.AuthenMiddleware())
	seller.Use(middleware.RequireRole("seller"))
	{
		sellerProducts := seller.Group("/products")
		{
			sellerProducts.GET("", productHandler.GetMyProducts)
			sellerProducts.PATCH("/:id/status", productHandler.UpdateProductStatus)
			sellerProducts.PATCH("/:id/stock", productHandler.UpdateProductStock)
			sellerProducts.PATCH("/bulk", productHandler.BulkUpdateProducts)
		}
	}

	// Category-specific routes
	categories := router.Group("/categories")
	{
		categories.GET("/:id/products", productHandler.ListProductsByCategory)
	}
}
