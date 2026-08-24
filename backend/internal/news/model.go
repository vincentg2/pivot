package news

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrFeedMissing   = errors.New("news feed not found")
	ErrClubMissing   = errors.New("club not found")
	ErrUnsafeFeedURL = errors.New("feed URL is not safe for server-side fetching")
)

type Feed struct {
	ID         uuid.UUID `json:"id"`
	ClubID     uuid.UUID `json:"clubId"`
	ClubName   string    `json:"clubName"`
	URL        string    `json:"url"`
	SourceName string    `json:"sourceName"`
	Enabled    bool      `json:"enabled"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type FeedInput struct {
	ClubID     uuid.UUID `json:"clubId" validate:"required"`
	URL        string    `json:"url" validate:"required,url,max=2000"`
	SourceName string    `json:"sourceName" validate:"required,max=120"`
	Enabled    bool      `json:"enabled"`
}

type Item struct {
	ID          uuid.UUID `json:"id"`
	ClubID      uuid.UUID `json:"clubId"`
	ClubName    string    `json:"clubName"`
	ClubTLA     string    `json:"clubTla"`
	Title       string    `json:"title"`
	SourceName  string    `json:"sourceName"`
	PublishedAt time.Time `json:"publishedAt"`
	LinkURL     string    `json:"linkUrl"`
}

type ImportedItem struct {
	ExternalID  string
	Title       string
	PublishedAt time.Time
	LinkURL     string
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
	ListItems(context.Context, uuid.UUID, *uuid.UUID, int) ([]Item, error)
	ListFeeds(context.Context) ([]Feed, error)
	SaveFeed(context.Context, uuid.UUID, FeedInput) (Feed, error)
	DeleteFeed(context.Context, uuid.UUID) error
	UpsertItems(context.Context, Feed, []ImportedItem) (int, error)
	Cleanup(context.Context, time.Time) (int64, error)
	StartCollection(context.Context) (CollectionRun, error)
	FinishCollection(context.Context, uuid.UUID, string, int, string) error
	LatestCollection(context.Context) (*CollectionRun, error)
}

type Connector interface {
	Fetch(context.Context, Feed) ([]ImportedItem, error)
}
