package user

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }
func (s *Service) UpdateProfile(ctx context.Context, id uuid.UUID, nickname, theme string) error {
	return s.repo.UpdateProfile(ctx, id, strings.TrimSpace(nickname), theme)
}
func (s *Service) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteAccount(ctx, id)
}
