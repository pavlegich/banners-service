package http_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pavlegich/banners-service/internal/controllers/handlers"
	"github.com/pavlegich/banners-service/internal/domains/banner"
	errs "github.com/pavlegich/banners-service/internal/errors"
	"github.com/pavlegich/banners-service/internal/infra/config"
	"github.com/pavlegich/banners-service/internal/mocks"
	"github.com/stretchr/testify/assert"
)

var bannersList = map[string]*banner.Banner{
	"ok": {
		ID:        1,
		TagIDs:    []int{1, 2, 3},
		FeatureID: 1,
		Content: &banner.Content{
			"title": "some_title",
			"text":  "some_text",
			"url":   "some_url",
		},
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now().Add(time.Duration(1) * time.Minute),
	},
	"not_active": {
		ID:        1,
		TagIDs:    []int{1, 2, 3},
		FeatureID: 1,
		Content: &banner.Content{
			"title": "some_title",
			"text":  "some_text",
			"url":   "some_url",
		},
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now().Add(time.Duration(1) * time.Minute),
	},
	"expired": {
		ID:        1,
		TagIDs:    []int{1, 2, 3},
		FeatureID: 1,
		Content: &banner.Content{
			"title": "some_title",
			"text":  "some_text",
			"url":   "some_url",
		},
		IsActive:  true,
		CreatedAt: time.Now().Add(-time.Duration(10) * time.Minute),
		UpdatedAt: time.Now().Add(-time.Duration(10) * time.Minute),
	},
}

func TestBannerHandler_HandleGetUserBanner(t *testing.T) {
	ctx := context.Background()

	// Initialize the mocks for storage
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := mocks.NewMockRepository(mockCtrl)
	mockCache := mocks.NewMockCache(mockCtrl)

	gomock.InOrder(
		// ok for user with cache
		mockCache.EXPECT().GetBannerByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(bannersList["ok"], nil),

		// ok with using last revision
		mockRepo.EXPECT().GetBannerByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(bannersList["ok"], nil),

		// banner not active
		mockRepo.EXPECT().GetBannerByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(bannersList["not_active"], nil),

		// banner expired in cache
		mockCache.EXPECT().GetBannerByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, errs.ErrBannerExpired),

		mockCache.EXPECT().DeleteBanner(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil),

		mockRepo.EXPECT().GetBannerByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(bannersList["expired"], nil),

		// ok with admin using last revision
		mockRepo.EXPECT().GetBannerByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(bannersList["ok"], nil),

		// ok with admin using last revision when is_active false
		mockRepo.EXPECT().GetBannerByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(bannersList["not_active"], nil),

		// banner for admin using last revision not found
		mockRepo.EXPECT().GetBannerByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, errs.ErrBannerNotFound),
	)

	cfg := &config.Config{}

	type args struct {
		token        string
		featureID    int
		tagID        int
		lastRevision bool
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "ok from cache",
			args: args{
				token:        "user_token",
				featureID:    1,
				tagID:        1,
				lastRevision: false,
			},
			wantCode: http.StatusOK,
		},
		{
			name: "ok with using last revision",
			args: args{
				token:        "user_token",
				featureID:    1,
				tagID:        1,
				lastRevision: true,
			},
			wantCode: http.StatusOK,
		},
		{
			name: "banner not active",
			args: args{
				token:        "user_token",
				featureID:    1,
				tagID:        1,
				lastRevision: true,
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "banner expired in cache",
			args: args{
				token:        "user_token",
				featureID:    1,
				tagID:        1,
				lastRevision: false,
			},
			wantCode: http.StatusOK,
		},
		{
			name: "ok with admin using last revision",
			args: args{
				token:        "admin_token",
				featureID:    1,
				tagID:        1,
				lastRevision: true,
			},
			wantCode: http.StatusOK,
		},
		{
			name: "ok with admin using last revision when banner not active",
			args: args{
				token:        "admin_token",
				featureID:    1,
				tagID:        1,
				lastRevision: true,
			},
			wantCode: http.StatusOK,
		},
		{
			name: "unexpected queries in url",
			args: args{
				token:        "user_token",
				featureID:    1,
				tagID:        -1,
				lastRevision: true,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "unexpected user role",
			args: args{
				token:        "unexpected_user",
				featureID:    1,
				tagID:        1,
				lastRevision: true,
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "banner for admin using last revision not found",
			args: args{
				token:        "admin_token",
				featureID:    1,
				tagID:        1,
				lastRevision: true,
			},
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Controller
			ctrl := handlers.NewController(ctx, mockRepo, mockCache, cfg)
			mh, err := ctrl.BuildRoute(ctx)
			assert.NoError(t, err)

			// Form the new request with requested queries
			url := `http://localhost:8080/user_banner`
			url += fmt.Sprintf("?feature_id=%d&tag_id=%d", tt.args.featureID, tt.args.tagID)
			if tt.args.lastRevision {
				url += fmt.Sprintf("&use_last_revision=%t", tt.args.lastRevision)
			}
			r := httptest.NewRequest(http.MethodGet, url, nil)
			r.Header.Set("token", tt.args.token)
			w := httptest.NewRecorder()

			mh.ServeHTTP(w, r)

			// Get response
			resp := w.Result()
			defer resp.Body.Close()

			gotBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			// Check status code
			gotCode := resp.StatusCode
			if gotCode != tt.wantCode {
				t.Errorf("BannerHandler.HandleGetUserBanner() = %v, want %v. Error: %s", gotCode, tt.wantCode, string(gotBody))
			}
		})
	}
}
