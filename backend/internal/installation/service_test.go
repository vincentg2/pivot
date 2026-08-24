package installation

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	required bool
	created  int
}

func (f *fakeRepository) Required(context.Context) (bool, error) { return f.required, nil }
func (f *fakeRepository) CreateFirstAdmin(_ context.Context, input NewAdmin) (Admin, error) {
	if !f.required {
		return Admin{}, ErrAlreadyInstalled
	}
	f.required = false
	f.created++
	return Admin{ID: input.ID, Email: input.Email, Nickname: input.Nickname, Role: "admin"}, nil
}

func TestInstallCreatesOnlyFirstAdministrator(t *testing.T) {
	repo := &fakeRepository{required: true}
	service := NewService(repo, "a-long-one-time-setup-token")
	admin, err := service.Install(context.Background(), "a-long-one-time-setup-token", "OWNER@Example.test", "secure-password", " Owner ")
	if err != nil {
		t.Fatal(err)
	}
	if admin.Email != "owner@example.test" || admin.Nickname != "Owner" || repo.created != 1 {
		t.Fatalf("unexpected administrator: %#v", admin)
	}
	if _, err = service.Install(context.Background(), "a-long-one-time-setup-token", "second@example.test", "secure-password", "Second"); !errors.Is(err, ErrAlreadyInstalled) {
		t.Fatalf("expected already installed, got %v", err)
	}
}
func TestInstallRejectsInvalidOrMissingToken(t *testing.T) {
	repo := &fakeRepository{required: true}
	if _, err := NewService(repo, "").Install(context.Background(), "anything", "a@example.test", "secure-password", "Admin"); !errors.Is(err, ErrSetupDisabled) {
		t.Fatalf("expected disabled, got %v", err)
	}
	if _, err := NewService(repo, "a-long-one-time-setup-token").Install(context.Background(), "wrong", "a@example.test", "secure-password", "Admin"); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected invalid token, got %v", err)
	}
}
