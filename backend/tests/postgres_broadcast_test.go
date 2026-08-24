package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/vincentg2/pivot/backend/internal/broadcast"
	"github.com/vincentg2/pivot/backend/internal/database"
)

func TestBroadcastCorrectionAuditAgainstPostgres(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := database.Open(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	adminID := uuid.New()
	_, err = pool.Exec(ctx, `INSERT INTO users(id,email,password_hash,nickname,avatar_seed,role) VALUES($1,$2,'test','Admin','seed','admin')`, adminID, uuid.NewString()+"@example.test")
	if err != nil {
		t.Fatal(err)
	}
	repository := broadcast.NewPostgresRepository(pool)
	startsAt := time.Now().Add(2 * time.Hour).Truncate(time.Second)
	sourceURL := "https://www.footao.tv/" + uuid.NewString()
	item := broadcast.ImportedListing{SourceKey: uuid.NewString(), StartsAt: startsAt, HomeName: "Home", AwayName: "Away", Label: "Home · Away", CompetitionName: "League", Kind: "live", Channels: []string{"Channel One"}, SourceURL: &sourceURL}
	if err = repository.Upsert(ctx, item, nil); err != nil {
		t.Fatal(err)
	}
	var listingID uuid.UUID
	if err = pool.QueryRow(ctx, `SELECT id FROM tv_listings WHERE source_key=$1`, item.SourceKey).Scan(&listingID); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM tv_listing_audit WHERE listing_id=$1`, listingID)
		_, _ = pool.Exec(ctx, `DELETE FROM tv_listings WHERE id=$1`, listingID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, adminID)
	}()
	item.SourceKey = uuid.NewString()
	item.StartsAt = startsAt.Add(15 * time.Minute)
	item.Channels = []string{"Channel Updated"}
	if err = repository.Upsert(ctx, item, nil); err != nil {
		t.Fatal(err)
	}
	var importedID uuid.UUID
	var importedCount int
	if err = pool.QueryRow(ctx, `SELECT min(id::text)::uuid,count(*) FROM tv_listings WHERE source='footao' AND source_url=$1`, sourceURL).Scan(&importedID, &importedCount); err != nil {
		t.Fatal(err)
	}
	if importedCount != 1 || importedID != listingID {
		t.Fatalf("schedule update created a duplicate listing: count=%d id=%s", importedCount, importedID)
	}
	startsAt = item.StartsAt

	input := broadcast.CorrectionInput{StartsAt: startsAt.Add(30 * time.Minute), HomeName: "Corrected Home", AwayName: "Away", Label: "Corrected Home · Away", CompetitionName: "League", Kind: "delayed", Channels: []string{"Channel Two"}, Hidden: true, Note: "integration"}
	if err = repository.Correct(ctx, listingID, adminID, input); err != nil {
		t.Fatal(err)
	}
	items, err := repository.List(ctx, startsAt.Add(-24*time.Hour), startsAt.Add(24*time.Hour), true)
	if err != nil {
		t.Fatal(err)
	}
	var corrected *broadcast.Listing
	for index := range items {
		if items[index].ID == listingID {
			corrected = &items[index]
			break
		}
	}
	if corrected == nil || !corrected.Corrected || !corrected.Hidden || corrected.Label != input.Label || corrected.Channels[0] != "Channel Two" {
		t.Fatalf("correction did not override import: %#v", corrected)
	}
	if err = repository.ClearCorrection(ctx, listingID, adminID); err != nil {
		t.Fatal(err)
	}
	audit, err := repository.Audit(ctx, 100)
	if err != nil {
		t.Fatal(err)
	}
	seen := 0
	for _, entry := range audit {
		if entry.ListingID == listingID {
			seen++
		}
	}
	if seen != 2 {
		t.Fatalf("expected correction and restore audit entries, got %d", seen)
	}
}
