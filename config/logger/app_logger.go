package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	appLogFile *os.File
)

func InitAppLogger() {
	err := os.MkdirAll("./config/logs/app", os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create app log directory: %v", err)
	}

	currentMonth := time.Now().Format("January")
	currentMonth = strings.ToLower(currentMonth)
	logFileName := currentMonth + "_app.log"

	appLogFile, err = os.OpenFile("./config/logs/app/"+logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open app log file: %v", err)
	}
}

func LogInfo(message string) {
	if appLogFile == nil {
		InitAppLogger()
	}
	logEntry := fmt.Sprintf("[%s] [INFO] %s\n", time.Now().Format("2006/01/02 15:04:05"), message)
	if _, err := appLogFile.WriteString(logEntry); err != nil {
		log.Printf("failed to write app log: %v", err)
	}
	// Also print to terminal
	log.Printf("[APP] [INFO] %s", message)
}

func LogError(message string) {
	if appLogFile == nil {
		InitAppLogger()
	}
	logEntry := fmt.Sprintf("[%s] [ERROR] %s\n", time.Now().Format("2006/01/02 15:04:05"), message)
	if _, err := appLogFile.WriteString(logEntry); err != nil {
		log.Printf("failed to write app log: %v", err)
	}
	// Also print to terminal
	log.Printf("[APP] [ERROR] %s", message)
}

func LogAuth(message string) {
	if appLogFile == nil {
		InitAppLogger()
	}
	logEntry := fmt.Sprintf("[%s] [AUTH] %s\n", time.Now().Format("2006/01/02 15:04:05"), message)
	if _, err := appLogFile.WriteString(logEntry); err != nil {
		log.Printf("failed to write app log: %v", err)
	}
	// Also print to terminal
	log.Printf("[APP] [AUTH] %s", message)
}
