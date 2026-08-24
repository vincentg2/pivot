package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment        string
	Address            string
	BaseURL            string
	DatabaseURL        string
	SessionCookieName  string
	SessionSecure      bool
	SessionTTL         time.Duration
	LoginRateLimit     int
	FootballDataAPIKey string
}

func Load() (Config, error) {
	ttl, err := time.ParseDuration(env("SESSION_TTL", "720h"))
	if err != nil {
		return Config{}, fmt.Errorf("parse SESSION_TTL: %w", err)
	}
	secure, err := strconv.ParseBool(env("SESSION_SECURE", "false"))
	if err != nil {
		return Config{}, fmt.Errorf("parse SESSION_SECURE: %w", err)
	}
	limit, err := strconv.Atoi(env("LOGIN_RATE_LIMIT", "10"))
	if err != nil || limit < 1 {
		return Config{}, fmt.Errorf("LOGIN_RATE_LIMIT must be a positive integer")
	}
	cfg := Config{
		Environment:        env("APP_ENV", "development"),
		Address:            env("APP_ADDR", ":8080"),
		BaseURL:            env("APP_BASE_URL", "http://localhost:5173"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		SessionCookieName:  env("SESSION_COOKIE_NAME", "pivot_session"),
		SessionSecure:      secure,
		SessionTTL:         ttl,
		LoginRateLimit:     limit,
		FootballDataAPIKey: os.Getenv("FOOTBALL_DATA_API_KEY"),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
