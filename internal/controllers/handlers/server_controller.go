// Package handlers contains server controller object and
// methods for building the server route.
package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pavlegich/banners-service/internal/controllers/middlewares"
	banners "github.com/pavlegich/banners-service/internal/domains/banner/controllers/http"
	"github.com/pavlegich/banners-service/internal/infra/config"
)

// Controller contains database and configuration
// for building the server router.
type Controller struct {
	db  *sql.DB
	cfg *config.Config
}

// NewController creates and returns new server controller.
func NewController(ctx context.Context, db *sql.DB, cfg *config.Config) *Controller {
	return &Controller{
		db:  db,
		cfg: cfg,
	}
}

// BuildRoute creates new router and appends handlers and middlewares to it.
func (c *Controller) BuildRoute(ctx context.Context) (*chi.Mux, error) {
	r := chi.NewRouter()

	r.Use(middlewares.WithLogging)
	r.Use(middlewares.Recovery)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})

	banners.Activate(ctx, r, c.cfg, c.db)

	return r, nil
}
