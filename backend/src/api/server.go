package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/src/api/routes"
	"backend/src/config"
	"backend/src/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// NewServer configures and returns the single Gin engine.
func NewServer() *gin.Engine {
	cfg := config.Get()
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestUUIDMiddleware())
	router.Use(corsMiddleware(cfg.AppAllowedOrigins))

	routes.RegisterRoutes(router)
	return router
}

func requestUUIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqUUID := c.GetHeader("X-Request-ID")
		if reqUUID == "" {
			reqUUID = uuid.NewString()
		}
		c.Set("request_uuid", reqUUID)
		c.Header("X-Request-ID", reqUUID)
		c.Next()
	}
}

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	originsMap := make(map[string]bool)
	allowAll := false
	for _, o := range allowedOrigins {
		if o == "*" {
			allowAll = true
		}
		originsMap[o] = true
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			if allowAll || originsMap[origin] {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-Request-ID")
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Start launches the HTTP server on APP_PORT and manages graceful shutdown.
func Start() error {
	cfg := config.Get()
	router := NewServer()

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

	// Listen and serve routine
	g.Go(func() error {
		logger.Info("Starting HTTP server", zap.String("port", cfg.AppPort), zap.String("env", cfg.AppEnv))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server failed", zap.Error(err))
			return err
		}
		return nil
	})

	// Graceful shutdown routine (10s timeout)
	g.Go(func() error {
		<-gCtx.Done()
		logger.Info("Shutdown signal received, shutting down HTTP server gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("Server forced to shutdown", zap.Error(err))
			return fmt.Errorf("forced server shutdown: %w", err)
		}

		logger.Info("Server stopped gracefully")
		return nil
	})

	return g.Wait()
}
