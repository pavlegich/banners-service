package middlewares

import (
	"context"
	"net/http"

	"github.com/pavlegich/banners-service/internal/infra/logger"
	"github.com/pavlegich/banners-service/internal/utils"
	"go.uber.org/zap"
)

// WithAuth checks and validates authorization token.
func WithAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")

		var role string
		switch token {
		case "admin_token":
			role = "admin"
		case "user_token":
			role = "user"
			if r.RequestURI != "/user_banner" {
				logger.Log.Error("WithAuth: no permissions to access resource",
					zap.String("role", role),
					zap.String("uri", r.RequestURI))
				w.WriteHeader(http.StatusForbidden)
			}
		default:
			logger.Log.Error("WithAuth: unknown user role",
				zap.String("role", token))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), utils.ContextRoleKey, role)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
