package review

import (
	"github.com/edynt/chogiare/veloras-api/internal/review/application/service"
	"github.com/edynt/chogiare/veloras-api/internal/review/controller/http"
	"github.com/edynt/chogiare/veloras-api/internal/review/infrastructure/persistence/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitReview(db *pgxpool.Pool) *http.ReviewHandler {
	// Initialize repository
	reviewRepo := repository.NewReviewRepository(db)

	// Initialize service
	reviewService := service.NewReviewService(reviewRepo)

	// Initialize handler
	reviewHandler := http.NewReviewHandler(reviewService)

	return reviewHandler
}
