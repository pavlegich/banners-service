// Package handlers contains server controller object and
// methods for building the server route.
package handlers

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/pavlegich/banners-service/internal/controllers/middlewares"
	"github.com/pavlegich/banners-service/internal/domains/banner"
	banners "github.com/pavlegich/banners-service/internal/domains/banner/controllers/http"
	"github.com/pavlegich/banners-service/internal/infra/config"
)

// Controller contains database and configuration
// for building the server router.
type Controller struct {
	repo  banner.Repository
	cache banner.Cache
	cfg   *config.Config
}

// NewController creates and returns new server controller.
func NewController(ctx context.Context, repo banner.Repository, cache banner.Cache, cfg *config.Config) *Controller {
	return &Controller{
		repo:  repo,
		cache: cache,
		cfg:   cfg,
	}
}

// BuildRoute creates new router and appends handlers and middlewares to it.
func (c *Controller) BuildRoute(ctx context.Context) (*chi.Mux, error) {
	r := chi.NewRouter()

	r.Use(middlewares.WithLogging)
	r.Use(middlewares.Recovery)
	r.Use(middlewares.WithAuth)

	banners.Activate(ctx, r, c.cfg, c.repo, c.cache)

	return r, nil
}
