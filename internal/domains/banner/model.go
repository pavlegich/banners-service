// Package banner contains object and methods
// for interacting with banners.
package banner

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Banner contains data for banners.
type Banner struct {
	ID        int       `json:"banner_id"`
	TagIDs    []int     `json:"tag_ids"`
	FeatureID int       `json:"feature_id"`
	Content   *Content  `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	UpdateBannerByID(ctx context.Context, banner *Banner) error
	DeleteBannerByID(ctx context.Context, id int) error
}

// Content type for implementing the Scanner interface.
type Content map[string]string

// Scan implements Scan method for scanning the banner content from the storage.
func (c *Content) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	switch data := v.(type) {
	case string:
		return json.Unmarshal([]byte(data), &c)
	case []byte:
		return json.Unmarshal(data, &c)
	default:
		return fmt.Errorf("cannot scan type %t into Map", v)
	}
}
