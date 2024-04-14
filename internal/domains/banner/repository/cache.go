package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pavlegich/banners-service/internal/domains/banner"
	errs "github.com/pavlegich/banners-service/internal/errors"
)

// Cache contains data for cache object.
type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	banners           map[bannerKey]cacheBanner
}

// cacheBanner contains data for store banner in cache.
type cacheBanner struct {
	banner  *banner.Banner
	expires time.Time
}

// bannerKey contains data for unique banner search.
type bannerKey struct {
	featureID int
	tagID     int
}

// NewBannerCache creates and returns new banner cache.
func NewBannerCache(ctx context.Context, defaultExpiration time.Duration, cleanupInterval time.Duration) *Cache {
	return &Cache{
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		banners:           make(map[bannerKey]cacheBanner, 0),
	}
}

// GetBannerByFilter finds and returns requested banner content by filter.
func (c *Cache) GetBannerByFilter(ctx context.Context, featureID int, tagID int) (*banner.Banner, error) {
	c.RLock()
	defer c.RUnlock()

	key := bannerKey{
		featureID: featureID,
		tagID:     tagID,
	}

	cb, ok := c.banners[key]
	if !ok {
		return nil, fmt.Errorf("GetBannerByFilter: banners with requested tag not found %w", errs.ErrBannerInCacheNotFound)
	}

	if time.Now().After(cb.expires) {
		return nil, fmt.Errorf("GetBannerByFilter: banner content usage expired %w", errs.ErrBannerExpired)
	}

	return cb.banner, nil
}

// CreateBanner creates new banner in cache.
func (c *Cache) CreateBanner(ctx context.Context, banner *banner.Banner) error {
	c.Lock()
	defer c.Unlock()

	for _, tagID := range banner.TagIDs {
		key := bannerKey{
			featureID: banner.FeatureID,
			tagID:     tagID,
		}

		c.banners[key] = cacheBanner{
			banner:  banner,
			expires: banner.UpdatedAt.Add(c.defaultExpiration),
		}
	}

	return nil
}

// DeleteBanner deletes banner from cache.
func (c *Cache) DeleteBanner(ctx context.Context, id int, featureID int, tagID int) error {
	c.Lock()
	defer c.Unlock()

	if id > 0 {
		for k := range c.banners {
			if c.banners[k].banner.ID == id {
				delete(c.banners, k)
			}
		}

		return nil
	}

	key := bannerKey{
		featureID: featureID,
		tagID:     tagID,
	}

	_, ok := c.banners[key]
	if !ok {
		return fmt.Errorf("DeleteBanner: requested banner not found %w", errs.ErrBannerInCacheNotFound)
	}

	delete(c.banners, key)

	return nil
}

// GarbageCollect cleans banners cache with requested interval.
func (c *Cache) GarbageCollect(ctx context.Context) {
	ticker := time.NewTicker(c.cleanupInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if c.banners == nil {
				return
			}

			keys := c.expiredKeys()
			if len(keys) != 0 {
				c.clearBanners(keys)
			}

		default:
			continue
		}
	}

}

// expiredKeys returns list of expired keys.
func (c *Cache) expiredKeys() (keys []bannerKey) {
	c.RLock()
	defer c.RUnlock()

	for key, b := range c.banners {
		if time.Now().After(b.expires) {
			keys = append(keys, key)
		}
	}

	return
}

// clearBanners deletes expired keys.
func (c *Cache) clearBanners(keys []bannerKey) {
	c.Lock()
	defer c.Unlock()

	for _, k := range keys {
		delete(c.banners, k)
	}
}
