package middleware

import (
	"net/http"
	"time"

	"github.com/SysoevAndrey/shortener/internal/logger"
	"go.uber.org/zap"
)

type responseData struct {
	statusCode int
	size       int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (lw *loggingResponseWriter) WriteHeader(statusCode int) {
	lw.ResponseWriter.WriteHeader(statusCode)
	lw.responseData.statusCode = statusCode
}

func (lw *loggingResponseWriter) Write(p []byte) (int, error) {
	size, err := lw.ResponseWriter.Write(p)
	lw.responseData.size += size
	return size, err
}

func WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rd := &responseData{
			statusCode: 0,
			size:       0,
		}

		lw := loggingResponseWriter{ResponseWriter: w, responseData: rd}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Info(
			"Request processed",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			zap.Int("status", lw.responseData.statusCode),
			zap.Int("size", lw.responseData.size),
		)
	})
}
