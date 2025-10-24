package logger

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type QueryLogger struct {
	gormLogger.Interface
	file *os.File
}

func NewQueryLogger() gormLogger.Interface {
	err := os.MkdirAll("./config/logs/query_log", os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}

	currentMonth := time.Now().Format("January")
	currentMonthLower := strings.ToLower(currentMonth)
	logFileName := currentMonthLower + "_query.log"

	file, err := os.OpenFile("./config/logs/query_log/"+logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	return &QueryLogger{
		Interface: gormLogger.Default.LogMode(gormLogger.Info),
		file:      file,
	}
}

func (l *QueryLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.Interface = l.Interface.LogMode(level)
	return &newLogger
}

func (l *QueryLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logEntry := fmt.Sprintf("%s %s\n", time.Now().Format("2006/01/02 15:04:05"), fmt.Sprintf(msg, data...))
	l.file.WriteString(logEntry)
	l.Interface.Info(ctx, msg, data...)
}

func (l *QueryLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logEntry := fmt.Sprintf("%s %s\n", time.Now().Format("2006/01/02 15:04:05"), fmt.Sprintf(msg, data...))
	l.file.WriteString(logEntry)
	l.Interface.Warn(ctx, msg, data...)
}

func (l *QueryLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logEntry := fmt.Sprintf("%s %s\n", time.Now().Format("2006/01/02 15:04:05"), fmt.Sprintf(msg, data...))
	l.file.WriteString(logEntry)
	l.Interface.Error(ctx, msg, data...)
}

func (l *QueryLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	var logEntry string
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logEntry = fmt.Sprintf("%s %s\n%s\n%.3fms %s\n\n",
			time.Now().Format("2006/01/02 15:04:05"),
			fmt.Sprintf("[ERROR] %v", err),
			sql,
			float64(elapsed.Nanoseconds())/1e6,
			fmt.Sprintf("[rows:%v]", rows))
	} else if elapsed > time.Second {
		logEntry = fmt.Sprintf("%s %s\n%s\n%.3fms %s\n\n",
			time.Now().Format("2006/01/02 15:04:05"),
			"[SLOW SQL]",
			sql,
			float64(elapsed.Nanoseconds())/1e6,
			fmt.Sprintf("[rows:%v]", rows))
	} else {
		logEntry = fmt.Sprintf("%s %s\n%s\n%.3fms %s\n\n",
			time.Now().Format("2006/01/02 15:04:05"),
			"[QUERY]",
			sql,
			float64(elapsed.Nanoseconds())/1e6,
			fmt.Sprintf("[rows:%v]", rows))
	}

	if _, err := l.file.WriteString(logEntry); err != nil {
		log.Printf("failed to write query log: %v", err)
	}

	// Also print to terminal with clean formatting
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[QUERY ERROR] %v", err)
		log.Printf("[SQL] %s", sql)
		log.Printf("[TIME] %.3fms [rows:%v]", float64(elapsed.Nanoseconds())/1e6, rows)
	} else if elapsed > time.Second {
		log.Printf("[SLOW QUERY] %.3fms [rows:%v]", float64(elapsed.Nanoseconds())/1e6, rows)
		log.Printf("[SQL] %s", sql)
	} else {
		log.Printf("[QUERY] %.3fms [rows:%v]", float64(elapsed.Nanoseconds())/1e6, rows)
		log.Printf("[SQL] %s", sql)
	}

	l.Interface.Trace(ctx, begin, fc, err)
}
