// Package banner contains object and methods
// for interacting with banners.
package banner

import (
	"context"
	"time"
)

// Banner contains data for banners.
type Banner struct {
	ID        int               `json:"id"`
	TagIDs    []int             `json:"tag_ids"`
	FeatureID int               `json:"feature_id"`
	Content   map[string]string `json:"content"`
	IsActive  bool              `json:"is_active"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Service describes methods for communication between
// handlers and repositories.
type Service interface {
	Unload(ctx context.Context, feature_id int, tag_id int, actual bool) (*Banner, error)
	Create(ctx context.Context, banner *Banner) (int, error)
	List(ctx context.Context, feature_id int, tag_id int, limit int, offset int) ([]*Banner, error)
	Update(ctx context.Context, banner *Banner) error
	Delete(ctx context.Context, id int) error
}

// Repository describes methods related with banners
// for interaction with the storage.
type Repository interface {
	GetBannerByFilter(ctx context.Context, name string, feature_id int, tag_id int) (*Banner, error)
	CreateBanner(ctx context.Context, banner *Banner) (int, error)
	GetBannersByFilter(ctx context.Context, feature_id int, tag_id int, limit int, offset int) ([]*Banner, error)
	UpdateBanner(ctx context.Context, banner *Banner) error
	DeleteBannerByID(ctx context.Context, id int) error
}
