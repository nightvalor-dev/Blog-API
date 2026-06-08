package logger

import (
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
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
