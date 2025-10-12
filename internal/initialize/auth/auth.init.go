package initialize

import (
	"github.com/edynt/chogiare/veloras-api/internal/auth/application/service"
	"github.com/edynt/chogiare/veloras-api/internal/auth/controller/http"
	authRepo "github.com/edynt/chogiare/veloras-api/internal/auth/infrastructure/persistence/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitAuth(db *pgxpool.Pool) *http.AuthHandler {
	authRepo := authRepo.NewAuthRepository(db)
	service := service.NewAuthService(authRepo)
	handler := http.NewAuthHandler(service)
	return handler
}
