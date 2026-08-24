package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/vincentg2/pivot/backend/internal/config"
	"github.com/vincentg2/pivot/backend/internal/database"
	"github.com/vincentg2/pivot/backend/internal/news"
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
	service := news.NewService(news.NewPostgresRepository(pool), news.NewRSSConnector(cfg.NewsUserAgent))
	run, err := service.Sync(ctx, time.Now())
	if err != nil {
		fail(err)
	}
	fmt.Printf("News collection succeeded: %d items\n", run.RecordsCount)
}
func fail(err error) { fmt.Fprintln(os.Stderr, "News collection failed:", err); os.Exit(1) }
