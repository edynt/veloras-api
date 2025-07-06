package initialize

import (
	"database/sql"

	"github.com/edynnt/veloras-api/internal/auth/application/service"
	"github.com/edynnt/veloras-api/internal/auth/controller/http"
	authRepo "github.com/edynnt/veloras-api/internal/auth/infrastructure/persistence/repository"
)

func InitAuth(db *sql.DB) *http.AuthHandler {
	authRepo := authRepo.NewAuthRepository(db)
	service := service.NewAuthService(authRepo)
	handler := http.NewAuthHandler(service)
	return handler
}
