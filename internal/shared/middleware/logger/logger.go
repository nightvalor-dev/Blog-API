package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	file        *os.File
}

func NewLogger(logFilePath string) (*Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// write to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, file)

	return &Logger{
		infoLogger:  log.New(multiWriter, "", 0),
		errorLogger: log.New(multiWriter, "", 0),
		file:        file,
	}, nil
}

func (l *Logger) Close() {
	err := l.file.Close()
	if err != nil {
		return
	}
}
