package initialize

import (
	authHttp "github.com/edynnt/veloras-api/internal/auth/controller/http"
	permissionHttp "github.com/edynnt/veloras-api/internal/auth/controller/http"
	cronHttp "github.com/edynnt/veloras-api/internal/cron/controller/http"
	initialize "github.com/edynnt/veloras-api/internal/initialize/auth"
	"github.com/edynnt/veloras-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(db *pgxpool.Pool, logLevel string, cronScheduler *cronHttp.CronHandler, logHandler *cronHttp.LogHandler) *gin.Engine {
	var r *gin.Engine

	if logLevel == "debug" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	r.Use(middleware.CORS) // cross
	r.Use(middleware.ValidatorMiddleware())

	v1 := r.Group("/api/v1")

	authHandler := initialize.InitAuth(db)
	authHttp.RegisterAuthRoutes(v1, authHandler)

	permissionHandler := initialize.InitPermission(db)
	permissionHttp.RegisterPermissionRoutes(v1, permissionHandler)

	roleHandler := initialize.InitRole(db)
	authHttp.RegisterRoleRoutes(v1, roleHandler)

	// Admin routes for cron management
	if cronScheduler != nil {
		cronHttp.RegisterCronRoutes(v1, cronScheduler)
	}

	// Admin routes for log management
	if logHandler != nil {
		cronHttp.RegisterLogRoutes(v1, logHandler)
	}

	// Serve static files (log viewer)
	r.Static("/web", "./web")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/web/logs.html")
	})

	return r
}
