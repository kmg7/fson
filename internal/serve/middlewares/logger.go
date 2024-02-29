package middlewares

import (
	"net/http"

	"github.com/kmg7/fson/internal/logger"
)

type Logger struct {
	log logger.AppLogger
}

func NewLogger(l logger.AppLogger) *Logger {
	return &Logger{
		log: l,
	}
}

func (m *Logger) Next(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.log.Info(r.Method, "\t", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
