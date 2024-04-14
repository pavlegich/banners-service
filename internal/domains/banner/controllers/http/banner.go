// Package http contains banners object functions
// for activating the handler in controller, and handlers.
package http

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/pavlegich/banners-service/internal/domains/banner"
	"github.com/pavlegich/banners-service/internal/infra/config"
)

// requestQuery contains data, which might be in request queries.
type requestQuery struct {
	tagID        int
	featureID    int
	lastRevision bool
	limit        int
	offset       int
}

// BannersHandler contains objects for work with banner handlers.
type BannerHandler struct {
	Config  *config.Config
	Service banner.Service
}

// Activate activates handler for banner object.
func Activate(ctx context.Context, r *chi.Mux, cfg *config.Config, repo banner.Repository, cache banner.Cache) {
	s := banner.NewBannerService(ctx, repo, cache)
	newHandler(r, cfg, s)
}

// newHandler initializes handler for banner object.
func newHandler(r *chi.Mux, cfg *config.Config, s banner.Service) {
	h := &BannerHandler{
		Config:  cfg,
		Service: s,
	}

	r.Get("/user_banner", h.HandleGetUserBanner)
	r.Get("/banner", h.HandleGetBanner)
	r.Post("/banner", h.HandleCreateBanner)
	r.Patch("/banner/{id}", h.HandleUpdateBanner)
	r.Delete("/banner/{id}", h.HandleDeleteBanner)
}
