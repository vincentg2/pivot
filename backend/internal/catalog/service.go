package catalog

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

var defaultCompetitionCodes = []string{"PL", "PD", "BL1", "SA", "FL1"}

type Service struct {
	repo      Repository
	connector Connector
}

func NewService(repo Repository, connector Connector) *Service {
	return &Service{repo: repo, connector: connector}
}

func (s *Service) ListCompetitions(ctx context.Context) ([]Competition, error) {
	return s.repo.ListCompetitions(ctx)
}
func (s *Service) ListClubs(ctx context.Context, code string) ([]Club, error) {
	return s.repo.ListClubs(ctx, strings.ToUpper(strings.TrimSpace(code)))
}
func (s *Service) GetClub(ctx context.Context, id uuid.UUID) (Club, error) {
	return s.repo.GetClub(ctx, id)
}
func (s *Service) ListFavorites(ctx context.Context, userID uuid.UUID) ([]Club, error) {
	return s.repo.ListFavorites(ctx, userID)
}

func (s *Service) ReplaceFavorites(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) ([]Club, error) {
	if len(ids) > 5 {
		return nil, ErrTooManyFavorites
	}
	seen := make(map[uuid.UUID]struct{}, len(ids))
	for _, id := range ids {
		if _, exists := seen[id]; exists {
			return nil, ErrClubNotFound
		}
		seen[id] = struct{}{}
	}
	return s.repo.ReplaceFavorites(ctx, userID, ids)
}

func (s *Service) CollectionStatus(ctx context.Context) (bool, *CollectionRun, error) {
	run, err := s.repo.LatestCollection(ctx)
	return s.connector.Enabled(), run, err
}

func (s *Service) Sync(ctx context.Context) (CollectionRun, error) {
	if !s.connector.Enabled() {
		return CollectionRun{}, ErrConnectorDisabled
	}
	run, err := s.repo.StartCollection(ctx)
	if err != nil {
		return CollectionRun{}, err
	}
	total := 0
	for _, code := range defaultCompetitionCodes {
		competition, clubs, fetchErr := s.connector.FetchCompetition(ctx, code)
		if fetchErr != nil {
			_ = s.repo.FinishCollection(ctx, run.ID, "failed", total, "Provider request failed; consult structured server logs.")
			return CollectionRun{}, fetchErr
		}
		count, saveErr := s.repo.UpsertCompetition(ctx, competition, clubs)
		if saveErr != nil {
			_ = s.repo.FinishCollection(ctx, run.ID, "failed", total, "Catalog persistence failed.")
			return CollectionRun{}, saveErr
		}
		total += count
	}
	if err := s.repo.FinishCollection(ctx, run.ID, "succeeded", total, ""); err != nil {
		return CollectionRun{}, err
	}
	run.Status, run.RecordsCount = "succeeded", total
	return run, nil
}
