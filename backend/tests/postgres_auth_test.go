package tests

import (
	"context"
	"crypto/sha256"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/vincentg2/pivot/backend/internal/auth"
	"github.com/vincentg2/pivot/backend/internal/database"
)

func TestInvitationRegistrationAgainstPostgres(t *testing.T) {
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
	code := "PIVOT-INTEGRATION-" + uuid.NewString()
	hash := sha256.Sum256([]byte(code))
	inviteID := uuid.New()
	_, err = pool.Exec(ctx, `INSERT INTO invitations(id,code_hash,label,max_uses) VALUES($1,$2,'integration',1)`, inviteID, hash[:])
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM invitations WHERE id=$1`, inviteID) })
	service := auth.NewService(auth.NewPostgresRepository(pool), time.Hour)
	user, err := service.Register(ctx, uuid.NewString()+"@example.test", "integration-password", "Integration", code)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, user.ID) })
	if user.Role != "user" {
		t.Fatalf("unexpected role %q", user.Role)
	}
	var uses int
	if err := pool.QueryRow(ctx, `SELECT uses FROM invitations WHERE id=$1`, inviteID).Scan(&uses); err != nil {
		t.Fatal(err)
	}
	if uses != 1 {
		t.Fatalf("expected one use, got %d", uses)
	}
}
