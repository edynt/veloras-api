package order

import (
	"github.com/edynt/chogiare/veloras-api/internal/order/application/service"
	"github.com/edynt/chogiare/veloras-api/internal/order/controller/http"
	"github.com/edynt/chogiare/veloras-api/internal/order/infrastructure/persistence/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitOrder(db *pgxpool.Pool) *http.OrderHandler {
	// Initialize repository
	orderRepo := repository.NewOrderRepository(db)

	// Initialize service
	orderService := service.NewOrderService(orderRepo)

	// Initialize handler
	orderHandler := http.NewOrderHandler(orderService)

	return orderHandler
}
