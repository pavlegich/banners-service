// Package repository contains repository object
// and methods for interaction with storage.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/pavlegich/banners-service/internal/domains/banner"
	errs "github.com/pavlegich/banners-service/internal/errors"
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

// GetBannerByFilter gets and returns banner content from the storage by the requested filters.
func (r *Repository) GetBannerByFilter(ctx context.Context, featureID int, tagID int) (*banner.Banner, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, tag_ids, feature_id, content, is_active, created_at, updated_at 
	FROM banners WHERE feature_id = $1 AND $2 = ANY (tag_ids) AND is_active = true 
	ORDER BY updated_at DESC LIMIT 1`, featureID, tagID)

	var b banner.Banner
	var tagIDs pq.Int64Array
	err := row.Scan(&b.ID, &tagIDs, &b.FeatureID, &b.Content, &b.IsActive, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("GetBannerByFilter: banner not found in database %w", errs.ErrBannerNotFound)
		}
		return nil, fmt.Errorf("GetBannerByFilter: scan row failed %w", err)
	}
	for _, v := range tagIDs {
		b.TagIDs = append(b.TagIDs, int(v))
	}

	err = row.Err()
	if err != nil {
		return nil, fmt.Errorf("GetBannerByFilter:: row.Err() %w", err)
	}

	return &b, nil
}

// CreateBanner stores new banner into the storage.
func (r *Repository) CreateBanner(ctx context.Context, b *banner.Banner) (*banner.Banner, error) {
	row := r.db.QueryRowContext(ctx, `INSERT INTO banners (tag_ids, feature_id, content, is_active) 
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`, b.TagIDs, b.FeatureID, b.Content, b.IsActive)

	var id int
	var createdAt, updatedAt time.Time
	err := row.Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("CreateBanner: scan row failed %w", err)
	}

	b.ID = id
	b.CreatedAt = createdAt
	b.UpdatedAt = updatedAt

	err = row.Err()
	if err != nil {
		return nil, fmt.Errorf("CreateBanner: row.Err %w", err)
	}

	return b, nil
}

// GetBannersByFilter gets and returns the banners by filter from the storage.
func (r *Repository) GetBannersByFilter(ctx context.Context, featureID int, tagID int, limit int, offset int) ([]*banner.Banner, error) {
	query := "SELECT id, tag_ids, feature_id, content, is_active, created_at, updated_at FROM banners"
	if featureID != 0 || tagID != 0 {
		query += " WHERE"
		if featureID != 0 && tagID == 0 {
			query += fmt.Sprintf(" feature_id = %d", featureID)
		} else if featureID == 0 && tagID != 0 {
			query += fmt.Sprintf(" %d = ANY (tag_ids)", tagID)
		} else {
			query += fmt.Sprintf(" feature_id = %d AND %d = ANY (tag_ids)", featureID, tagID)
		}
	}

	query += " ORDER BY updated_at DESC"

	if limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	query += fmt.Sprintf(" OFFSET %d", offset)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetBannersByFilter: read rows from table failed %w", err)
	}
	defer rows.Close()

	bannersList := make([]*banner.Banner, 0)
	for rows.Next() {
		var b banner.Banner
		var tagIDs pq.Int64Array
		err = rows.Scan(&b.ID, &tagIDs, &b.FeatureID, &b.Content, &b.IsActive, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("GetBannersByFilter: scan row failed %w", err)
		}
		for _, v := range tagIDs {
			b.TagIDs = append(b.TagIDs, int(v))
		}
		bannersList = append(bannersList, &b)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("GetBannersByFilter: rows.Err %w", err)
	}

	return bannersList, nil
}

// UpdateBanner updates requested banner in the storage.
func (r *Repository) UpdateBanner(ctx context.Context, b *banner.Banner) (*banner.Banner, error) {
	row := r.db.QueryRowContext(ctx, `UPDATE banners SET tag_ids = $1, feature_id = $2, content = $3, is_active = $4,
	updated_at = NOW() WHERE id = $5 RETURNING updated_at`, b.TagIDs, b.FeatureID, b.Content, b.IsActive, b.ID)

	var updatedAt time.Time
	err := row.Scan(&updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("UpdateBanner: nothing to update, %w", errs.ErrBannerNotFound)
		}
		return nil, fmt.Errorf("UpdateBanner: scan row failed %w", err)
	}

	b.UpdatedAt = updatedAt

	err = row.Err()
	if err != nil {
		return nil, fmt.Errorf("UpdateBanner: row.Err %w", err)
	}

	return b, nil
}

// DeleteBannerByID deletes the requested by ID banner from the storage.
func (r *Repository) DeleteBannerByID(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM banners WHERE id = $1`, id)

	if err != nil {
		return fmt.Errorf("DeleteBannerByID: delete data failed %w", err)
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteBannerByID: couldn't get rows affected %w", err)
	}
	if rowsCount == 0 {
		return fmt.Errorf("DeleteBannerByID: nothing to delete, %w", errs.ErrBannerNotFound)
	}

	return nil
}
