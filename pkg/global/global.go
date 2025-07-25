package global

import (
	"github.com/edynnt/veloras-api/pkg/config"
	"github.com/edynnt/veloras-api/pkg/logger"
)

var (
	Config config.Config
	Logger *logger.LoggerZap
)
