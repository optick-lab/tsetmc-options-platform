package config

import (
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv            string
	AppPort           string
	AppAllowedOrigins []string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
}

var (
	cfg  *Config
	once sync.Once
)

// Load reads configuration from environment variables and .env file.
func Load() *Config {
	once.Do(func() {
		_ = godotenv.Load()

		appEnv := getEnv("APP_ENV", "local")
		appPort := getEnv("APP_PORT", "8085")
		originsRaw := getEnv("APP_ALLOWED_ORIGINS", "http://localhost:3080,http://127.0.0.1:3080,http://localhost:3000,http://127.0.0.1:3000")
		var origins []string
		for _, o := range strings.Split(originsRaw, ",") {
			trimmed := strings.TrimSpace(o)
			if trimmed != "" {
				origins = append(origins, trimmed)
			}
		}

		cfg = &Config{
			AppEnv:            appEnv,
			AppPort:           appPort,
			AppAllowedOrigins: origins,
			DBHost:            getEnv("DB_HOST", "localhost"),
			DBPort:            getEnv("DB_PORT", "5440"),
			DBUser:            getEnv("DB_USER", "postgres"),
			DBPassword:        getEnv("DB_PASSWORD", "change_me"),
			DBName:            getEnv("DB_NAME", "tsetmc_options"),
		}
	})
	return cfg
}

// Get returns the singleton configuration instance.
func Get() *Config {
	if cfg == nil {
		return Load()
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
