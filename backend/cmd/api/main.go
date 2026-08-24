package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vincentg2/pivot/backend/internal/auth"
	"github.com/vincentg2/pivot/backend/internal/broadcast"
	"github.com/vincentg2/pivot/backend/internal/catalog"
	"github.com/vincentg2/pivot/backend/internal/config"
	"github.com/vincentg2/pivot/backend/internal/database"
	"github.com/vincentg2/pivot/backend/internal/football"
	"github.com/vincentg2/pivot/backend/internal/httpx"
	"github.com/vincentg2/pivot/backend/internal/installation"
	"github.com/vincentg2/pivot/backend/internal/invitation"
	"github.com/vincentg2/pivot/backend/internal/news"
	"github.com/vincentg2/pivot/backend/internal/platform/validate"
	"github.com/vincentg2/pivot/backend/internal/user"
	"golang.org/x/time/rate"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("configuration failed", "error", err)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	pool, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	authService := auth.NewService(auth.NewPostgresRepository(pool), cfg.SessionTTL)
	authHandler := auth.NewHandler(authService, cfg.SessionCookieName, cfg.SessionSecure, cfg.SessionTTL)
	inviteHandler := invitation.NewHandler(invitation.NewService(invitation.NewPostgresRepository(pool)))
	userHandler := user.NewHandler(user.NewService(user.NewPostgresRepository(pool)))
	catalogService := catalog.NewService(catalog.NewPostgresRepository(pool), catalog.NewFootballDataConnector(cfg.FootballDataAPIKey))
	catalogHandler := catalog.NewHandler(catalogService, logger)
	footballService := football.NewService(football.NewPostgresRepository(pool), football.NewFootballDataConnector(cfg.FootballDataAPIKey))
	footballHandler := football.NewHandler(footballService, logger)
	broadcastService := broadcast.NewService(broadcast.NewPostgresRepository(pool), broadcast.NewFootaoConnector(cfg.FootaoEnabled, cfg.FootaoUserAgent))
	broadcastHandler := broadcast.NewHandler(broadcastService, logger)
	newsService := news.NewService(news.NewPostgresRepository(pool), news.NewRSSConnector(cfg.NewsUserAgent))
	newsHandler := news.NewHandler(newsService, logger)
	installationHandler := installation.NewHandler(installation.NewService(installation.NewPostgresRepository(pool), cfg.SetupToken))

	e := echo.New()
	e.HideBanner = true
	e.Validator = validate.New()
	e.HTTPErrorHandler = httpx.ErrorHandler
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{cfg.BaseURL}, AllowCredentials: true, AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept}}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			started := time.Now()
			err := next(c)
			logger.Info("http request", "request_id", c.Response().Header().Get(echo.HeaderXRequestID), "method", c.Request().Method, "path", c.Path(), "status", c.Response().Status, "duration_ms", time.Since(started).Milliseconds())
			return err
		}
	})

	e.GET("/health", func(c echo.Context) error {
		pingCtx, cancel := context.WithTimeout(c.Request().Context(), time.Second)
		defer cancel()
		if err := pool.Ping(pingCtx); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unhealthy"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	api := e.Group("/api/v1")
	loginLimiter := middleware.RateLimiter(newLoginRateLimiter(cfg.LoginRateLimit))
	api.GET("/setup/status", installationHandler.Status)
	api.POST("/setup", installationHandler.Install, loginLimiter)
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login, loginLimiter)
	api.POST("/auth/logout", authHandler.Logout)
	api.GET("/auth/me", authHandler.Me, authHandler.RequireSession)
	api.PATCH("/profile", userHandler.Update, authHandler.RequireSession)
	api.DELETE("/profile", userHandler.Delete, authHandler.RequireSession)
	api.GET("/competitions", catalogHandler.Competitions, authHandler.RequireSession)
	api.GET("/clubs", catalogHandler.Clubs, authHandler.RequireSession)
	api.GET("/clubs/:id", catalogHandler.Club, authHandler.RequireSession)
	api.GET("/favorites", catalogHandler.Favorites, authHandler.RequireSession)
	api.PUT("/favorites", catalogHandler.ReplaceFavorites, authHandler.RequireSession)
	api.GET("/matches", footballHandler.Matches, authHandler.RequireSession)
	api.GET("/standings", footballHandler.Standing, authHandler.RequireSession)
	api.GET("/broadcasts", broadcastHandler.List, authHandler.RequireSession)
	api.GET("/news", newsHandler.List, authHandler.RequireSession)
	admin := api.Group("/admin", authHandler.RequireAdmin)
	admin.GET("/invitations", inviteHandler.List)
	admin.POST("/invitations", inviteHandler.Create)
	admin.DELETE("/invitations/:id", inviteHandler.Revoke)
	admin.GET("/collections/football-data", catalogHandler.CollectionStatus)
	admin.POST("/collections/football-data", catalogHandler.Sync)
	admin.GET("/collections/football-data/sport", footballHandler.CollectionStatus)
	admin.POST("/collections/football-data/sport", footballHandler.Sync)
	admin.GET("/collections/footao", broadcastHandler.CollectionStatus)
	admin.POST("/collections/footao", broadcastHandler.Sync)
	admin.GET("/broadcasts", broadcastHandler.AdminList)
	admin.PUT("/broadcasts/:id/correction", broadcastHandler.Correct)
	admin.DELETE("/broadcasts/:id/correction", broadcastHandler.Clear)
	admin.GET("/broadcasts/audit", broadcastHandler.Audit)
	admin.GET("/news/feeds", newsHandler.Feeds)
	admin.POST("/news/feeds", newsHandler.SaveFeed)
	admin.DELETE("/news/feeds/:id", newsHandler.DeleteFeed)
	admin.GET("/collections/news", newsHandler.CollectionStatus)
	admin.POST("/collections/news", newsHandler.Sync)

	go func() {
		logger.Info("server starting", "address", cfg.Address, "environment", cfg.Environment)
		if err := e.Start(cfg.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			stop()
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", "error", err)
	}
}

func newLoginRateLimiter(requestsPerMinute int) middleware.RateLimiterStore {
	return middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:  rate.Limit(float64(requestsPerMinute) / 60),
		Burst: requestsPerMinute,
	})
}
