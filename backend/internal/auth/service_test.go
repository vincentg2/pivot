package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type fakeRepository struct {
	user           User
	invitationHash []byte
	sessionHash    []byte
	expires        time.Time
	resetHash      []byte
	resetExpires   time.Time
	resetUserID    uuid.UUID
	resetAdminID   uuid.UUID
	passwordHash   string
	resetAt        time.Time
}

func (f *fakeRepository) FindUserByEmail(context.Context, string) (User, error) {
	if f.user.ID == uuid.Nil {
		return User{}, ErrNotFound
	}
	return f.user, nil
}
func (f *fakeRepository) FindUserByID(context.Context, uuid.UUID) (User, error) { return f.user, nil }
func (f *fakeRepository) RegisterWithInvitation(_ context.Context, input NewUser, hash []byte, _ time.Time) (User, error) {
	f.invitationHash = hash
	f.user = User{ID: input.ID, Email: input.Email, PasswordHash: input.PasswordHash, Nickname: input.Nickname, AvatarSeed: input.AvatarSeed, Theme: "system", Role: "user"}
	return f.user, nil
}
func (f *fakeRepository) CreateSession(_ context.Context, _ uuid.UUID, hash []byte, expires time.Time, _, _ string) error {
	f.sessionHash = hash
	f.expires = expires
	return nil
}
func (f *fakeRepository) FindUserBySession(context.Context, []byte) (User, error) { return f.user, nil }
func (f *fakeRepository) DeleteSession(context.Context, []byte) error             { return nil }
func (f *fakeRepository) CreatePasswordReset(_ context.Context, userID, adminID uuid.UUID, hash []byte, expires time.Time) error {
	f.resetUserID, f.resetAdminID, f.resetHash, f.resetExpires = userID, adminID, hash, expires
	return nil
}
func (f *fakeRepository) ConsumePasswordReset(_ context.Context, hash []byte, passwordHash string, now time.Time) error {
	f.resetHash, f.passwordHash, f.resetAt = hash, passwordHash, now
	return nil
}

func TestRegisterNormalizesIdentityAndHashesSecrets(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, time.Hour)
	user, err := service.Register(context.Background(), "  Friend@Example.COM ", "a-long-safe-password", "  Camille  ", " pivot-example ")
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != "friend@example.com" || user.Nickname != "Camille" {
		t.Fatalf("unexpected normalized user: %#v", user)
	}
	if string(repo.invitationHash) == "PIVOT-EXAMPLE" {
		t.Fatal("invitation code was stored in plaintext")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("a-long-safe-password")) != nil {
		t.Fatal("password was not bcrypt hashed")
	}
}

func TestLoginReturnsOpaqueTokenAndPersistsOnlyHash(t *testing.T) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("a-long-safe-password"), bcrypt.MinCost)
	repo := &fakeRepository{user: User{ID: uuid.New(), Email: "friend@example.com", PasswordHash: string(passwordHash)}}
	service := NewService(repo, 2*time.Hour)
	now := time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return now }
	_, token, err := service.Login(context.Background(), "friend@example.com", "a-long-safe-password", "test", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if token == "" || string(repo.sessionHash) == token {
		t.Fatal("session token is missing or stored in plaintext")
	}
	if !repo.expires.Equal(now.Add(2 * time.Hour)) {
		t.Fatalf("unexpected expiry: %s", repo.expires)
	}
	_, _, err = service.Login(context.Background(), "friend@example.com", "wrong", "test", "")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestAdministratorCreatesOneTimePasswordReset(t *testing.T) {
	userID, adminID := uuid.New(), uuid.New()
	repo := &fakeRepository{user: User{ID: userID, Email: "friend@example.com"}}
	service := NewService(repo, time.Hour)
	now := time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return now }
	reset, err := service.CreatePasswordReset(context.Background(), " Friend@Example.COM ", adminID)
	if err != nil {
		t.Fatal(err)
	}
	if reset.Token == "" || string(repo.resetHash) == reset.Token {
		t.Fatal("reset token is missing or stored in plaintext")
	}
	if repo.resetUserID != userID || repo.resetAdminID != adminID {
		t.Fatal("reset link was not associated with the target and administrator")
	}
	if !repo.resetExpires.Equal(now.Add(30 * time.Minute)) {
		t.Fatalf("unexpected reset expiry: %s", repo.resetExpires)
	}
}

func TestPasswordResetHashesPasswordAndToken(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, time.Hour)
	now := time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return now }
	if err := service.ResetPassword(context.Background(), "opaque-reset-token", "a-new-safe-password"); err != nil {
		t.Fatal(err)
	}
	if string(repo.resetHash) == "opaque-reset-token" || !repo.resetAt.Equal(now) {
		t.Fatal("reset token was not hashed or reset time was not recorded")
	}
	if bcrypt.CompareHashAndPassword([]byte(repo.passwordHash), []byte("a-new-safe-password")) != nil {
		t.Fatal("new password was not bcrypt hashed")
	}
}
