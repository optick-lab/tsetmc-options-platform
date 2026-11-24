package services

import (
	"context"

	"backend/src/api/http/dto"
	"backend/src/repositories"
)

// UnderlyingService defines business logic for underlying assets.
type UnderlyingService interface {
	GetUnderlyings(ctx context.Context, page, pageSize int) (*dto.PaginatedUnderlyingsDto, error)
}

type underlyingService struct {
	repo repositories.UnderlyingRepository
}

// NewUnderlyingService creates a new UnderlyingService instance.
func NewUnderlyingService(repo repositories.UnderlyingRepository) UnderlyingService {
	return &underlyingService{repo: repo}
}

// GetUnderlyings computes offset and retrieves paginated underlying assets.
func (s *underlyingService) GetUnderlyings(ctx context.Context, page, pageSize int) (*dto.PaginatedUnderlyingsDto, error) {
	offset := (page - 1) * pageSize
	assets, total, err := s.repo.GetPaginated(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]dto.UnderlyingItemDto, len(assets))
	for i, a := range assets {
		items[i] = dto.UnderlyingItemDto{
			ID:        a.ID,
			Symbol:    a.Symbol,
			InsCode:   a.InsCode,
			AssetType: a.AssetType,
		}
	}

	return &dto.PaginatedUnderlyingsDto{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}
