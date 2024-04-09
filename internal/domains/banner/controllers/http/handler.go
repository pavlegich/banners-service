// Package http contains banners object functions
// for activating the handler in controller, and handlers.
package http

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pavlegich/banners-service/internal/domains/banner"
	repo "github.com/pavlegich/banners-service/internal/domains/banner/repository"
	errs "github.com/pavlegich/banners-service/internal/errors"
	"github.com/pavlegich/banners-service/internal/infra/config"
	"github.com/pavlegich/banners-service/internal/infra/logger"
	"github.com/pavlegich/banners-service/internal/utils"
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
	r.Patch("/banner/{id}", h.HandleUpdateBanner)
	r.Delete("/banner/{id}", h.HandleDeleteBanner)
}

// HandleCreateBanner handles request to create new banner.
func (h *BannerHandler) HandleCreateBanner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req banner.Banner
	var buf bytes.Buffer

	w.Header().Set("Content-Type", "application/json")
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.Error("HandleCreateBanner: read request body failed",
			zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		resp := utils.ParamToJSON("error", err.Error())
		w.Write(resp)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		logger.Log.Error("HandleCreateBanner: request unmarshal failed",
			zap.String("body", buf.String()),
			zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		resp := utils.ParamToJSON("error", err.Error())
		w.Write(resp)
		return
	}

	bannerID, err := h.Service.Create(ctx, &req)
	if err != nil {
		logger.Log.Error("HandleCreateBanner: create banner failed",
			zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		resp := utils.ParamToJSON("error", err.Error())
		w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp := utils.ParamToJSON("banner_id", strconv.Itoa(bannerID))
	w.Write(resp)
}

// HandleUpdateBanner handles request to update the requested banner.
func (h *BannerHandler) HandleUpdateBanner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	idString := chi.URLParam(r, "id")
	if idString == "" {
		logger.Log.Error("HandleUpdateBanner: id parameter is empty")

		w.WriteHeader(http.StatusBadRequest)
		resp := utils.ParamToJSON("error", "id parameter is empty")
		w.Write(resp)
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		logger.Log.Error("HandleUpdateBanner: convert id parameter to integer failed",
			zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		resp := utils.ParamToJSON("error", err.Error())
		w.Write(resp)
		return
	}

	var req banner.Banner
	var buf bytes.Buffer

	req.ID = id

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.Error("HandleUpdateBanner: read request body failed",
			zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		resp := utils.ParamToJSON("error", err.Error())
		w.Write(resp)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		logger.Log.Error("HandleUpdateBanner: request unmarshal failed",
			zap.String("body", buf.String()),
			zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		resp := utils.ParamToJSON("error", err.Error())
		w.Write(resp)
		return
	}

	err = h.Service.Update(ctx, &req)
	if err != nil {
		logger.Log.Error("HandleUpdateBanner: update data failed",
			zap.Error(err))

		if errors.Is(err, errs.ErrBannerNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		resp := utils.ParamToJSON("error", err.Error())
		w.Write(resp)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

// HandleDeleteBanner handles request to delete banner.
func (h *BannerHandler) HandleDeleteBanner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	idString := chi.URLParam(r, "id")
	if idString == "" {
		logger.Log.Error("HandleUpdateBanner: id parameter is empty")

		w.WriteHeader(http.StatusBadRequest)
		resp := utils.ParamToJSON("error", "id parameter is empty")
		w.Write(resp)
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		logger.Log.Error("HandleUpdateBanner: convert id parameter to integer failed",
			zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		resp := utils.ParamToJSON("error", err.Error())
		w.Write(resp)
		return
	}

	err = h.Service.Delete(ctx, id)
	if err != nil {
		logger.Log.Error("HandleDeleteBanner: delete data failed",
			zap.Error(err))

		if errors.Is(err, errs.ErrBannerNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		resp := utils.ParamToJSON("error", err.Error())
		w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
