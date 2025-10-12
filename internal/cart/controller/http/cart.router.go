package http

import (
	"github.com/edynt/chogiare/veloras-api/internal/middleware"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterCartRoutes(r *gin.RouterGroup, handler *CartHandler) {
	cart := r.Group("/cart")
	{
		// Protected routes
		cart.Use(middleware.AuthenMiddleware())
		{
			// Cart operations
			cart.GET("/", response.Wrap(handler.GetCart))
			cart.POST("/clear", response.Wrap(handler.ClearCart))
			cart.GET("/stats", response.Wrap(handler.GetCartStats))

			// Cart item operations
			cart.POST("/items", response.Wrap(handler.AddCartItem))
			cart.PUT("/items/:id", response.Wrap(handler.UpdateCartItemQuantity))
			cart.DELETE("/items/:id", response.Wrap(handler.RemoveCartItem))
		}
	}
}
