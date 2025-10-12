package category

import (
	"github.com/edynt/chogiare/veloras-api/internal/category/application/service"
	"github.com/edynt/chogiare/veloras-api/internal/category/controller/http"
	"github.com/edynt/chogiare/veloras-api/internal/category/infrastructure/persistence/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitCategory(db *pgxpool.Pool) *http.CategoryHandler {
	// Initialize repository
	categoryRepo := repository.NewCategoryRepository(db)

	// Initialize service
	categoryService := service.NewCategoryService(categoryRepo)

	// Initialize handler
	categoryHandler := http.NewCategoryHandler(categoryService)

	return categoryHandler
}
