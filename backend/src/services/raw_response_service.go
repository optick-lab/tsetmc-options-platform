package services

import (
	"context"
	"errors"
	"strings"

	"backend/src/api/http/dto"
	"backend/src/repositories"

	"gorm.io/gorm"
)

var (
	// ErrRawResponseNotFound indicates that the requested raw response ID does not exist.
	ErrRawResponseNotFound = errors.New("raw response not found")

	// ErrNotGetTradeHistory indicates that the raw response endpoint is not a GetTradeHistory endpoint.
	ErrNotGetTradeHistory = errors.New("endpoint is not GetTradeHistory")
)

// RawResponseService defines business logic for stored raw responses.
type RawResponseService interface {
	GetRawResponses(ctx context.Context) ([]dto.RawResponseListItemDto, error)
	GetTradesByRawResponseID(ctx context.Context, id int64) (*dto.TradesResponseDto, error)
}

type rawResponseService struct {
	repo repositories.RawResponseRepository
}

// NewRawResponseService creates a new RawResponseService instance.
func NewRawResponseService(repo repositories.RawResponseRepository) RawResponseService {
	return &rawResponseService{repo: repo}
}

// GetRawResponses returns the list of archived responses without the body.
func (s *rawResponseService) GetRawResponses(ctx context.Context) ([]dto.RawResponseListItemDto, error) {
	return s.repo.GetAll(ctx)
}

// GetTradesByRawResponseID validates the endpoint and parses trades from the response body.
func (s *rawResponseService) GetTradesByRawResponseID(ctx context.Context, id int64) (*dto.TradesResponseDto, error) {
	record, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRawResponseNotFound
		}
		return nil, err
	}

	if !strings.Contains(record.Endpoint, "GetTradeHistory") {
		return nil, ErrNotGetTradeHistory
	}

	return ParseTradeHistoryJSON([]byte(record.Body))
}
