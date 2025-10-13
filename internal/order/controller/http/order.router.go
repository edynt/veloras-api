package http

import (
	"github.com/edynt/chogiare/veloras-api/internal/middleware"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(r *gin.RouterGroup, handler *OrderHandler) {
	orders := r.Group("/orders")

	// Public routes (no authentication required)
	// None for now

	// Protected routes (require authentication)
	ordersProtected := orders.Group("")
	ordersProtected.Use(middleware.AuthenMiddleware())
	{
		// User-specific order operations
		ordersProtected.POST("/", response.Wrap(handler.CreateOrder))
		ordersProtected.GET("/my", response.Wrap(handler.ListUserOrders))
		ordersProtected.GET("/:id", response.Wrap(handler.GetOrder))
		ordersProtected.PUT("/:id", response.Wrap(handler.UpdateOrder))
		ordersProtected.DELETE("/:id", response.Wrap(handler.DeleteOrder))

		// User statistics
		ordersProtected.GET("/stats/my", response.Wrap(handler.GetUserOrderStats))
	}

	// Admin routes (require admin role)
	ordersAdmin := orders.Group("")
	ordersAdmin.Use(middleware.AuthenMiddleware())
	ordersAdmin.Use(middleware.RequireRole("admin"))
	{
		// Admin order operations
		ordersAdmin.GET("/", response.Wrap(handler.ListOrders))
		ordersAdmin.GET("/store/:store_id", response.Wrap(handler.ListStoreOrders))
		ordersAdmin.PATCH("/:id/status", response.Wrap(handler.UpdateOrderStatus))
		ordersAdmin.PATCH("/:id/payment-status", response.Wrap(handler.UpdateOrderPaymentStatus))

		// Admin statistics
		ordersAdmin.GET("/stats", response.Wrap(handler.GetOrderStats))
		ordersAdmin.GET("/stats/store/:store_id", response.Wrap(handler.GetStoreOrderStats))
	}
}
