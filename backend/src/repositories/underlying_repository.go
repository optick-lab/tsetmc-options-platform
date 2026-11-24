package repositories

import (
	"context"

	"backend/src/models"

	"gorm.io/gorm"
)

// UnderlyingRepository defines database operations for underlying assets.
type UnderlyingRepository interface {
	GetPaginated(ctx context.Context, offset int, limit int) ([]models.UnderlyingAsset, int64, error)
}

type underlyingRepository struct {
	db *gorm.DB
}

// NewUnderlyingRepository creates a new UnderlyingRepository instance.
func NewUnderlyingRepository(db *gorm.DB) UnderlyingRepository {
	return &underlyingRepository{db: db}
}

// GetPaginated retrieves underlying assets ordered by symbol ASC with total count.
func (r *underlyingRepository) GetPaginated(ctx context.Context, offset int, limit int) ([]models.UnderlyingAsset, int64, error) {
	var items []models.UnderlyingAsset
	var total int64

	db := r.db.WithContext(ctx).Model(&models.UnderlyingAsset{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("symbol ASC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
