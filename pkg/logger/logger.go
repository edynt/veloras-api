package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/edynnt/veloras-api/pkg/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerZap struct {
	*zap.Logger
}

func NewLogger(config config.LoggerSetting) *LoggerZap {
	logLevel := config.Log_level
	// debug -> info -> warning -> error -> fatal -> panic
	var level zapcore.Level

	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warning":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	case "panic":
		level = zap.PanicLevel
	default:
		level = zap.InfoLevel
	}

	encoder := getEncoderLog()

	// Generate daily log filename
	dailyLogFile := generateDailyLogFile(config.File_log_name)

	hook := lumberjack.Logger{
		Filename:   dailyLogFile,
		MaxSize:    config.Max_size, // megabytes
		MaxBackups: config.Max_backups,
		MaxAge:     config.Max_age,  //days
		Compress:   config.Compress, // disabled by default
	}

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), level)
	return &LoggerZap{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
}

// generateDailyLogFile creates a log filename with date suffix
func generateDailyLogFile(basePath string) string {
	dir := filepath.Dir(basePath)
	filename := filepath.Base(basePath)
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]

	// Add date suffix: YYYY-MM-DD
	dateStr := time.Now().Format("2006-01-02")

	return filepath.Join(dir, name+"-"+dateStr+ext)
}

func getEncoderLog() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()

	// timestamp => dd/mm/yyyy
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// ts -> time
	encodeConfig.TimeKey = "time"

	// level
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// caller
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewJSONEncoder(encodeConfig)
}
