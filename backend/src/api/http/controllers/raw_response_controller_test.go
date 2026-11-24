package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend/src/api/http/dto"
	"backend/src/services"

	"github.com/gin-gonic/gin"
)

type mockRawResponseService struct {
	responses []dto.RawResponseListItemDto
	tradesMap map[int64]*dto.TradesResponseDto
	errMap    map[int64]error
}

func (m *mockRawResponseService) GetRawResponses(ctx context.Context) ([]dto.RawResponseListItemDto, error) {
	return m.responses, nil
}

func (m *mockRawResponseService) GetTradesByRawResponseID(ctx context.Context, id int64) (*dto.TradesResponseDto, error) {
	if err, exists := m.errMap[id]; exists {
		return nil, err
	}
	if trades, exists := m.tradesMap[id]; exists {
		return trades, nil
	}
	return nil, services.ErrRawResponseNotFound
}

func setupRawResponseTestRouter(svc services.RawResponseService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	ctrl := NewRawResponseController(svc)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/raw-responses", ctrl.List)
		v1.GET("/raw-responses/:id/trades", ctrl.GetTrades)
	}

	return router
}

func TestRawResponseController_List_Success(t *testing.T) {
	insCode := "31459957843277944"
	tradeDate := "20261111"
	mockSvc := &mockRawResponseService{
		responses: []dto.RawResponseListItemDto{
			{
				ID:         1,
				Endpoint:   "/Trade/GetTradeHistory/31459957843277944/20261111/false",
				InsCode:    &insCode,
				HttpStatus: 200,
				IsBlocked:  false,
				FetchedAt:  time.Now(),
				BodySize:   2723,
				TradeDate:  &tradeDate,
			},
		},
	}

	router := setupRawResponseTestRouter(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/raw-responses", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp struct {
		IsSuccessful bool                         `json:"is_successful"`
		Data         []dto.RawResponseListItemDto `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.IsSuccessful || len(resp.Data) != 1 {
		t.Errorf("expected 1 raw response, got %d", len(resp.Data))
	}

	if resp.Data[0].TradeDate == nil || *resp.Data[0].TradeDate != "20261111" {
		t.Errorf("expected trade_date='20261111', got %v", resp.Data[0].TradeDate)
	}
}

func TestRawResponseController_GetTrades_Success(t *testing.T) {
	mockTrades := &dto.TradesResponseDto{
		Summary: dto.TradeSummaryDto{
			TradeCount:  14,
			TotalVolume: 1251,
			TotalValue:  377073.0,
			MinPrice:    290.0,
			MaxPrice:    313.0,
			FirstPrice:  290.0,
			LastPrice:   299.0,
		},
		Trades: []dto.TradeItemDto{
			{
				NTran:    14,
				Time:     "12:29:40",
				Volume:   10,
				Price:    299.0,
				Canceled: false,
			},
		},
	}

	mockSvc := &mockRawResponseService{
		tradesMap: map[int64]*dto.TradesResponseDto{
			1: mockTrades,
		},
	}

	router := setupRawResponseTestRouter(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/raw-responses/1/trades", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		IsSuccessful bool                  `json:"is_successful"`
		Data         dto.TradesResponseDto `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Summary.TradeCount != 14 || resp.Data.Summary.TotalVolume != 1251 {
		t.Errorf("unexpected summary in response: %+v", resp.Data.Summary)
	}
}

func TestRawResponseController_GetTrades_NotFound_404(t *testing.T) {
	mockSvc := &mockRawResponseService{}
	router := setupRawResponseTestRouter(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/raw-responses/999/trades", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestRawResponseController_GetTrades_NotGetTradeHistory_400(t *testing.T) {
	mockSvc := &mockRawResponseService{
		errMap: map[int64]error{
			2: services.ErrNotGetTradeHistory,
		},
	}
	router := setupRawResponseTestRouter(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/raw-responses/2/trades", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestRawResponseController_GetTrades_InvalidId_400(t *testing.T) {
	mockSvc := &mockRawResponseService{}
	router := setupRawResponseTestRouter(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/raw-responses/invalid_id/trades", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
