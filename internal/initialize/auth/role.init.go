package initialize

import (
	"github.com/edynnt/veloras-api/internal/auth/application/service"
	"github.com/edynnt/veloras-api/internal/auth/controller/http"
	permissionRepo "github.com/edynnt/veloras-api/internal/auth/infrastructure/persistence/repository"
	roleRepo "github.com/edynnt/veloras-api/internal/auth/infrastructure/persistence/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRole(db *pgxpool.Pool) *http.RoleHandler {
	rRepo := roleRepo.NewRoleRepository(db)
	pRepo := permissionRepo.NewPermissionRepository(db)

	svc := service.NewRoleService(rRepo, pRepo)
	handler := http.NewRoleHandler(svc)

	return handler
}
