package middlewares

import (
	"net/http"

	"github.com/pavlegich/banners-service/internal/infra/logger"
	"go.uber.org/zap"
)

// Recovery recovers server operation when a server running panic occurs.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				logger.Log.Error("server panic",
					zap.Any("error", err),
				)

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
