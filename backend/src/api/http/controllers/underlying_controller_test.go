package controllers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/src/api/http/controllers"
	"backend/src/api/http/dto"
	"backend/src/api/routes"
	"backend/src/providers"

	"github.com/gin-gonic/gin"
)

// Unit tests for ParsePaginationParams
func TestParsePaginationParams(t *testing.T) {
	tests := []struct {
		name         string
		pageStr      string
		pageSizeStr  string
		expectedPage int
		expectedSize int
		expectErr    bool
	}{
		{
			name:         "Default values when empty",
			pageStr:      "",
			pageSizeStr:  "",
			expectedPage: 1,
			expectedSize: 20,
			expectErr:    false,
		},
		{
			name:         "Valid custom page and page_size",
			pageStr:      "2",
			pageSizeStr:  "10",
			expectedPage: 2,
			expectedSize: 10,
			expectErr:    false,
		},
		{
			name:         "Valid boundary max page_size of 100",
			pageStr:      "5",
			pageSizeStr:  "100",
			expectedPage: 5,
			expectedSize: 100,
			expectErr:    false,
		},
		{
			name:        "Invalid page zero",
			pageStr:     "0",
			pageSizeStr: "20",
			expectErr:   true,
		},
		{
			name:        "Invalid negative page",
			pageStr:     "-1",
			pageSizeStr: "20",
			expectErr:   true,
		},
		{
			name:        "Invalid non-integer page",
			pageStr:     "abc",
			pageSizeStr: "20",
			expectErr:   true,
		},
		{
			name:        "Invalid page_size zero",
			pageStr:     "1",
			pageSizeStr: "0",
			expectErr:   true,
		},
		{
			name:        "Invalid negative page_size",
			pageStr:     "1",
			pageSizeStr: "-10",
			expectErr:   true,
		},
		{
			name:        "Invalid page_size exceeding 100",
			pageStr:     "1",
			pageSizeStr: "101",
			expectErr:   true,
		},
		{
			name:        "Invalid non-integer page_size",
			pageStr:     "1",
			pageSizeStr: "xyz",
			expectErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p, ps, err := controllers.ParsePaginationParams(tc.pageStr, tc.pageSizeStr)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("Expected error for inputs page=%q, page_size=%q, got nil", tc.pageStr, tc.pageSizeStr)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if p != tc.expectedPage {
					t.Errorf("Expected page %d, got %d", tc.expectedPage, p)
				}
				if ps != tc.expectedSize {
					t.Errorf("Expected page_size %d, got %d", tc.expectedSize, ps)
				}
			}
		})
	}
}

// MockUnderlyingService simulates business logic responses.
type mockUnderlyingService struct {
	getUnderlyingsFn func(ctx context.Context, page, pageSize int) (*dto.PaginatedUnderlyingsDto, error)
}

func (m *mockUnderlyingService) GetUnderlyings(ctx context.Context, page, pageSize int) (*dto.PaginatedUnderlyingsDto, error) {
	if m.getUnderlyingsFn != nil {
		return m.getUnderlyingsFn(ctx, page, pageSize)
	}
	return nil, errors.New("not implemented")
}

type paginatedEnvelopeResponse struct {
	IsSuccessful bool                    `json:"is_successful"`
	RequestUUID  string                  `json:"request_uuid"`
	StatusCode   int                     `json:"status_code"`
	Message      string                  `json:"message"`
	Data         dto.PaginatedUnderlyingsDto `json:"data"`
	Items        []dto.UnderlyingItemDto `json:"items"`
	Page         int                     `json:"page"`
	PageSize     int                     `json:"page_size"`
	Total        int64                   `json:"total"`
	Errors       []string                `json:"errors"`
	ErrorCode    int                     `json:"error_code"`
}

func TestUnderlyingController_List_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockUnderlyingService{
		getUnderlyingsFn: func(ctx context.Context, page, pageSize int) (*dto.PaginatedUnderlyingsDto, error) {
			items := make([]dto.UnderlyingItemDto, pageSize)
			for i := 0; i < pageSize; i++ {
				items[i] = dto.UnderlyingItemDto{
					ID:        int64(i + 1),
					Symbol:    "خودرو",
					InsCode:   nil,
					AssetType: "stock",
				}
			}
			return &dto.PaginatedUnderlyingsDto{
				Items:    items,
				Page:     page,
				PageSize: pageSize,
				Total:    25,
			}, nil
		},
	}

	providers.SetUnderlyingService(mockSvc)

	router := gin.New()
	routes.RegisterRoutes(router)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/underlyings?page=1&page_size=10", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var res paginatedEnvelopeResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !res.IsSuccessful {
		t.Errorf("Expected is_successful=true, got false")
	}

	// Verify both data and top-level mirrored pagination fields
	if len(res.Data.Items) != 10 {
		t.Errorf("Expected 10 items in data.items, got %d", len(res.Data.Items))
	}
	if res.Data.Total != 25 {
		t.Errorf("Expected data.total=25, got %d", res.Data.Total)
	}
	if len(res.Items) != 10 {
		t.Errorf("Expected 10 items at root, got %d", len(res.Items))
	}
	if res.Total != 25 {
		t.Errorf("Expected total=25 at root, got %d", res.Total)
	}
	if res.Page != 1 || res.PageSize != 10 {
		t.Errorf("Expected page=1, page_size=10, got page=%d, page_size=%d", res.Page, res.PageSize)
	}
}

func TestUnderlyingController_List_Page3_Returns5Items(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockUnderlyingService{
		getUnderlyingsFn: func(ctx context.Context, page, pageSize int) (*dto.PaginatedUnderlyingsDto, error) {
			// For page 3 with page_size 10 out of 25 total, return remaining 5
			items := make([]dto.UnderlyingItemDto, 5)
			for i := 0; i < 5; i++ {
				items[i] = dto.UnderlyingItemDto{
					ID:        int64(20 + i + 1),
					Symbol:    "همتراز",
					InsCode:   nil,
					AssetType: "stock",
				}
			}
			return &dto.PaginatedUnderlyingsDto{
				Items:    items,
				Page:     page,
				PageSize: pageSize,
				Total:    25,
			}, nil
		},
	}

	providers.SetUnderlyingService(mockSvc)

	router := gin.New()
	routes.RegisterRoutes(router)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/underlyings?page=3&page_size=10", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var res paginatedEnvelopeResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(res.Data.Items) != 5 {
		t.Errorf("Expected 5 items, got %d", len(res.Data.Items))
	}
	if res.Data.Total != 25 {
		t.Errorf("Expected total=25, got %d", res.Data.Total)
	}
}

func TestUnderlyingController_List_InvalidQueryParams_400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	routes.RegisterRoutes(router)

	invalidUrls := []string{
		"/api/v1/underlyings?page_size=0",
		"/api/v1/underlyings?page=-1",
		"/api/v1/underlyings?page=0",
		"/api/v1/underlyings?page_size=101",
		"/api/v1/underlyings?page=abc",
		"/api/v1/underlyings?page_size=xyz",
	}

	for _, u := range invalidUrls {
		t.Run(u, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, u, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status 400 for %s, got %d: %s", u, w.Code, w.Body.String())
			}

			var res paginatedEnvelopeResponse
			if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
				t.Fatalf("Failed to parse error response: %v", err)
			}

			if res.IsSuccessful {
				t.Errorf("Expected is_successful=false for %s", u)
			}
			if len(res.Errors) == 0 {
				t.Errorf("Expected non-empty errors array for %s", u)
			}
		})
	}
}
