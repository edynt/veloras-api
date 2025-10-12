package global

import (
	"github.com/edynt/chogiare/veloras-api/pkg/config"
	"github.com/edynt/chogiare/veloras-api/pkg/logger"
)

var (
	Config config.Config
	Logger *logger.LoggerZap
)
