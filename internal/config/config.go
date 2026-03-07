package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL   string
	Port          string
	CheckInterval time.Duration
}

func Load() *Config {
	intervalSec := getEnvInt("CHECK_INTERVAL_SEC", 30)

	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://monitor:monitor@localhost:5432/monitor?sslmode=disable"),
		Port:          getEnv("PORT", "8080"),
		CheckInterval: time.Duration(intervalSec) * time.Second,
	}
}

func (c *Config) Addr() string {
	return fmt.Sprintf(":%s", c.Port)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}
