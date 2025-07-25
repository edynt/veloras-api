package initialize

import (
	"log"

	"github.com/edynnt/veloras-api/pkg/global"
	"github.com/gin-gonic/gin"
)

func Run() (*gin.Engine, string) {

	global.Config = MustLoadConfig()

	InitLogger()

	db, err := InitDB(&global.Config)
	if err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}

	r := InitRouter(db, global.Config.Logger.Log_level)
	return r, global.Config.Server.Port
}
