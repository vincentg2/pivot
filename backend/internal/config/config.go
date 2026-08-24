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
	WebDistDir         string
	DatabaseURL        string
	SessionCookieName  string
	SessionSecure      bool
	SessionTTL         time.Duration
	LoginRateLimit     int
	FootballDataAPIKey string
	FootaoEnabled      bool
	FootaoUserAgent    string
	NewsUserAgent      string
	SetupToken         string
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
	address := os.Getenv("APP_ADDR")
	if address == "" {
		address = ":" + env("PORT", "8080")
	}
	baseURL := os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		baseURL = env("RENDER_EXTERNAL_URL", "http://localhost:5173")
	}
	cfg := Config{
		Environment:        env("APP_ENV", "development"),
		Address:            address,
		BaseURL:            strings.TrimRight(baseURL, "/"),
		WebDistDir:         os.Getenv("WEB_DIST_DIR"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		SessionCookieName:  env("SESSION_COOKIE_NAME", "pivot_session"),
		SessionSecure:      secure,
		SessionTTL:         ttl,
		LoginRateLimit:     limit,
		FootballDataAPIKey: os.Getenv("FOOTBALL_DATA_API_KEY"),
		FootaoEnabled:      footaoEnabled,
		FootaoUserAgent:    env("FOOTAO_USER_AGENT", "Pivot/0.1 (+https://github.com/OWNER/pivot)"),
		NewsUserAgent:      env("NEWS_USER_AGENT", "Pivot/0.1 (+https://github.com/OWNER/pivot)"),
		SetupToken:         os.Getenv("SETUP_TOKEN"),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.FootaoEnabled && (strings.TrimSpace(cfg.FootaoUserAgent) == "" || strings.Contains(cfg.FootaoUserAgent, "OWNER")) {
		return Config{}, fmt.Errorf("FOOTAO_USER_AGENT must identify the operator when FOOTAO_ENABLED is true")
	}
	if cfg.SetupToken != "" && len(cfg.SetupToken) < 20 {
		return Config{}, fmt.Errorf("SETUP_TOKEN must contain at least 20 characters")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
