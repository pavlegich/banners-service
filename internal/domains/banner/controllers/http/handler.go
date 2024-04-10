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

// requestQuery contains data, which might be in request queries.
type requestQuery struct {
	tag_id        int
	feature_id    int
	last_revision bool
	limit         int
	offset        int
}

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

	r.Get("/banner", h.HandleGetBanner)
	r.Post("/banner", h.HandleCreateBanner)
	r.Patch("/banner/{id}", h.HandleUpdateBanner)
	r.Delete("/banner/{id}", h.HandleDeleteBanner)
}

// HandleGetBanner handles admin's request to get list of banners.
func (h *BannerHandler) HandleGetBanner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req requestQuery
	want := map[string]struct{}{
		"feature_id": {},
		"tag_id":     {},
		"limit":      {},
		"offset":     {},
	}

	w.Header().Set("Content-Type", "application/json")

	queries := r.URL.Query()
	for val := range queries {
		_, ok := want[val]
		if !ok {
			logger.Log.Error("HandleGetBanner: incorrect query",
				zap.String("query", val))

			w.WriteHeader(http.StatusBadRequest)
			resp := utils.ParamToJSON("error", "incorrect query in request url")
			w.Write(resp)
			return
		}

		if len(queries[val]) != 1 {
			logger.Log.Error("HandleGetBanner: incorrect queries number",
				zap.String("query_name", val),
				zap.Int("query_number", len(queries[val])))

			w.WriteHeader(http.StatusBadRequest)
			resp := utils.ParamToJSON("error", "incorrect query number in request url")
			w.Write(resp)
			return
		}

		current, err := strconv.Atoi(queries[val][0])
		if err != nil {
			logger.Log.Error("HandleGetBanner: convert query to integer failed",
				zap.String("query_name", val),
				zap.String("query_number", queries[val][0]))

			w.WriteHeader(http.StatusBadRequest)
			resp := utils.ParamToJSON("error", "convert query to integer failed")
			w.Write(resp)
			return
		}

		switch val {
		case "feature_id":
			req.feature_id = current
		case "tag_id":
			req.tag_id = current
		case "limit":
			req.limit = current
		case "offset":
			req.offset = current
		}
	}

	bannersList, err := h.Service.List(ctx, req.feature_id, req.tag_id, req.limit, req.offset)
	if err != nil {
		logger.Log.Error("HandleGetBanner: get banners list failed",
			zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		resp := utils.ParamToJSON("error", err.Error())
		w.Write(resp)
		return
	}

	bannersJSON, err := json.Marshal(bannersList)
	if err != nil {
		logger.Log.Error("HandleGetBanner: marshal banner data failed",
			zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		resp := utils.ParamToJSON("error", err.Error())
		w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bannersJSON)
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
