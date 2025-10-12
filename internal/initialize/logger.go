package initialize

import (
	"github.com/edynt/chogiare/veloras-api/pkg/global"
	"github.com/edynt/chogiare/veloras-api/pkg/logger"
)

func InitLogger() {
	global.Logger = logger.NewLogger(global.Config.Logger)
}
