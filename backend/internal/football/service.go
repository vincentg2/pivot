package football

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

var competitionCodes = []string{"PL", "PD", "BL1", "SA", "FL1"}

type Service struct {
	repo      Repository
	connector Connector
}

func NewService(repo Repository, connector Connector) *Service {
	return &Service{repo: repo, connector: connector}
}

func (s *Service) ListMatches(ctx context.Context, userID uuid.UUID, from, to time.Time, code string, clubID *uuid.UUID) ([]Match, error) {
	if to.Before(from) || to.Sub(from) > 31*24*time.Hour {
		return nil, ErrInvalidWindow
	}
	return s.repo.ListMatches(ctx, userID, from, to, strings.ToUpper(strings.TrimSpace(code)), clubID)
}
func (s *Service) Standing(ctx context.Context, code string) (Standing, error) {
	return s.repo.CurrentStanding(ctx, strings.ToUpper(strings.TrimSpace(code)))
}
func (s *Service) CollectionStatus(ctx context.Context) (bool, *CollectionRun, error) {
	run, err := s.repo.LatestCollection(ctx)
	return s.connector.Enabled(), run, err
}

func (s *Service) Sync(ctx context.Context, now time.Time) (CollectionRun, error) {
	if !s.connector.Enabled() {
		return CollectionRun{}, ErrConnectorDisabled
	}
	run, err := s.repo.StartCollection(ctx)
	if err != nil {
		return CollectionRun{}, err
	}
	from, to := sportCollectionWindow(now)
	total := 0
	for _, code := range competitionCodes {
		season, matches, fetchErr := s.connector.FetchMatches(ctx, code, from, to)
		if fetchErr != nil {
			return CollectionRun{}, s.fail(ctx, run.ID, total, fetchErr)
		}
		standingSeason, rows, standingErr := s.connector.FetchStanding(ctx, code)
		if standingErr != nil {
			return CollectionRun{}, s.fail(ctx, run.ID, total, standingErr)
		}
		if standingSeason.ProviderID != 0 {
			season = standingSeason
		}
		count, saveErr := s.repo.SaveCompetitionData(ctx, code, season, matches, rows)
		if saveErr != nil {
			return CollectionRun{}, s.fail(ctx, run.ID, total, saveErr)
		}
		total += count
	}
	if err := s.repo.FinishCollection(ctx, run.ID, "succeeded", total, ""); err != nil {
		return CollectionRun{}, err
	}
	run.Status, run.RecordsCount = "succeeded", total
	return run, nil
}

func sportCollectionWindow(now time.Time) (time.Time, time.Time) {
	return now.AddDate(0, 0, -30), now.AddDate(0, 0, 30)
}
func (s *Service) fail(ctx context.Context, id uuid.UUID, count int, cause error) error {
	_ = s.repo.FinishCollection(ctx, id, "failed", count, "Provider synchronization failed; consult structured server logs.")
	return cause
}
