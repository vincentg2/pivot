package tests

import (
	"context"
	"crypto/sha256"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/vincentg2/pivot/backend/internal/auth"
	"github.com/vincentg2/pivot/backend/internal/database"
	"golang.org/x/crypto/bcrypt"
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
	code := strings.ToUpper("PIVOT-INTEGRATION-" + uuid.NewString())
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

func TestPasswordResetAgainstPostgres(t *testing.T) {
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
	adminID, userID := uuid.New(), uuid.New()
	oldHash, _ := bcrypt.GenerateFromPassword([]byte("old-safe-password"), bcrypt.MinCost)
	_, err = pool.Exec(ctx, `INSERT INTO users(id,email,password_hash,nickname,avatar_seed,role) VALUES($1,$2,$3,'Admin','admin-seed','admin'),($4,$5,$3,'Friend','friend-seed','user')`, adminID, uuid.NewString()+"@example.test", string(oldHash), userID, uuid.NewString()+"@example.test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id=ANY($1)`, []uuid.UUID{adminID, userID})
	})
	var email string
	if err := pool.QueryRow(ctx, `SELECT email FROM users WHERE id=$1`, userID).Scan(&email); err != nil {
		t.Fatal(err)
	}
	service := auth.NewService(auth.NewPostgresRepository(pool), time.Hour)
	reset, err := service.CreatePasswordReset(ctx, email, adminID)
	if err != nil {
		t.Fatal(err)
	}
	sessionHash := sha256.Sum256([]byte("existing-session"))
	if _, err = pool.Exec(ctx, `INSERT INTO sessions(id,token_hash,user_id,expires_at) VALUES($1,$2,$3,now()+interval '1 hour')`, uuid.New(), sessionHash[:], userID); err != nil {
		t.Fatal(err)
	}
	if err = service.ResetPassword(ctx, reset.Token, "new-integration-password"); err != nil {
		t.Fatal(err)
	}
	var passwordHash string
	var sessions int
	if err = pool.QueryRow(ctx, `SELECT password_hash FROM users WHERE id=$1`, userID).Scan(&passwordHash); err != nil {
		t.Fatal(err)
	}
	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte("new-integration-password")) != nil {
		t.Fatal("new password was not persisted")
	}
	if err = pool.QueryRow(ctx, `SELECT count(*) FROM sessions WHERE user_id=$1`, userID).Scan(&sessions); err != nil || sessions != 0 {
		t.Fatalf("existing sessions were not revoked: count=%d err=%v", sessions, err)
	}
	if err = service.ResetPassword(ctx, reset.Token, "another-safe-password"); !errors.Is(err, auth.ErrPasswordResetInvalid) {
		t.Fatalf("expected used reset link to be rejected, got %v", err)
	}
	expired, err := service.CreatePasswordReset(ctx, email, adminID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `UPDATE password_reset_links SET expires_at=now()-interval '1 minute' WHERE user_id=$1 AND used_at IS NULL`, userID); err != nil {
		t.Fatal(err)
	}
	if err = service.ResetPassword(ctx, expired.Token, "another-safe-password"); !errors.Is(err, auth.ErrPasswordResetInvalid) {
		t.Fatalf("expected expired reset link to be rejected, got %v", err)
	}
}
