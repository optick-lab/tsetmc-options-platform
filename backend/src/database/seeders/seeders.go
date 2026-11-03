package seeders

import (
	"backend/src/pkg/logger"

	"go.uber.org/zap"
)

// RunAll executes all registered database seeders (no-op for initial skeleton).
func RunAll() error {
	logger.Info("Database seeders executed", zap.String("status", "no-op"))
	return nil
}
