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

// ParsePaginationParams extracts and validates page and page_size query parameters.
func ParsePaginationParams(pageStr, pageSizeStr string) (int, int, error) {
	page := 1
	pageSize := 20

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			return 0, 0, errors.New("page must be an integer greater than or equal to 1")
		}
		page = p
	}

	if pageSizeStr != "" {
		ps, err := strconv.Atoi(pageSizeStr)
		if err != nil || ps < 1 {
			return 0, 0, errors.New("page_size must be an integer greater than or equal to 1")
		}
		if ps > 100 {
			return 0, 0, errors.New("page_size must not exceed 100")
		}
		pageSize = ps
	}

	return page, pageSize, nil
}

// UnderlyingController handles HTTP requests for underlying assets.
type UnderlyingController struct {
	service services.UnderlyingService
}

// NewUnderlyingController creates a new UnderlyingController instance.
func NewUnderlyingController(svc ...services.UnderlyingService) *UnderlyingController {
	var s services.UnderlyingService
	if len(svc) > 0 {
		s = svc[0]
	}
	return &UnderlyingController{service: s}
}

// List handles GET /api/v1/underlyings
func (ctrl *UnderlyingController) List(c *gin.Context) {
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	page, pageSize, err := ParsePaginationParams(pageStr, pageSizeStr)
	if err != nil {
		response.Api(c).
			SetStatusCode(http.StatusBadRequest).
			SetMessage(err.Error()).
			SetErrors([]string{err.Error()}).
			SetErrorCode(40000).
			Send()
		return
	}

	svc := ctrl.service
	if svc == nil {
		svc = providers.GetUnderlyingService()
	}

	result, err := svc.GetUnderlyings(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Api(c).
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("failed to retrieve underlying assets").
			SetErrors([]string{err.Error()}).
			SetErrorCode(50000).
			Send()
		return
	}

	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage("success").
		SetPaginatedData(result.Items, result.Page, result.PageSize, result.Total).
		Send()
}
