package football

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestListMatchesRejectsLongWindow(t *testing.T) {
	service := NewService(fakeRepository{}, fakeConnector{})
	from := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	if _, err := service.ListMatches(context.Background(), uuid.New(), from, from.AddDate(0, 0, 32), "", nil); err != ErrInvalidWindow {
		t.Fatalf("expected ErrInvalidWindow, got %v", err)
	}
}
func TestSyncRequiresConnector(t *testing.T) {
	service := NewService(fakeRepository{}, fakeConnector{})
	if _, err := service.Sync(context.Background(), time.Now()); err != ErrConnectorDisabled {
		t.Fatalf("expected ErrConnectorDisabled, got %v", err)
	}
}

func TestSportCollectionWindowIncludesNextMonth(t *testing.T) {
	now := time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC)
	from, to := sportCollectionWindow(now)
	if from.Format("2006-01-02") != "2026-07-25" || to.Format("2006-01-02") != "2026-09-23" {
		t.Fatalf("unexpected sports window: %s to %s", from, to)
	}
}

func TestProviderSeasonInfersCurrentFromDates(t *testing.T) {
	now := time.Now().UTC()
	season := providerSeason{ID: 1, StartDate: now.AddDate(0, -1, 0).Format("2006-01-02"), EndDate: now.AddDate(0, 1, 0).Format("2006-01-02")}.model()
	if !season.Current {
		t.Fatal("expected a season containing today to be current")
	}
}

type fakeConnector struct{}

func (fakeConnector) Enabled() bool { return false }
func (fakeConnector) FetchMatches(context.Context, string, time.Time, time.Time) (Season, []Match, error) {
	return Season{}, nil, nil
}
func (fakeConnector) FetchStanding(context.Context, string) (Season, []StandingRow, error) {
	return Season{}, nil, nil
}

type fakeRepository struct{}

func (fakeRepository) ListMatches(context.Context, uuid.UUID, time.Time, time.Time, string, *uuid.UUID) ([]Match, error) {
	return nil, nil
}
func (fakeRepository) CurrentStanding(context.Context, string) (Standing, error) {
	return Standing{}, nil
}
func (fakeRepository) SaveCompetitionData(context.Context, string, Season, []Match, []StandingRow) (int, error) {
	return 0, nil
}
func (fakeRepository) StartCollection(context.Context) (CollectionRun, error) {
	return CollectionRun{}, nil
}
func (fakeRepository) FinishCollection(context.Context, uuid.UUID, string, int, string) error {
	return nil
}
func (fakeRepository) LatestCollection(context.Context) (*CollectionRun, error) { return nil, nil }
