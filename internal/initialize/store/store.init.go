package store

import (
	"github.com/edynt/chogiare/veloras-api/internal/store/application/service"
	"github.com/edynt/chogiare/veloras-api/internal/store/controller/http"
	"github.com/edynt/chogiare/veloras-api/internal/store/infrastructure/persistence/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitStore(db *pgxpool.Pool) *http.StoreHandler {
	// Initialize repository
	storeRepo := repository.NewStoreRepository(db)

	// Initialize service
	storeService := service.NewStoreService(storeRepo)

	// Initialize handler
	storeHandler := http.NewStoreHandler(storeService)

	return storeHandler
}
