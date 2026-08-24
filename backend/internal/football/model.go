package football

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrConnectorDisabled  = errors.New("football-data.org connector is not configured")
	ErrCompetitionMissing = errors.New("competition is not in the local catalog")
	ErrInvalidWindow      = errors.New("invalid match date window")
)

type Competition struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

type ClubRef struct {
	ID        *uuid.UUID `json:"id"`
	Name      string     `json:"name"`
	ShortName string     `json:"shortName"`
	TLA       string     `json:"tla"`
	CrestURL  *string    `json:"crestUrl"`
}

type Season struct {
	ProviderID int       `json:"-"`
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
	Current    bool      `json:"current"`
}

type Match struct {
	ID             uuid.UUID   `json:"id"`
	ProviderID     int         `json:"-"`
	Competition    Competition `json:"competition"`
	Season         Season      `json:"season"`
	UTCDate        time.Time   `json:"utcDate"`
	Status         string      `json:"status"`
	Stage          string      `json:"stage"`
	Matchday       *int        `json:"matchday"`
	Home           ClubRef     `json:"home"`
	Away           ClubRef     `json:"away"`
	HomeProviderID int         `json:"-"`
	AwayProviderID int         `json:"-"`
	HomeScore      *int        `json:"homeScore"`
	AwayScore      *int        `json:"awayScore"`
	Favorite       bool        `json:"favorite"`
}

type StandingRow struct {
	Position       int     `json:"position"`
	Club           ClubRef `json:"club"`
	ClubProviderID int     `json:"-"`
	Played         int     `json:"played"`
	Won            int     `json:"won"`
	Drawn          int     `json:"drawn"`
	Lost           int     `json:"lost"`
	GoalsFor       int     `json:"goalsFor"`
	GoalsAgainst   int     `json:"goalsAgainst"`
	GoalDifference int     `json:"goalDifference"`
	Points         int     `json:"points"`
}

type Standing struct {
	Competition Competition   `json:"competition"`
	Season      Season        `json:"season"`
	Rows        []StandingRow `json:"rows"`
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
	ListMatches(context.Context, uuid.UUID, time.Time, time.Time, string, *uuid.UUID) ([]Match, error)
	CurrentStanding(context.Context, string) (Standing, error)
	SaveCompetitionData(context.Context, string, Season, []Match, []StandingRow) (int, error)
	StartCollection(context.Context) (CollectionRun, error)
	FinishCollection(context.Context, uuid.UUID, string, int, string) error
	LatestCollection(context.Context) (*CollectionRun, error)
}

type Connector interface {
	Enabled() bool
	FetchMatches(context.Context, string, time.Time, time.Time) (Season, []Match, error)
	FetchStanding(context.Context, string) (Season, []StandingRow, error)
}
