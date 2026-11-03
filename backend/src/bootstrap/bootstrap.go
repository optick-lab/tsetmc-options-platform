package bootstrap

import (
	"fmt"

	"backend/src/api"
	"backend/src/config"
	"backend/src/database"
	"backend/src/pkg/logger"

	"go.uber.org/zap"
)

// Bootstrap initializes logger, database, runs migrations, and starts the API server.
func Bootstrap() error {
	cfg := config.Get()
	logger.Init(cfg.AppEnv)
	logger.Info("Bootstrapping application...", zap.String("env", cfg.AppEnv))

	if _, err := database.Connect(); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	if err := database.MigrateUp(); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	if err := api.Start(); err != nil {
		return fmt.Errorf("api server error: %w", err)
	}

	return nil
}

// Serve initializes logger, database, and starts the API server without running migrations.
func Serve() error {
	cfg := config.Get()
	logger.Init(cfg.AppEnv)
	logger.Info("Starting application server...", zap.String("env", cfg.AppEnv))

	if _, err := database.Connect(); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	if err := api.Start(); err != nil {
		return fmt.Errorf("api server error: %w", err)
	}

	return nil
}
