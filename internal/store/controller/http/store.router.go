package http

import (
	"github.com/edynt/chogiare/veloras-api/internal/middleware"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterStoreRoutes(r *gin.RouterGroup, handler *StoreHandler) {
	stores := r.Group("/stores")
	{
		// Public routes (if any)
		// stores.GET("/", handler.ListStores) // This might be public

		// Protected routes
		stores.Use(middleware.AuthenMiddleware())
		{
			// Store CRUD operations
			stores.POST("/", response.Wrap(handler.CreateStore))
			stores.GET("/", response.Wrap(handler.ListStores))
			stores.GET("/my", response.Wrap(handler.GetMyStore))
			stores.GET("/search", response.Wrap(handler.SearchStores))
			stores.GET("/:id", response.Wrap(handler.GetStore))
			stores.PUT("/:id", response.Wrap(handler.UpdateStore))
			stores.DELETE("/:id", response.Wrap(handler.DeleteStore))

			// Statistics
			stores.GET("/stats", response.Wrap(handler.GetStoreStats))
		}
	}
}
