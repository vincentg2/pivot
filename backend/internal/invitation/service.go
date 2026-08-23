package invitation

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) Create(ctx context.Context, label string, expiresAt *time.Time, maxUses int, createdBy uuid.UUID) (Invitation, string, error) {
	buffer := make([]byte, 10)
	if _, err := rand.Read(buffer); err != nil {
		return Invitation{}, "", err
	}
	code := "PIVOT-" + base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buffer)
	hash := sha256.Sum256([]byte(strings.ToUpper(code)))
	item, err := s.repo.Create(ctx, Invitation{ID: uuid.New(), Label: strings.TrimSpace(label), ExpiresAt: expiresAt, MaxUses: maxUses}, hash[:], createdBy)
	return item, code, err
}

func (s *Service) List(ctx context.Context) ([]Invitation, error) { return s.repo.List(ctx) }
func (s *Service) Revoke(ctx context.Context, id uuid.UUID) error { return s.repo.Revoke(ctx, id) }
