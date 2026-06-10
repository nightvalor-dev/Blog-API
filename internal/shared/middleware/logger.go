package middleware

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
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

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	file        *os.File
}

// responseWriter wraps standard http.ResponseWriter to capture status codes
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool // Tracks if header was explicitly written
}

func NewLogger(logFilePath string) (*Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Stream to both stdout and your persistent file
	multiWriter := io.MultiWriter(os.Stdout, file)

	return &Logger{
		infoLogger:  log.New(multiWriter, "", 0),
		errorLogger: log.New(multiWriter, "", 0),
		file:        file,
	}, nil
}

func (logger *Logger) write(level string, entry LogEntry) {
	entry.Level = level
	entry.Timestamp = time.Now().Format(time.RFC3339)

	b, err := json.Marshal(entry)
	if err != nil {
		// Fallback if JSON encoding fails entirely
		log.Printf("failed to marshal log entry: %v", err)
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

func (logger *Logger) Close() {
	if err := logger.file.Close(); err != nil {
		// Log to standard error if closing the file fails
		log.SetOutput(os.Stderr)
		log.Printf("error closing log file: %v", err)
	}
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		// If Write is called before WriteHeader, status is implicitly 200 OK
		rw.wroteHeader = true
	}
	return rw.ResponseWriter.Write(b)
}

func LoggerMiddleware(logger *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := newResponseWriter(w)

			next.ServeHTTP(ww, r)

			logger.Info(LogEntry{
				Method:   r.Method,
				Path:     r.URL.Path,
				Status:   ww.status,
				Duration: time.Since(start).String(),
				IP:       r.RemoteAddr,
			})
		})
	}
}
