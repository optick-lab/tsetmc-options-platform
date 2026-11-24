package repositories

import (
	"context"
	"errors"

	"backend/src/api/http/dto"
	"backend/src/models"

	"gorm.io/gorm"
)

// RawResponseRepository defines database operations for raw responses.
type RawResponseRepository interface {
	GetAll(ctx context.Context) ([]dto.RawResponseListItemDto, error)
	GetByID(ctx context.Context, id int64) (*models.RawResponse, error)
}

type rawResponseRepository struct {
	db *gorm.DB
}

// NewRawResponseRepository creates a new RawResponseRepository instance.
func NewRawResponseRepository(db *gorm.DB) RawResponseRepository {
	return &rawResponseRepository{db: db}
}

// GetAll retrieves all raw responses without the body, ordered by trade_date DESC, fetched_at DESC.
func (r *rawResponseRepository) GetAll(ctx context.Context) ([]dto.RawResponseListItemDto, error) {
	if r.db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var items []dto.RawResponseListItemDto
	err := r.db.WithContext(ctx).
		Model(&models.RawResponse{}).
		Select("id, endpoint, ins_code, http_status, is_blocked, fetched_at, LENGTH(body) as body_size, COALESCE(SUBSTRING(endpoint FROM '/GetTradeHistory/[^/]+/([0-9]{8})'), SUBSTRING(url FROM '/GetTradeHistory/[^/]+/([0-9]{8})'), SUBSTRING(endpoint FROM '([0-9]{8})')) as trade_date").
		Order("trade_date DESC NULLS LAST, fetched_at DESC, id DESC").
		Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetByID retrieves a single raw response by primary key ID (including body).
func (r *rawResponseRepository) GetByID(ctx context.Context, id int64) (*models.RawResponse, error) {
	if r.db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var record models.RawResponse
	err := r.db.WithContext(ctx).First(&record, id).Error
	if err != nil {
		return nil, err
	}

	return &record, nil
}
