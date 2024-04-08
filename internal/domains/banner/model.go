// Package banner contains object and methods
// for interacting with banners.
package banner

import "context"

// Banner contains data for banners.
type Banner struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Script string `json:"script"`
}

// Service describes methods for communication between
// handlers and repositories.
type Service interface {
	Create(ctx context.Context, banner *Banner) error
	List(ctx context.Context) ([]*Banner, error)
	Unload(ctx context.Context, name string) (*Banner, error)
}

// Repository describes methods related with banners
// for interaction with database.
type Repository interface {
	CreateBanner(ctx context.Context, banner *Banner) error
	GetAllBanners(ctx context.Context) ([]*Banner, error)
	GetBannerByName(ctx context.Context, name string) (*Banner, error)
}
