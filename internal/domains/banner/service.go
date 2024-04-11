package banner

import (
	"context"
	"fmt"
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
func (s *BannerService) Unload(ctx context.Context, featureID int, tagID int, lastRevision bool) (*Content, error) {
	bannerContent, err := s.repo.GetBannerContentByFilter(ctx, featureID, tagID)
	if err != nil {
		return nil, fmt.Errorf("Unload: get user banner content failed %w", err)
	}

	return bannerContent, nil
}

// Create creates new banner and puts it into the storage.
func (s *BannerService) Create(ctx context.Context, banner *Banner) (int, error) {
	id, err := s.repo.CreateBanner(ctx, banner)
	if err != nil {
		return -1, fmt.Errorf("Create: create banner failed %w", err)
	}

	return id, nil
}

// List returns list of banners by filter stored in the storage.
func (s *BannerService) List(ctx context.Context, featureID int, tagID int, limit int, offset int) ([]*Banner, error) {
	bannersList, err := s.repo.GetBannersByFilter(ctx, featureID, tagID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("List: get banners list by filter failed %w", err)
	}

	return bannersList, nil
}

// Update updates the requested banner.
func (s *BannerService) Update(ctx context.Context, banner *Banner) error {
	err := s.repo.UpdateBannerByID(ctx, banner)
	if err != nil {
		return fmt.Errorf("Update: update banner failed %w", err)
	}

	return nil
}

// Delete deletes the requested banner by ID from the storage.
func (s *BannerService) Delete(ctx context.Context, id int) error {
	err := s.repo.DeleteBannerByID(ctx, id)
	if err != nil {
		return fmt.Errorf("Delete: delete banner failed %w", err)
	}

	return nil
}
