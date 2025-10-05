package logger

import (
	"gorm.io/gorm/logger"
)

const (
	LOG_DIR = "./config/logs/query_log"
)

func SetupLogger() logger.Interface {
	return NewQueryLogger()
}