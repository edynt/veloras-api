package http

import (
	"github.com/edynt/chogiare/veloras-api/internal/middleware"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(r *gin.RouterGroup, handler *OrderHandler) {
	orders := r.Group("/orders")
	{
		// Public routes (if any)
		// orders.GET("/", handler.ListOrders) // This might be admin-only

		// Protected routes
		orders.Use(middleware.AuthenMiddleware())
		{
			// Order CRUD operations
			orders.POST("/", response.Wrap(handler.CreateOrder))
			orders.GET("/", response.Wrap(handler.ListOrders))
			orders.GET("/my", response.Wrap(handler.ListUserOrders))
			orders.GET("/store/:store_id", response.Wrap(handler.ListStoreOrders))
			orders.GET("/:id", response.Wrap(handler.GetOrder))
			orders.PUT("/:id", response.Wrap(handler.UpdateOrder))
			orders.PATCH("/:id/status", response.Wrap(handler.UpdateOrderStatus))
			orders.PATCH("/:id/payment-status", response.Wrap(handler.UpdateOrderPaymentStatus))
			orders.DELETE("/:id", response.Wrap(handler.DeleteOrder))

			// Statistics
			orders.GET("/stats", response.Wrap(handler.GetOrderStats))
			orders.GET("/stats/my", response.Wrap(handler.GetUserOrderStats))
			orders.GET("/stats/store/:store_id", response.Wrap(handler.GetStoreOrderStats))
		}
	}
}
