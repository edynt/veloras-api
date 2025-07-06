package initialize

import (
	"database/sql"

	"github.com/edynnt/veloras-api/internal/auth/controller/http"
	initialize "github.com/edynnt/veloras-api/internal/initialize/auth"
	"github.com/gin-gonic/gin"
)

func InitRouter(db *sql.DB, logLevel string) *gin.Engine {
	var r *gin.Engine

	if logLevel == "debug" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	v1 := r.Group("/api/v1")

	authHandler := initialize.InitAuth(db)
	http.RegisterAuthRoutes(v1, authHandler)

	// userHandler := initialize.InitUser(db)
	// http.RegisterUserRoutes(v1, userHandler)

	return r

}
