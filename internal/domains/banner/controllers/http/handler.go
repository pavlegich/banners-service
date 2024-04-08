// Package http contains banners object functions
// for activating the handler in controller, and handlers.
package http

import (
	"context"
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/pavlegich/banners-service/internal/domains/banner"
	repo "github.com/pavlegich/banners-service/internal/domains/banner/repository"
	"github.com/pavlegich/banners-service/internal/infra/config"
)

// BannersHandler contains objects for work with banner handlers.
type BannersHandler struct {
	Config  *config.Config
	Service banner.Service
}

// Activate activates handler for banner object.
func Activate(ctx context.Context, r *chi.Mux, cfg *config.Config, db *sql.DB) {
	s := banner.NewBannerService(ctx, repo.NewBannerRepository(ctx, db))
	newHandler(r, cfg, s)
}

// newHandler initializes handler for banner object.
func newHandler(r *chi.Mux, cfg *config.Config, s banner.Service) {
	_ = &BannersHandler{
		Config:  cfg,
		Service: s,
	}

	// r.HandleFunc("/api/banner", h.HandleBanner)
	// r.HandleFunc("/api/banners", h.HandleBanners)
}
