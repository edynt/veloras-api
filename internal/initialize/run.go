package initialize

import (
	"log"

	cronHttp "github.com/edynnt/veloras-api/internal/cron/controller/http"
	"github.com/edynnt/veloras-api/pkg/global"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/gin-gonic/gin"
)

func Run() (*gin.Engine, string) {

	global.Config = MustLoadConfig()

	InitLogger()

	db, err := InitDB(&global.Config)
	if err != nil {
		log.Fatalf("%s: %v", msg.FailedToInitDB, err)
	}

	// Initialize and start cron scheduler
	cronScheduler := InitCronScheduler(db)
	if err := cronScheduler.Start(); err != nil {
		log.Fatalf("Failed to start cron scheduler: %v", err)
	}

	// Create cron handler for HTTP endpoints
	cronHandler := cronHttp.NewCronHandler(cronScheduler)

	// Create log handler for log viewing
	logHandler := cronHttp.NewLogHandler("./storage/logs")

	r := InitRouter(db, global.Config.Logger.Log_level, cronHandler, logHandler)
	return r, global.Config.Server.Port
}
