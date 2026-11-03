package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/src/api/routes"
	"backend/src/config"

	"github.com/gin-gonic/gin"
)

type responseEnvelope struct {
	IsSuccessful bool                   `json:"is_successful"`
	RequestUUID  string                 `json:"request_uuid"`
	StatusCode   int                    `json:"status_code"`
	Message      string                 `json:"message"`
	Data         map[string]interface{} `json:"data"`
	Errors       interface{}            `json:"errors"`
	ErrorCode    int                    `json:"error_code"`
}

func TestHealthController_Health(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = config.Get()

	router := gin.New()
	routes.RegisterRoutes(router)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	var res responseEnvelope
	if err := json.Unmarshal(recorder.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to parse response body JSON: %v. Body: %s", err, recorder.Body.String())
	}

	if !res.IsSuccessful {
		t.Errorf("Expected is_successful=true, got %v", res.IsSuccessful)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status_code=%d, got %d", http.StatusOK, res.StatusCode)
	}

	if res.Message != "ok" {
		t.Errorf("Expected message='ok', got %s", res.Message)
	}

	if res.Data["status"] != "ok" {
		t.Errorf("Expected data.status='ok', got %v", res.Data["status"])
	}

	if res.RequestUUID == "" {
		t.Errorf("Expected non-empty request_uuid")
	}
}
