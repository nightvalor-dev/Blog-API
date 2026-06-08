package logger

import (
	"encoding/json"
	"time"
)

type LogEntry struct {
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
	Method    string `json:"method,omitempty"`
	Path      string `json:"path,omitempty"`
	Status    int    `json:"status,omitempty"`
	Duration  string `json:"duration,omitempty"`
	IP        string `json:"ip,omitempty"`
	Message   string `json:"message,omitempty"`
}

func (logger *Logger) write(level string, entry LogEntry) {
	entry.Level = level
	entry.Timestamp = time.Now().Format(time.RFC3339)

	b, err := json.Marshal(entry)
	if err != nil {
		return
	}

	if level == "ERROR" {
		logger.errorLogger.Println(string(b))
	} else {
		logger.infoLogger.Println(string(b))
	}
}

func (logger *Logger) Info(entry LogEntry) {
	logger.write("INFO", entry)
}

func (logger *Logger) Error(entry LogEntry) {
	logger.write("ERROR", entry)
}
