// Package http contains banners object functions
// for activating the handler in controller, and handlers.
package http

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pavlegich/banners-service/internal/domains/banner"
	repo "github.com/pavlegich/banners-service/internal/domains/banner/repository"
	"github.com/pavlegich/banners-service/internal/infra/config"
	"github.com/pavlegich/banners-service/internal/infra/logger"
	"go.uber.org/zap"
)

// BannersHandler contains objects for work with banner handlers.
type BannerHandler struct {
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
	h := &BannerHandler{
		Config:  cfg,
		Service: s,
	}

	r.Post("/banner", h.HandleCreateBanner)
	// r.HandleFunc("/api/banners", h.HandleBanners)
}

// HandleCreateBanner handles request to create new banner.
func (h *BannerHandler) HandleCreateBanner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req banner.Banner
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.Error("HandleCreateBanner: read request body failed",
			zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		logger.Log.Error("HandleCreateBanner: request unmarshal failed",
			zap.String("body", buf.String()),
			zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bannerID, err := h.Service.Create(ctx, &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := map[string]int{
		"banner_id": bannerID,
	}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(respJSON))
}
