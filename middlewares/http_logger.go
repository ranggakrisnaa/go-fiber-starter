package middlewares

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

var (
	httpLogFile *os.File
)

func initHTTPLogger() {
	err := os.MkdirAll("./config/logs/http", os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create http log directory: %v", err)
	}

	currentMonth := time.Now().Format("January")
	currentMonth = strings.ToLower(currentMonth)
	logFileName := currentMonth + "_http.log"

	httpLogFile, err = os.OpenFile("./config/logs/http/"+logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open http log file: %v", err)
	}
}

func HTTPLogger() fiber.Handler {
	if httpLogFile == nil {
		initHTTPLogger()
	}

	return func(ctx fiber.Ctx) error {
		start := time.Now()

		err := ctx.Next()

		duration := time.Since(start)
		status := ctx.Response().StatusCode()
		method := ctx.Method()
		path := ctx.Path()
		ip := ctx.IP()
		userAgent := ctx.Get("User-Agent")

		logEntry := fmt.Sprintf(
			"[%s] %s %s %s %d %v %s",
			time.Now().Format("2006/01/02 15:04:05"),
			ip,
			method,
			path,
			status,
			duration,
			userAgent,
		)

		if _, err := httpLogFile.WriteString(logEntry + "\n"); err != nil {
			log.Printf("failed to write http log: %v", err)
		}

		// Also print to terminal
		log.Printf("[HTTP] %s", strings.TrimSpace(logEntry))

		return err
	}
}