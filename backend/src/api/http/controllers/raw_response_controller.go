package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"backend/src/api/http/response"
	"backend/src/providers"
	"backend/src/services"

	"github.com/gin-gonic/gin"
)

// RawResponseController handles HTTP requests for stored raw responses and trade extraction.
type RawResponseController struct {
	service services.RawResponseService
}

// NewRawResponseController creates a new RawResponseController instance.
func NewRawResponseController(svc ...services.RawResponseService) *RawResponseController {
	var s services.RawResponseService
	if len(svc) > 0 {
		s = svc[0]
	}
	return &RawResponseController{service: s}
}

func (ctrl *RawResponseController) getService() services.RawResponseService {
	if ctrl.service != nil {
		return ctrl.service
	}
	return providers.GetRawResponseService()
}

// List handles GET /api/v1/raw-responses
// Returns the list of stored responses without their bodies.
func (ctrl *RawResponseController) List(c *gin.Context) {
	svc := ctrl.getService()
	items, err := svc.GetRawResponses(c.Request.Context())
	if err != nil {
		response.Api(c).
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("failed to retrieve raw responses").
			SetErrors([]string{err.Error()}).
			SetErrorCode(50000).
			Send()
		return
	}

	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage("success").
		SetData(items).
		Send()
}

// GetTrades handles GET /api/v1/raw-responses/:id/trades
// Extracts and summarizes trades from the stored GetTradeHistory response body.
func (ctrl *RawResponseController) GetTrades(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		response.Api(c).
			SetStatusCode(http.StatusBadRequest).
			SetMessage("invalid raw response id").
			SetErrors([]string{"id must be a positive integer"}).
			SetErrorCode(40000).
			Send()
		return
	}

	svc := ctrl.getService()
	result, err := svc.GetTradesByRawResponseID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrRawResponseNotFound) {
			response.Api(c).
				SetStatusCode(http.StatusNotFound).
				SetMessage("raw response not found").
				SetErrors([]string{err.Error()}).
				SetErrorCode(40400).
				Send()
			return
		}

		if errors.Is(err, services.ErrNotGetTradeHistory) {
			response.Api(c).
				SetStatusCode(http.StatusBadRequest).
				SetMessage("endpoint is not GetTradeHistory").
				SetErrors([]string{err.Error()}).
				SetErrorCode(40000).
				Send()
			return
		}

		response.Api(c).
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("failed to parse trades from raw response").
			SetErrors([]string{err.Error()}).
			SetErrorCode(50000).
			Send()
		return
	}

	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage("success").
		SetData(result).
		Send()
}
