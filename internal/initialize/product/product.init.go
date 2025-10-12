package product

import (
	"github.com/edynt/chogiare/veloras-api/internal/product/application/service"
	"github.com/edynt/chogiare/veloras-api/internal/product/controller/http"
	"github.com/edynt/chogiare/veloras-api/internal/product/infrastructure/persistence/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitProduct(db *pgxpool.Pool) *http.ProductHandler {
	// Initialize repository
	productRepo := repository.NewProductRepository(db)

	// Initialize service
	productService := service.NewProductService(productRepo)

	// Initialize handler
	productHandler := http.NewProductHandler(productService)

	return productHandler
}
