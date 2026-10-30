package controllers

import (
	"net/http"

	"backend/src/api/http/response"
	"backend/src/database"

	"github.com/gin-gonic/gin"
)

// HealthController handles health check endpoints.
type HealthController struct{}

// NewHealthController creates a new HealthController instance.
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Health handles GET /api/v1/health.
func (h *HealthController) Health(c *gin.Context) {
	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage("ok").
		SetData(gin.H{"status": "ok"}).
		Send()
}

// DBHealth handles GET /api/v1/health/db by executing a ping against PostgreSQL.
// Note: Per architecture rules, health is the single permitted exception where controller touches DB directly.
func (h *HealthController) DBHealth(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.Api(c).
			SetStatusCode(http.StatusServiceUnavailable).
			SetMessage("database connection unavailable").
			SetErrors([]string{"database connection is nil"}).
			SetErrorCode(50301).
			Send()
		return
	}

	var result int
	if err := db.Raw("SELECT 1").Scan(&result).Error; err != nil {
		response.Api(c).
			SetStatusCode(http.StatusServiceUnavailable).
			SetMessage("database ping failed").
			SetErrors([]string{err.Error()}).
			SetErrorCode(50302).
			Send()
		return
	}

	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage("database connection healthy").
		SetData(gin.H{"status": "ok", "db": "connected"}).
		Send()
}
