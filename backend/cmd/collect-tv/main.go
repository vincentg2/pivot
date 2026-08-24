package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/vincentg2/pivot/backend/internal/broadcast"
	"github.com/vincentg2/pivot/backend/internal/config"
	"github.com/vincentg2/pivot/backend/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fail(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	pool, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		fail(err)
	}
	defer pool.Close()
	service := broadcast.NewService(broadcast.NewPostgresRepository(pool), broadcast.NewFootaoConnector(cfg.FootaoEnabled, cfg.FootaoUserAgent))
	run, err := service.Sync(ctx, time.Now())
	if err != nil {
		fail(err)
	}
	fmt.Printf("TV collection succeeded: %d listings\n", run.RecordsCount)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "TV collection failed:", err)
	os.Exit(1)
}
