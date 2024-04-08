package banner

import (
	"context"
)

// BannerService contains objects for banner service.
type BannerService struct {
	repo Repository
}

// NewBannerService returns new banner service.
func NewBannerService(ctx context.Context, repo Repository) *BannerService {
	return &BannerService{
		repo: repo,
	}
}

// Create creates new requested banner and requests repository to put it into the storage.
func (s *BannerService) Create(ctx context.Context, banner *Banner) error {
	return nil
}

// List returns list of available banners stored in the database.
func (s *BannerService) List(ctx context.Context) ([]*Banner, error) {
	return nil, nil
}

// Unload gets banner by banner's name and returns it.
func (s *BannerService) Unload(ctx context.Context, name string) (*Banner, error) {
	return nil, nil
}
