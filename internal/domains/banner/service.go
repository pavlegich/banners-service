package banner

import (
	"context"
	"errors"
	"fmt"

	errs "github.com/pavlegich/banners-service/internal/errors"
	"github.com/pavlegich/banners-service/internal/utils"
)

// BannerService contains objects for banner service.
type BannerService struct {
	repo  Repository
	cache Cache
}

// NewBannerService returns new banner service.
func NewBannerService(ctx context.Context, repo Repository, cache Cache) *BannerService {
	return &BannerService{
		repo:  repo,
		cache: cache,
	}
}

// Unload gets banner by filter and returns it.
func (s *BannerService) Unload(ctx context.Context, featureID int, tagID int, lastRevision bool) (*Content, error) {
	userRole, err := utils.GetUserRoleFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unload: get user role from context failed %w", err)
	}

	if !lastRevision {
		banner, err := s.cache.GetBannerByFilter(ctx, featureID, tagID)

		if err != nil {
			// If banner expired, delete banner from cache and get banner from the database
			if errors.Is(err, errs.ErrBannerExpired) {
				err := s.cache.DeleteBanner(ctx, 0, featureID, tagID)
				if err != nil {
					return nil, fmt.Errorf("Unload: delete banner from cache failed %w", err)
				}
				// If banner not found, get banner from the database
			} else if !errors.Is(err, errs.ErrBannerInCacheNotFound) {
				return nil, fmt.Errorf("Unload: get user banner content from cache failed %w", err)
			}
		} else { // If banner found, check whether the banner is active for user and return it
			if !banner.IsActive && userRole == "user" {
				return nil, fmt.Errorf("Unload: banner currently not active for users %w", errs.ErrBannerNotAllowed)
			}
			return banner.Content, nil
		}
	}

	banner, err := s.repo.GetBannerByFilter(ctx, featureID, tagID)
	if err != nil {
		return nil, fmt.Errorf("Unload: get actual user banner content failed %w", err)
	}

	if !banner.IsActive && userRole == "user" {
		return nil, fmt.Errorf("Unload: banner currently not active for users %w", errs.ErrBannerNotAllowed)
	}
	return banner.Content, nil
}

// Create creates new banner and puts it into the storage.
func (s *BannerService) Create(ctx context.Context, banner *Banner) (int, error) {
	storedBanner, err := s.repo.CreateBanner(ctx, banner)
	if err != nil {
		return -1, fmt.Errorf("Create: create banner failed %w", err)
	}

	err = s.cache.CreateBanner(ctx, storedBanner)
	if err != nil {
		return -1, fmt.Errorf("Create: create banner in cache failed %w", err)
	}

	return storedBanner.ID, nil
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
	storedBanner, err := s.repo.UpdateBanner(ctx, banner)
	if err != nil {
		return fmt.Errorf("Update: update banner failed %w", err)
	}

	err = s.cache.CreateBanner(ctx, storedBanner)
	if err != nil {
		return fmt.Errorf("Update: create banner in cache failed %w", err)
	}

	return nil
}

// Delete deletes the requested banner by ID from the storage.
func (s *BannerService) Delete(ctx context.Context, id int) error {
	err := s.repo.DeleteBannerByID(ctx, id)
	if err != nil {
		return fmt.Errorf("Delete: delete banner failed %w", err)
	}

	err = s.cache.DeleteBanner(ctx, id, 0, 0)
	if err != nil {
		return fmt.Errorf("Delete: delete banner from cache failed %w", err)
	}

	return nil
}
