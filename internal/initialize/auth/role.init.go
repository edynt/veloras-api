package initialize

import (
	"github.com/edynnt/veloras-api/internal/auth/application/service"
	"github.com/edynnt/veloras-api/internal/auth/controller/http"
	roleRepo "github.com/edynnt/veloras-api/internal/auth/infrastructure/persistence/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRole(db *pgxpool.Pool) *http.RoleHandler {
	roleRepo := roleRepo.NewRoleRepository(db)
	service := service.NewRoleService(roleRepo)
	handler := http.NewRoleHandler(service)
	return handler
}
