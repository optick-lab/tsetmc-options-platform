package routes

import (
	"net/http"

	"backend/src/api/http/controllers"
	"backend/src/api/http/response"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts all application routes and custom error handlers.
func RegisterRoutes(router *gin.Engine) {
	// NoRoute handler -> 404 envelope with error_code 40400
	router.NoRoute(func(c *gin.Context) {
		response.Api(c).
			SetStatusCode(http.StatusNotFound).
			SetMessage("Route not found").
			SetErrorCode(40400).
			SetErrors([]string{"The requested endpoint does not exist"}).
			Send()
	})

	healthController := controllers.NewHealthController()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthController.Health)
		v1.GET("/health/db", healthController.DBHealth)
	}
}
