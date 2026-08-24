package broadcast

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo      Repository
	connector Connector
}

func NewService(repo Repository, connector Connector) *Service {
	return &Service{repo: repo, connector: connector}
}

func (s *Service) List(ctx context.Context, from, to time.Time, includeHidden bool) ([]Listing, error) {
	if to.Before(from) || to.Sub(from) > 31*24*time.Hour {
		return nil, ErrInvalidWindow
	}
	return s.repo.List(ctx, from, to, includeHidden)
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
	location, _ := time.LoadLocation("Europe/Paris")
	localNow := now.In(location)
	from := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, location)
	to := from.AddDate(0, 0, 7)
	if _, err = s.repo.CleanupPast(ctx, now); err != nil {
		return CollectionRun{}, s.fail(ctx, run.ID, 0, err)
	}
	items, err := s.connector.Fetch(ctx, from, to)
	if err != nil {
		return CollectionRun{}, s.fail(ctx, run.ID, 0, err)
	}
	for index, item := range items {
		candidates, candidateErr := s.repo.FindMatchCandidates(ctx, item.StartsAt)
		if candidateErr != nil {
			return CollectionRun{}, s.fail(ctx, run.ID, index, candidateErr)
		}
		matchID := chooseMatch(item, candidates)
		if saveErr := s.repo.Upsert(ctx, item, matchID); saveErr != nil {
			return CollectionRun{}, s.fail(ctx, run.ID, index, saveErr)
		}
	}
	if err = s.repo.FinishCollection(ctx, run.ID, "succeeded", len(items), ""); err != nil {
		return CollectionRun{}, err
	}
	run.Status = "succeeded"
	run.RecordsCount = len(items)
	return run, nil
}
func (s *Service) Correct(ctx context.Context, listingID, adminID uuid.UUID, input CorrectionInput) error {
	input.Label = strings.TrimSpace(input.Label)
	input.HomeName = strings.TrimSpace(input.HomeName)
	input.AwayName = strings.TrimSpace(input.AwayName)
	input.CompetitionName = strings.TrimSpace(input.CompetitionName)
	input.Note = strings.TrimSpace(input.Note)
	channels := make([]string, 0, len(input.Channels))
	for _, channel := range input.Channels {
		channel = strings.TrimSpace(channel)
		if channel != "" && !contains(channels, channel) {
			channels = append(channels, channel)
		}
	}
	input.Channels = channels
	return s.repo.Correct(ctx, listingID, adminID, input)
}
func (s *Service) ClearCorrection(ctx context.Context, listingID, adminID uuid.UUID) error {
	return s.repo.ClearCorrection(ctx, listingID, adminID)
}
func (s *Service) Audit(ctx context.Context) ([]AuditEntry, error) { return s.repo.Audit(ctx, 100) }
func (s *Service) fail(ctx context.Context, id uuid.UUID, count int, cause error) error {
	_ = s.repo.FinishCollection(ctx, id, "failed", count, "Provider synchronization failed; consult structured server logs.")
	return cause
}

func chooseMatch(item ImportedListing, candidates []MatchCandidate) *uuid.UUID {
	if item.HomeName == "" || item.AwayName == "" {
		return nil
	}
	bestScore := 0
	var best *uuid.UUID
	for _, candidate := range candidates {
		score := nameScore(item.HomeName, candidate.HomeName) + nameScore(item.AwayName, candidate.AwayName)
		if score > bestScore {
			value := candidate.ID
			best = &value
			bestScore = score
		} else if score == bestScore && score > 0 {
			best = nil
		}
	}
	if bestScore < 4 {
		return nil
	}
	return best
}
func nameScore(left, right string) int {
	a, b := normalizeName(left), normalizeName(right)
	if a == "" || b == "" {
		return 0
	}
	if a == b {
		return 3
	}
	if strings.Contains(a, b) || strings.Contains(b, a) {
		return 2
	}
	return 0
}
