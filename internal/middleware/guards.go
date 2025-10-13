package middleware

import (
	"net/http"

	"github.com/edynt/chogiare/veloras-api/internal/auth/domain/repository"
	"github.com/edynt/chogiare/veloras-api/pkg/response/msg"
	"github.com/edynt/chogiare/veloras-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

// RequirePermissionWithRepo creates a middleware that checks if the authenticated user has a specific permission
func RequirePermissionWithRepo(authRepo repository.AuthRepository, permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthenMiddleware)
		subject := c.Value("subjectID")
		if subject == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized, "message": msg.Unauthorized, "error_details": msg.HeaderAuthenticationNotFound, "error": true,
			})
			return
		}

		userID := utils.StringToInt(subject.(string))
		if userID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized, "message": msg.Unauthorized, "error_details": msg.UserIdInvalid, "error": true,
			})
			return
		}

		// Check if user has the required permission
		hasPermission, err := authRepo.UserHasPermission(c.Request.Context(), userID, permissionName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError, "message": "Internal server error", "error_details": "Failed to check user permissions", "error": true,
			})
			return
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden, "message": "Forbidden", "error_details": "Insufficient permissions", "error": true,
			})
			return
		}

		c.Next()
	}
}

// RequireRoleWithRepo creates a middleware that checks if the authenticated user has a specific role
func RequireRoleWithRepo(authRepo repository.AuthRepository, roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthenMiddleware)
		subject := c.Value("subjectID")
		if subject == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized, "message": msg.Unauthorized, "error_details": msg.HeaderAuthenticationNotFound, "error": true,
			})
			return
		}

		userID := utils.StringToInt(subject.(string))
		if userID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized, "message": msg.Unauthorized, "error_details": msg.UserIdInvalid, "error": true,
			})
			return
		}

		// Check if user has the required role
		hasRole, err := authRepo.UserHasRole(c.Request.Context(), userID, roleName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError, "message": "Internal server error", "error_details": "Failed to check user roles", "error": true,
			})
			return
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden, "message": "Forbidden", "error_details": "Insufficient role privileges", "error": true,
			})
			return
		}

		c.Next()
	}
}
