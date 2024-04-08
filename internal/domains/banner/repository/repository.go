// Package repository contains repository object
// and methods for interaction with storage.
package repository

import (
	"context"
	"database/sql"

	"github.com/pavlegich/banners-service/internal/domains/banner"
)

// Repository contains storage objects for storing the banners.
type Repository struct {
	db *sql.DB
}

// NewBannerRepository returns new banners repository object.
func NewBannerRepository(ctx context.Context, db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetBannerByFilter gets benner from the storage by the requested filters and returns it.
func (r *Repository) GetBannerByFilter(ctx context.Context, name string, feature_id int, tag_id int) (*banner.Banner, error) {
	return nil, nil
}

// CreateBanner stores new banner into the storage.
func (r *Repository) CreateBanner(ctx context.Context, banner *banner.Banner) (int, error) {
	return -1, nil
}

// GetBannersByFilter gets and returns the banners by filter from the storage.
func (r *Repository) GetBannersByFilter(ctx context.Context, feature_id int, tag_id int, limit int, offset int) ([]*banner.Banner, error) {
	return nil, nil
}

// UpdateBanner updates requested banner in the storage.
func (r *Repository) UpdateBanner(ctx context.Context, banner *banner.Banner) error {
	return nil
}

// DeleteBannerByID deletes the requested by ID banner from the storage.
func (r *Repository) DeleteBannerByID(ctx context.Context, id int) error {
	return nil
}
