package catalog

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestReplaceFavoritesRejectsMoreThanFive(t *testing.T) {
	service := NewService(&fakeRepository{}, fakeConnector{})
	ids := []uuid.UUID{uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()}
	if _, err := service.ReplaceFavorites(context.Background(), uuid.New(), ids); err != ErrTooManyFavorites {
		t.Fatalf("expected ErrTooManyFavorites, got %v", err)
	}
}

func TestSyncRequiresConfiguredConnector(t *testing.T) {
	service := NewService(&fakeRepository{}, fakeConnector{})
	if _, err := service.Sync(context.Background()); err != ErrConnectorDisabled {
		t.Fatalf("expected ErrConnectorDisabled, got %v", err)
	}
}

type fakeConnector struct{}

func (fakeConnector) Enabled() bool { return false }
func (fakeConnector) FetchCompetition(context.Context, string) (Competition, []Club, error) {
	return Competition{}, nil, nil
}

type fakeRepository struct{}

func (*fakeRepository) ListCompetitions(context.Context) ([]Competition, error)  { return nil, nil }
func (*fakeRepository) ListClubs(context.Context, string) ([]Club, error)        { return nil, nil }
func (*fakeRepository) GetClub(context.Context, uuid.UUID) (Club, error)         { return Club{}, nil }
func (*fakeRepository) ListFavorites(context.Context, uuid.UUID) ([]Club, error) { return nil, nil }
func (*fakeRepository) ReplaceFavorites(context.Context, uuid.UUID, []uuid.UUID) ([]Club, error) {
	return nil, nil
}
func (*fakeRepository) UpsertCompetition(context.Context, Competition, []Club) (int, error) {
	return 0, nil
}
func (*fakeRepository) StartCollection(context.Context) (CollectionRun, error) {
	return CollectionRun{}, nil
}
func (*fakeRepository) FinishCollection(context.Context, uuid.UUID, string, int, string) error {
	return nil
}
func (*fakeRepository) LatestCollection(context.Context) (*CollectionRun, error) { return nil, nil }
