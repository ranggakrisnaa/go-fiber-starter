package logs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

type LogController struct{}

type LogEntry struct {
	Type    string
	Content string
	Time    time.Time
}

func NewLogController() *LogController {
	return &LogController{}
}

func (lc *LogController) GetLogs(c fiber.Ctx) error {
	month := c.Params("month", "")
	if month == "" {
		month = time.Now().Format("January")
		month = strings.ToLower(month)
	}

	logTypes := []string{"query", "http", "app"}
	var allLogs []LogEntry

	for _, logType := range logTypes {
		logs := lc.readLogFile(logType, month)
		allLogs = append(allLogs, logs...)
	}

	sort.Slice(allLogs, func(i, j int) bool {
		return allLogs[i].Time.After(allLogs[j].Time)
	})
	fmt.Println(allLogs)

	return c.Render("./public/logs.html", fiber.Map{
		"Month": month,
		"Logs":  allLogs,
	})
}

func (lc *LogController) readLogFile(logType, month string) []LogEntry {
	var logDir string
	var logFileName string

	switch logType {
	case "query":
		logDir = "./config/logs/query_log"
		logFileName = fmt.Sprintf("%s_query.log", month)
	case "http":
		logDir = "./config/logs/http"
		logFileName = fmt.Sprintf("%s_http.log", month)
	case "app":
		logDir = "./config/logs/app"
		logFileName = fmt.Sprintf("%s_app.log", month)
	default:
		return []LogEntry{}
	}

	logFilePath := filepath.Join(logDir, logFileName)

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		return []LogEntry{}
	}

	file, err := os.Open(logFilePath)
	if err != nil {
		return []LogEntry{}
	}
	defer file.Close()

	var logs []LogEntry
	scanner := bufio.NewScanner(file)
	var currentLog strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if lc.isTimestampLine(line) {
			if currentLog.Len() > 0 {
				logContent := strings.TrimSpace(currentLog.String())
				if logContent != "" {
					logTime := lc.parseTimeFromLog(logContent)
					logs = append(logs, LogEntry{
						Type:    logType,
						Content: logContent,
						Time:    logTime,
					})
				}
			}
			currentLog.Reset()
			currentLog.WriteString(line)
		} else if line != "" {
			if currentLog.Len() > 0 {
				currentLog.WriteString("\n")
				currentLog.WriteString(line)
			}
		}
	}

	if currentLog.Len() > 0 {
		logContent := strings.TrimSpace(currentLog.String())
		if logContent != "" {
			logTime := lc.parseTimeFromLog(logContent)
			logs = append(logs, LogEntry{
				Type:    logType,
				Content: logContent,
				Time:    logTime,
			})
		}
	}

	return logs
}

func (lc *LogController) isTimestampLine(line string) bool {
	if len(line) < 19 {
		return false
	}
	var timeStr string
	if strings.HasPrefix(line, "[") {
		if len(line) < 21 {
			return false
		}
		timeStr = line[1:20]
	} else {
		timeStr = line[:19]
	}
	_, err := time.Parse("2006/01/02 15:04:05", timeStr)
	return err == nil
}

func (lc *LogController) parseTimeFromLog(line string) time.Time {
	var timeStr string
	if strings.HasPrefix(line, "[") {
		timeStr = line[1:20]
	} else if len(line) >= 19 {
		timeStr = line[:19]
	} else {
		return time.Now()
	}

	t, err := time.Parse("2006/01/02 15:04:05", timeStr)
	if err != nil {
		return time.Now()
	}
	return t
}
