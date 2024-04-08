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

// CreateBanner stores new banner into the storage.
func (r *Repository) CreateBanner(ctx context.Context, banner *banner.Banner) error {
	return nil
}

// GetAllBanners gets and returns all the banners from the storage.
func (r *Repository) GetAllBanners(ctx context.Context) ([]*banner.Banner, error) {
	return nil, nil
}

// GetBannerByName gets and returns the requested by name banner from the storage.
func (r *Repository) GetBannerByName(ctx context.Context, name string) (*banner.Banner, error) {
	return nil, nil
}
