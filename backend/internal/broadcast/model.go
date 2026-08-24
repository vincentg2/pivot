package broadcast

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrConnectorDisabled = errors.New("Footao connector is disabled")
	ErrInvalidWindow     = errors.New("invalid broadcast date window")
	ErrListingMissing    = errors.New("broadcast listing not found")
)

type Listing struct {
	ID              uuid.UUID  `json:"id"`
	MatchID         *uuid.UUID `json:"matchId"`
	StartsAt        time.Time  `json:"startsAt"`
	HomeName        string     `json:"homeName"`
	AwayName        string     `json:"awayName"`
	Label           string     `json:"label"`
	CompetitionName string     `json:"competitionName"`
	Kind            string     `json:"kind"`
	Channels        []string   `json:"channels"`
	SourceURL       *string    `json:"sourceUrl"`
	External        bool       `json:"external"`
	Corrected       bool       `json:"corrected"`
	Hidden          bool       `json:"hidden,omitempty"`
}

type ImportedListing struct {
	SourceKey       string
	StartsAt        time.Time
	HomeName        string
	AwayName        string
	Label           string
	CompetitionName string
	Kind            string
	Channels        []string
	SourceURL       *string
}

type MatchCandidate struct {
	ID       uuid.UUID
	HomeName string
	AwayName string
}

type CorrectionInput struct {
	StartsAt        time.Time `json:"startsAt" validate:"required"`
	HomeName        string    `json:"homeName" validate:"max=140"`
	AwayName        string    `json:"awayName" validate:"max=140"`
	Label           string    `json:"label" validate:"required,max=280"`
	CompetitionName string    `json:"competitionName" validate:"max=140"`
	Kind            string    `json:"kind" validate:"required,oneof=live delayed replay"`
	Channels        []string  `json:"channels" validate:"max=12,dive,required,max=100"`
	Hidden          bool      `json:"hidden"`
	Note            string    `json:"note" validate:"max=300"`
}

type AuditEntry struct {
	ID        uuid.UUID      `json:"id"`
	ListingID uuid.UUID      `json:"listingId"`
	AdminID   *uuid.UUID     `json:"adminId"`
	Action    string         `json:"action"`
	Before    map[string]any `json:"before,omitempty"`
	After     map[string]any `json:"after,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
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
	List(context.Context, time.Time, time.Time, bool) ([]Listing, error)
	FindMatchCandidates(context.Context, time.Time) ([]MatchCandidate, error)
	Upsert(context.Context, ImportedListing, *uuid.UUID) error
	CleanupPast(context.Context, time.Time) (int64, error)
	Correct(context.Context, uuid.UUID, uuid.UUID, CorrectionInput) error
	ClearCorrection(context.Context, uuid.UUID, uuid.UUID) error
	Audit(context.Context, int) ([]AuditEntry, error)
	StartCollection(context.Context) (CollectionRun, error)
	FinishCollection(context.Context, uuid.UUID, string, int, string) error
	LatestCollection(context.Context) (*CollectionRun, error)
}

type Connector interface {
	Enabled() bool
	Fetch(context.Context, time.Time, time.Time) ([]ImportedListing, error)
}
