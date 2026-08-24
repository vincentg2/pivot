package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/vincentg2/pivot/backend/internal/catalog"
	"github.com/vincentg2/pivot/backend/internal/config"
	"github.com/vincentg2/pivot/backend/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fail(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	pool, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		fail(err)
	}
	defer pool.Close()
	service := catalog.NewService(catalog.NewPostgresRepository(pool), catalog.NewFootballDataConnector(cfg.FootballDataAPIKey))
	run, err := service.Sync(ctx)
	if err != nil {
		fail(err)
	}
	fmt.Printf("Catalog collection succeeded: %d clubs\n", run.RecordsCount)
}
func fail(err error) { fmt.Fprintln(os.Stderr, "Catalog collection failed:", err); os.Exit(1) }
