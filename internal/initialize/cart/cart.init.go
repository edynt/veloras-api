package cart

import (
	"github.com/edynt/chogiare/veloras-api/internal/cart/application/service"
	"github.com/edynt/chogiare/veloras-api/internal/cart/controller/http"
	"github.com/edynt/chogiare/veloras-api/internal/cart/infrastructure/persistence/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitCart(db *pgxpool.Pool) *http.CartHandler {
	// Initialize repository
	cartRepo := repository.NewCartRepository(db)

	// Initialize service
	cartService := service.NewCartService(cartRepo)

	// Initialize handler
	cartHandler := http.NewCartHandler(cartService)

	return cartHandler
}
