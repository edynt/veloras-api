package initialize

import (
	"github.com/edynnt/veloras-api/internal/permission/application/service"
	"github.com/edynnt/veloras-api/internal/permission/controller/http"
	permissionRepo "github.com/edynnt/veloras-api/internal/permission/infrastructure/persistence/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPermission(db *pgxpool.Pool) *http.PermissionHandler {
	permissionRepo := permissionRepo.NewPermissionRepository(db)
	service := service.NewPermissionService(permissionRepo)
	handler := http.NewPermissionHandler(service)
	return handler
}
