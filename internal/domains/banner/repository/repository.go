// Package repository contains repository object
// and methods for interaction with storage.
package repository

import (
	"context"
	"database/sql"
	"fmt"

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

// GetBannerByFilter gets benner from the storage by the requested filters and returns it.
func (r *Repository) GetBannerByFilter(ctx context.Context, name string, feature_id int, tag_id int) (*banner.Banner, error) {
	return nil, nil
}

// CreateBanner stores new banner into the storage.
func (r *Repository) CreateBanner(ctx context.Context, b *banner.Banner) (int, error) {
	row := r.db.QueryRowContext(ctx, `INSERT INTO banners (tag_ids, feature_id, content, is_active) 
	VALUES ($1, $2, $3, $4) RETURNING id`, b.TagIDs, b.FeatureID, b.Content, b.IsActive)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("CreateBanner: scan row failed %w", err)
	}

	err = row.Err()
	if err != nil {
		return -1, fmt.Errorf("CreateBanner: row.Err %w", err)
	}

	return id, nil
}

// GetBannersByFilter gets and returns the banners by filter from the storage.
func (r *Repository) GetBannersByFilter(ctx context.Context, feature_id int, tag_id int, limit int, offset int) ([]*banner.Banner, error) {
	query := "SELECT id, tag_ids, feature_id, content, is_active, created_at, updated_at FROM banners"
	if feature_id != 0 || tag_id != 0 {
		query += " WHERE"
		if feature_id != 0 && tag_id == 0 {
			query += fmt.Sprintf(" feature_id = %d", feature_id)
		} else if feature_id == 0 && tag_id != 0 {
			query += fmt.Sprintf(" %d = ANY (tag_ids)", tag_id)
		} else {
			query += fmt.Sprintf(" feature_id = %d AND %d = ANY (tag_ids)", feature_id, tag_id)
		}
	}

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

// UpdateBannerByID updates requested banner in the storage.
func (r *Repository) UpdateBannerByID(ctx context.Context, b *banner.Banner) error {
	res, err := r.db.ExecContext(ctx, `UPDATE banners SET tag_ids = $1, feature_id = $2, content = $3, is_active = $4,
	updated_at = NOW() WHERE id = $5`, b.TagIDs, b.FeatureID, b.Content, b.IsActive, b.ID)

	if err != nil {
		return fmt.Errorf("UpdateBannerByID: update data failed %w", err)
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateBannerByID: couldn't get rows affected %w", err)
	}
	if rowsCount == 0 {
		return fmt.Errorf("UpdateBannerByID: nothing to update, %w", errs.ErrBannerNotFound)
	}

	return nil
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
