package database

import (
	"fmt"
	"sync"
	"time"

	"backend/src/config"
	"backend/src/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

// Connect initializes the GORM PostgreSQL connection singleton.
func Connect() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		cfg := config.Get()
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
		)

		logLevel := gormlogger.Warn
		if cfg.AppEnv == "local" {
			logLevel = gormlogger.Info
		}

		conn, openErr := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormlogger.Default.LogMode(logLevel),
		})
		if openErr != nil {
			err = openErr
			logger.Error("Failed to connect to database", zap.Error(err))
			return
		}

		sqlDB, dbErr := conn.DB()
		if dbErr != nil {
			err = dbErr
			logger.Error("Failed to get generic sql.DB", zap.Error(err))
			return
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetConnMaxLifetime(time.Hour)

		db = conn
	})

	if err != nil {
		once = sync.Once{}
		return nil, err
	}

	return db, nil
}

// GetDB returns the active GORM database instance.
func GetDB() *gorm.DB {
	if db == nil {
		_, _ = Connect()
	}
	return db
}

// SetDB allows overriding the database instance (e.g. for testing).
func SetDB(override *gorm.DB) {
	db = override
}
