package http

import (
	"github.com/edynt/chogiare/veloras-api/internal/middleware"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterReviewRoutes(r *gin.RouterGroup, handler *ReviewHandler) {
	reviews := r.Group("/reviews")
	{
		// Public routes (if any)
		// reviews.GET("/", handler.ListReviews) // This might be public

		// Protected routes
		reviews.Use(middleware.AuthenMiddleware())
		{
			// Review CRUD operations
			reviews.POST("/", response.Wrap(handler.CreateReview))
			reviews.GET("/", response.Wrap(handler.ListReviews))
			reviews.GET("/my", response.Wrap(handler.ListUserReviews))
			reviews.GET("/product/:product_id", response.Wrap(handler.ListReviewsByProduct))
			reviews.GET("/:id", response.Wrap(handler.GetReview))
			reviews.PUT("/:id", response.Wrap(handler.UpdateReview))
			reviews.DELETE("/:id", response.Wrap(handler.DeleteReview))

			// Review interactions
			reviews.POST("/:id/helpful", response.Wrap(handler.MarkReviewHelpful))
			reviews.DELETE("/:id/helpful", response.Wrap(handler.UnmarkReviewHelpful))

			// Statistics
			reviews.GET("/stats/product/:product_id", response.Wrap(handler.GetReviewStats))
			reviews.GET("/stats/my", response.Wrap(handler.GetUserReviewStats))
		}
	}
}
