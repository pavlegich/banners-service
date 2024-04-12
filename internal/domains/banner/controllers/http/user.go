package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	errs "github.com/pavlegich/banners-service/internal/errors"
	"github.com/pavlegich/banners-service/internal/infra/logger"
	"github.com/pavlegich/banners-service/internal/utils"
	"go.uber.org/zap"
)

// HandleGetUserBanner handles user's request to get banner by filter.
func (h *BannerHandler) HandleGetUserBanner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req requestQuery
	want := map[string]bool{
		"feature_id":        false,
		"tag_id":            false,
		"use_last_revision": false,
	}

	w.Header().Set("Content-Type", "application/json")

	queries := r.URL.Query()
	for val := range queries {
		_, ok := want[val]
		if !ok {
			logger.Log.Error("HandleGetUserBanner: incorrect query",
				zap.String("query", val))

			w.WriteHeader(http.StatusBadRequest)
			resp := utils.ParamToJSON("error", "incorrect query in request url")
			w.Write(resp)
			return
		}

		if len(queries[val]) != 1 {
			logger.Log.Error("HandleGetUserBanner: incorrect queries number",
				zap.String("query_name", val),
				zap.Int("query_number", len(queries[val])))

			w.WriteHeader(http.StatusBadRequest)
			resp := utils.ParamToJSON("error", "incorrect query number in request url")
			w.Write(resp)
			return
		}

		switch val {
		case "feature_id", "tag_id":
			current, err := strconv.Atoi(queries[val][0])
			if err != nil {
				logger.Log.Error("HandleGetUserBanner: convert query to integer failed",
					zap.String("query_name", val),
					zap.String("query_value", queries[val][0]))

				w.WriteHeader(http.StatusBadRequest)
				resp := utils.ParamToJSON("error", "convert query to integer failed")
				w.Write(resp)
				return
			}

			if val == "feature_id" {
				want["feature_id"] = true
				req.featureID = current
			}
			if val == "tag_id" {
				want["tag_id"] = true
				req.tagID = current
			}

		case "use_last_revision":
			current, err := strconv.ParseBool(queries[val][0])
			if err != nil {
				logger.Log.Error("HandleGetUserBanner: convert query to bool failed",
					zap.String("query_name", val),
					zap.String("query_value", queries[val][0]))

				w.WriteHeader(http.StatusBadRequest)
				resp := utils.ParamToJSON("error", "convert query to bool failed")
				w.Write(resp)
				return
			}

			req.lastRevision = current
		}
	}

	if !want["feature_id"] || !want["tag_id"] {
		logger.Log.Error("HandleGetUserBanner: required queries not set",
			zap.Bool("feature_id", want["feature_id"]),
			zap.Bool("tag_id", want["tag_id"]))

		w.WriteHeader(http.StatusBadRequest)
		resp := utils.ParamToJSON("error", "required queries not set")
		w.Write(resp)
		return
	}

	bannerContent, err := h.Service.Unload(ctx, req.featureID, req.tagID, req.lastRevision)
	if err != nil {
		logger.Log.Error("HandleGetUserBanner: get user banner failed",
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

	bannerJSON, err := json.Marshal(bannerContent)
	if err != nil {
		logger.Log.Error("HandleGetUserBanner: marshal banner content failed",
			zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		resp := utils.ParamToJSON("error", err.Error())
		w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bannerJSON)
}
