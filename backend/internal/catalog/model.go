package catalog

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrClubNotFound      = errors.New("club not found")
	ErrTooManyFavorites  = errors.New("a maximum of five favorite clubs is allowed")
	ErrConnectorDisabled = errors.New("football-data.org connector is not configured")
)

type Competition struct {
	ID           uuid.UUID `json:"id"`
	ProviderCode string    `json:"code"`
	Name         string    `json:"name"`
	Country      string    `json:"country"`
	EmblemURL    *string   `json:"emblemUrl"`
	ClubCount    int       `json:"clubCount"`
}

type Club struct {
	ID           uuid.UUID     `json:"id"`
	ProviderID   int           `json:"-"`
	Name         string        `json:"name"`
	ShortName    string        `json:"shortName"`
	TLA          string        `json:"tla"`
	CrestURL     *string       `json:"crestUrl"`
	WebsiteURL   *string       `json:"websiteUrl"`
	Venue        string        `json:"venue"`
	Competitions []Competition `json:"competitions,omitempty"`
}

type CollectionRun struct {
	ID           uuid.UUID  `json:"id"`
	Status       string     `json:"status"`
	RecordsCount int        `json:"recordsCount"`
	ErrorMessage string     `json:"errorMessage,omitempty"`
	StartedAt    time.Time  `json:"startedAt"`
	FinishedAt   *time.Time `json:"finishedAt"`
}

type Repository interface {
	ListCompetitions(context.Context) ([]Competition, error)
	ListClubs(context.Context, string) ([]Club, error)
	GetClub(context.Context, uuid.UUID) (Club, error)
	ListFavorites(context.Context, uuid.UUID) ([]Club, error)
	ReplaceFavorites(context.Context, uuid.UUID, []uuid.UUID) ([]Club, error)
	UpsertCompetition(context.Context, Competition, []Club) (int, error)
	StartCollection(context.Context) (CollectionRun, error)
	FinishCollection(context.Context, uuid.UUID, string, int, string) error
	LatestCollection(context.Context) (*CollectionRun, error)
}

type Connector interface {
	Enabled() bool
	FetchCompetition(context.Context, string) (Competition, []Club, error)
}
