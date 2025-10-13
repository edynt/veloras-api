package middleware

import (
	authRepo "github.com/edynt/chogiare/veloras-api/internal/auth/domain/repository"
	authRepoImpl "github.com/edynt/chogiare/veloras-api/internal/auth/infrastructure/persistence/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuthRepoInstance holds the auth repository instance for middleware
var AuthRepoInstance authRepo.AuthRepository

// InitAuthorizationMiddleware initializes the authorization middleware with database connection
func InitAuthorizationMiddleware(db *pgxpool.Pool) {
	AuthRepoInstance = authRepoImpl.NewAuthRepository(db)
}

// RequirePermission creates a middleware that checks if the authenticated user has a specific permission
func RequirePermission(permissionName string) gin.HandlerFunc {
	if AuthRepoInstance == nil {
		panic("Authorization middleware not initialized. Call InitAuthorizationMiddleware first.")
	}
	return RequirePermissionWithRepo(AuthRepoInstance, permissionName)
}

// RequireRole creates a middleware that checks if the authenticated user has a specific role
func RequireRole(roleName string) gin.HandlerFunc {
	if AuthRepoInstance == nil {
		panic("Authorization middleware not initialized. Call InitAuthorizationMiddleware first.")
	}
	return RequireRoleWithRepo(AuthRepoInstance, roleName)
}
