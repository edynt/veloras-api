package initialize

import (
	authHttp "github.com/edynnt/veloras-api/internal/auth/controller/http"
	permissionHttp "github.com/edynnt/veloras-api/internal/auth/controller/http"
	cronHttp "github.com/edynnt/veloras-api/internal/cron/controller/http"
	authInit "github.com/edynnt/veloras-api/internal/initialize/auth"
	userInit "github.com/edynnt/veloras-api/internal/initialize/user"
	"github.com/edynnt/veloras-api/internal/middleware"
	userHttp "github.com/edynnt/veloras-api/internal/user/controller/http"
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

	authHandler := authInit.InitAuth(db)
	authHttp.RegisterAuthRoutes(v1, authHandler)

	permissionHandler := authInit.InitPermission(db)
	permissionHttp.RegisterPermissionRoutes(v1, permissionHandler)

	roleHandler := authInit.InitRole(db)
	authHttp.RegisterRoleRoutes(v1, roleHandler)

	userHandler := userInit.InitUser(db)
	userHttp.RegisterUserRoutes(v1, userHandler)

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
