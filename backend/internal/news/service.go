package news

import (
	"context"
	"net/url"
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
func (s *Service) List(ctx context.Context, userID uuid.UUID, clubID *uuid.UUID, limit int) ([]Item, error) {
	if limit < 1 || limit > 100 {
		limit = 30
	}
	return s.repo.ListItems(ctx, userID, clubID, limit)
}
func (s *Service) Feeds(ctx context.Context) ([]Feed, error) { return s.repo.ListFeeds(ctx) }
func (s *Service) SaveFeed(ctx context.Context, adminID uuid.UUID, input FeedInput) (Feed, error) {
	parsed, err := url.Parse(strings.TrimSpace(input.URL))
	if err != nil || validateRemoteURL(parsed) != nil {
		return Feed{}, ErrUnsafeFeedURL
	}
	input.URL = parsed.String()
	input.SourceName = strings.TrimSpace(input.SourceName)
	return s.repo.SaveFeed(ctx, adminID, input)
}
func (s *Service) DeleteFeed(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteFeed(ctx, id)
}
func (s *Service) CollectionStatus(ctx context.Context) (bool, *CollectionRun, error) {
	feeds, err := s.repo.ListFeeds(ctx)
	if err != nil {
		return false, nil, err
	}
	enabled := false
	for _, feed := range feeds {
		enabled = enabled || feed.Enabled
	}
	run, err := s.repo.LatestCollection(ctx)
	return enabled, run, err
}
func (s *Service) Sync(ctx context.Context, now time.Time) (CollectionRun, error) {
	run, err := s.repo.StartCollection(ctx)
	if err != nil {
		return CollectionRun{}, err
	}
	feeds, err := s.repo.ListFeeds(ctx)
	if err != nil {
		return CollectionRun{}, s.fail(ctx, run.ID, 0, err)
	}
	total := 0
	for _, feed := range feeds {
		if !feed.Enabled {
			continue
		}
		items, fetchErr := s.connector.Fetch(ctx, feed)
		if fetchErr != nil {
			return CollectionRun{}, s.fail(ctx, run.ID, total, fetchErr)
		}
		count, saveErr := s.repo.UpsertItems(ctx, feed, items)
		if saveErr != nil {
			return CollectionRun{}, s.fail(ctx, run.ID, total, saveErr)
		}
		total += count
	}
	if _, err = s.repo.Cleanup(ctx, now); err != nil {
		return CollectionRun{}, s.fail(ctx, run.ID, total, err)
	}
	if err = s.repo.FinishCollection(ctx, run.ID, "succeeded", total, ""); err != nil {
		return CollectionRun{}, err
	}
	run.Status = "succeeded"
	run.RecordsCount = total
	return run, nil
}
func (s *Service) fail(ctx context.Context, id uuid.UUID, count int, cause error) error {
	_ = s.repo.FinishCollection(ctx, id, "failed", count, "Official feed synchronization failed; consult structured server logs.")
	return cause
}
