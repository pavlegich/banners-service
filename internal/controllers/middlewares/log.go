// Package middlewares contains function for logging actions from the handlers.
package middlewares

import (
	"bytes"
	"net/http"
	"time"

	"github.com/pavlegich/banners-service/internal/infra/logger"
	"go.uber.org/zap"
)

// WithLogging logs actions from the handlers.
func WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &logger.ResponseData{
			Status: 0,
			Size:   0,
			Body:   bytes.NewBufferString(""),
		}
		lw := logger.LoggingResponseWriter{
			ResponseWriter: w,
			ResponseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)
		// role, _ := utils.GetUserRoleFromContext(r.Context())

		logger.Log.Info("incoming HTTP request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			// zap.String("user_role", role),
			zap.Int("status", responseData.Status),
			zap.Int("size", responseData.Size),
			// zap.String("body", responseData.Body.String()),
		)
	})
}
