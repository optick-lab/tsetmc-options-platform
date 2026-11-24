package database

import (
	"embed"
	"errors"
	"fmt"

	"backend/src/config"
	"backend/src/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func getMigrateInstance() (*migrate.Migrate, error) {
	driver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to load embedded migrations: %w", err)
	}

	cfg := config.Get()
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	m, err := migrate.NewWithSourceInstance("iofs", driver, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return m, nil
}

// MigrateUp executes all pending migrations.
func MigrateUp() error {
	m, err := getMigrateInstance()
	if err != nil {
		return err
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			logger.Warn("Failed to close migration source", zap.Error(srcErr))
		}
		if dbErr != nil {
			logger.Warn("Failed to close migration database", zap.Error(dbErr))
		}
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration up failed: %w", err)
	}

	logger.Info("Database migrations applied successfully")
	return nil
}

// MigrateDown rolls back migrations (all by default, or by specified number of steps).
func MigrateDown(steps ...int) error {
	m, err := getMigrateInstance()
	if err != nil {
		return err
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			logger.Warn("Failed to close migration source", zap.Error(srcErr))
		}
		if dbErr != nil {
			logger.Warn("Failed to close migration database", zap.Error(dbErr))
		}
	}()

	if len(steps) > 0 && steps[0] > 0 {
		if err := m.Steps(-steps[0]); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migration down failed: %w", err)
		}
		logger.Info("Database migrations rolled back partially", zap.Int("steps", steps[0]))
		return nil
	}

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration down failed: %w", err)
	}

	logger.Info("Database migrations rolled back successfully")
	return nil
}

// MigrateForce sets a specific migration version.
func MigrateForce(version int) error {
	m, err := getMigrateInstance()
	if err != nil {
		return err
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			logger.Warn("Failed to close migration source", zap.Error(srcErr))
		}
		if dbErr != nil {
			logger.Warn("Failed to close migration database", zap.Error(dbErr))
		}
	}()

	if err := m.Force(version); err != nil {
		return fmt.Errorf("migration force failed: %w", err)
	}

	logger.Info("Database migration version forced successfully", zap.Int("version", version))
	return nil
}
