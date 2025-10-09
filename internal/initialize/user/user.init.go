package initialize

import (
	"github.com/edynnt/veloras-api/internal/user/application/service"
	"github.com/edynnt/veloras-api/internal/user/controller/http"
	userRepo "github.com/edynnt/veloras-api/internal/user/infrastructure/persistence/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitUser(db *pgxpool.Pool) *http.UserHandler {
	userRepo := userRepo.NewUserRepository(db)
	service := service.NewUserService(userRepo)
	handler := http.NewUserHandler(service)
	return handler
}
