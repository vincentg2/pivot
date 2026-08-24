package news

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/vincentg2/pivot/backend/internal/auth"
	"github.com/vincentg2/pivot/backend/internal/httpx"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}
func (h *Handler) List(c echo.Context) error {
	user, _ := auth.UserFromContext(c)
	var clubID *uuid.UUID
	if raw := c.QueryParam("club"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			return httpx.NewProblem(400, "Invalid club filter", "Use a valid club identifier.")
		}
		clubID = &id
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	items, err := h.service.List(c.Request().Context(), user.ID, clubID, limit)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"news": items, "retentionDays": 30, "contentPolicy": "Titles and links from official club feeds only"})
}
func (h *Handler) Feeds(c echo.Context) error {
	items, err := h.service.Feeds(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"feeds": items})
}
func (h *Handler) SaveFeed(c echo.Context) error {
	var input FeedInput
	if err := c.Bind(&input); err != nil || c.Validate(input) != nil {
		return httpx.NewProblem(400, "Invalid official feed", "Choose a club, source name, and valid public feed URL.")
	}
	user, _ := auth.UserFromContext(c)
	item, err := h.service.SaveFeed(c.Request().Context(), user.ID, input)
	if errors.Is(err, ErrUnsafeFeedURL) {
		return httpx.NewProblem(422, "Unsafe feed URL", "Only public HTTP or HTTPS feed URLs can be fetched.")
	}
	if errors.Is(err, ErrClubMissing) {
		return httpx.NewProblem(404, "Club not found", "The selected club is not in the local catalog.")
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]any{"feed": item})
}
func (h *Handler) DeleteFeed(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return feedMissing()
	}
	if err = h.service.DeleteFeed(c.Request().Context(), id); errors.Is(err, ErrFeedMissing) {
		return feedMissing()
	} else if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
func (h *Handler) CollectionStatus(c echo.Context) error {
	enabled, latest, err := h.service.CollectionStatus(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"enabled": enabled, "latestRun": latest})
}
func (h *Handler) Sync(c echo.Context) error {
	run, err := h.service.Sync(c.Request().Context(), time.Now())
	if err != nil {
		h.logger.Error("official news synchronization failed", "error", err)
		return httpx.NewProblem(502, "Synchronization failed", "Official club feeds could not be synchronized. Check their configuration and try again.")
	}
	return c.JSON(http.StatusOK, map[string]any{"run": run})
}
func feedMissing() error {
	return httpx.NewProblem(404, "Feed not found", "The requested official feed does not exist.")
}
