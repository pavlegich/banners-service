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

// Unload gets banner by filter and returns it.
func (s *BannerService) Unload(ctx context.Context, feature_id int, tag_id int, actual bool) (*Banner, error) {
	return nil, nil
}

// Create creates new banner and puts it into the storage.
func (s *BannerService) Create(ctx context.Context, banner *Banner) (int, error) {
	return -1, nil
}

// List returns list of banners by filter stored in the storage.
func (s *BannerService) List(ctx context.Context, feature_id int, tag_id int, limit int, offset int) ([]*Banner, error) {
	return nil, nil
}

// Update updates the requested banner.
func (s *BannerService) Update(ctx context.Context, banner *Banner) error {
	return nil
}

// Delete deletes the requested banner by ID from the storage.
func (s *BannerService) Delete(ctx context.Context, id int) error {
	return nil
}
