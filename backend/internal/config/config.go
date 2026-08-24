package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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
	FootaoEnabled      bool
	FootaoUserAgent    string
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
	footaoEnabled, err := strconv.ParseBool(env("FOOTAO_ENABLED", "false"))
	if err != nil {
		return Config{}, fmt.Errorf("parse FOOTAO_ENABLED: %w", err)
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
		FootaoEnabled:      footaoEnabled,
		FootaoUserAgent:    env("FOOTAO_USER_AGENT", "Pivot/0.1 (+https://github.com/OWNER/pivot)"),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.FootaoEnabled && (strings.TrimSpace(cfg.FootaoUserAgent) == "" || strings.Contains(cfg.FootaoUserAgent, "OWNER")) {
		return Config{}, fmt.Errorf("FOOTAO_USER_AGENT must identify the operator when FOOTAO_ENABLED is true")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
