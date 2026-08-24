package football

import (
	"errors"
	"log/slog"
	"net/http"
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

func (h *Handler) Matches(c echo.Context) error {
	location, _ := time.LoadLocation("Europe/Paris")
	today := time.Now().In(location)
	from, err := parseDate(c.QueryParam("from"), today)
	if err != nil {
		return invalidDates()
	}
	to, err := parseDate(c.QueryParam("to"), from.AddDate(0, 0, 7))
	if err != nil {
		return invalidDates()
	}
	var clubID *uuid.UUID
	if raw := c.QueryParam("club"); raw != "" {
		parsed, parseErr := uuid.Parse(raw)
		if parseErr != nil {
			return httpx.NewProblem(400, "Invalid club filter", "Use a valid club identifier.")
		}
		clubID = &parsed
	}
	current, _ := auth.UserFromContext(c)
	items, err := h.service.ListMatches(c.Request().Context(), current.ID, from, to, c.QueryParam("competition"), clubID)
	if errors.Is(err, ErrInvalidWindow) {
		return invalidDates()
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"matches": items, "from": from.Format("2006-01-02"), "to": to.Format("2006-01-02"), "attribution": "Football data provided by football-data.org"})
}

func (h *Handler) Standing(c echo.Context) error {
	if c.QueryParam("competition") == "" {
		return httpx.NewProblem(400, "Competition required", "Choose a competition to view its table.")
	}
	item, err := h.service.Standing(c.Request().Context(), c.QueryParam("competition"))
	if errors.Is(err, ErrCompetitionMissing) {
		return httpx.NewProblem(404, "Standing unavailable", "No current standing has been synchronized for this competition.")
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"standing": item, "attribution": "Football data provided by football-data.org"})
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
	if errors.Is(err, ErrConnectorDisabled) {
		return httpx.NewProblem(409, "Connector disabled", "Set FOOTBALL_DATA_API_KEY before synchronizing.")
	}
	if errors.Is(err, ErrCompetitionMissing) {
		return httpx.NewProblem(409, "Catalog required", "Synchronize the club catalog before sports data.")
	}
	if err != nil {
		h.logger.Error("sports synchronization failed", "error", err)
		return httpx.NewProblem(502, "Synchronization failed", "Match and standing data could not be synchronized. Try again later.")
	}
	return c.JSON(http.StatusOK, map[string]any{"run": run})
}

func parseDate(value string, fallback time.Time) (time.Time, error) {
	if value == "" {
		return time.Date(fallback.Year(), fallback.Month(), fallback.Day(), 0, 0, 0, 0, fallback.Location()), nil
	}
	return time.Parse("2006-01-02", value)
}
func invalidDates() error {
	return httpx.NewProblem(400, "Invalid date window", "Use YYYY-MM-DD dates spanning no more than 31 days.")
}
