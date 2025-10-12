package initialize

import (
	authHttp "github.com/edynt/chogiare/veloras-api/internal/auth/controller/http"
	permissionHttp "github.com/edynt/chogiare/veloras-api/internal/auth/controller/http"
	cartHttp "github.com/edynt/chogiare/veloras-api/internal/cart/controller/http"
	categoryHttp "github.com/edynt/chogiare/veloras-api/internal/category/controller/http"
	chatHttp "github.com/edynt/chogiare/veloras-api/internal/chat/controller/http"
	cronHttp "github.com/edynt/chogiare/veloras-api/internal/cron/controller/http"
	authInit "github.com/edynt/chogiare/veloras-api/internal/initialize/auth"
	cartInit "github.com/edynt/chogiare/veloras-api/internal/initialize/cart"
	categoryInit "github.com/edynt/chogiare/veloras-api/internal/initialize/category"
	chatInit "github.com/edynt/chogiare/veloras-api/internal/initialize/chat"
	orderInit "github.com/edynt/chogiare/veloras-api/internal/initialize/order"
	productInit "github.com/edynt/chogiare/veloras-api/internal/initialize/product"
	reviewInit "github.com/edynt/chogiare/veloras-api/internal/initialize/review"
	storeInit "github.com/edynt/chogiare/veloras-api/internal/initialize/store"
	userInit "github.com/edynt/chogiare/veloras-api/internal/initialize/user"
	"github.com/edynt/chogiare/veloras-api/internal/middleware"
	orderHttp "github.com/edynt/chogiare/veloras-api/internal/order/controller/http"
	productHttp "github.com/edynt/chogiare/veloras-api/internal/product/controller/http"
	reviewHttp "github.com/edynt/chogiare/veloras-api/internal/review/controller/http"
	storeHttp "github.com/edynt/chogiare/veloras-api/internal/store/controller/http"
	userHttp "github.com/edynt/chogiare/veloras-api/internal/user/controller/http"
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

	// E-commerce routes
	productHandler := productInit.InitProduct(db)
	productHttp.RegisterProductRoutes(v1, productHandler)

	categoryHandler := categoryInit.InitCategory(db)
	categoryHttp.RegisterCategoryRoutes(v1, categoryHandler)

	cartHandler := cartInit.InitCart(db)
	cartHttp.RegisterCartRoutes(v1, cartHandler)

	reviewHandler := reviewInit.InitReview(db)
	reviewHttp.RegisterReviewRoutes(v1, reviewHandler)

	storeHandler := storeInit.InitStore(db)
	storeHttp.RegisterStoreRoutes(v1, storeHandler)

	chatHandler := chatInit.InitChat(db)
	chatHttp.RegisterChatRoutes(v1, chatHandler)

	orderHandler := orderInit.InitOrder(db)
	orderHttp.RegisterOrderRoutes(v1, orderHandler)

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
