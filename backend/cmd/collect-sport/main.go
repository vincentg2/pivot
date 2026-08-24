package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/vincentg2/pivot/backend/internal/config"
	"github.com/vincentg2/pivot/backend/internal/database"
	"github.com/vincentg2/pivot/backend/internal/football"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fail(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	pool, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		fail(err)
	}
	defer pool.Close()
	service := football.NewService(football.NewPostgresRepository(pool), football.NewFootballDataConnector(cfg.FootballDataAPIKey))
	run, err := service.Sync(ctx, time.Now())
	if err != nil {
		fail(err)
	}
	fmt.Printf("Sports collection succeeded: %d records\n", run.RecordsCount)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "Sports collection failed:", err)
	os.Exit(1)
}
