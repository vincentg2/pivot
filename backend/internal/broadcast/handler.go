package broadcast

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

func (h *Handler) List(c echo.Context) error      { return h.list(c, false) }
func (h *Handler) AdminList(c echo.Context) error { return h.list(c, true) }
func (h *Handler) list(c echo.Context, includeHidden bool) error {
	location, _ := time.LoadLocation("Europe/Paris")
	today := time.Now().In(location)
	from, err := parseBroadcastDate(c.QueryParam("from"), today)
	if err != nil {
		return invalidBroadcastDates()
	}
	to, err := parseBroadcastDate(c.QueryParam("to"), from.AddDate(0, 0, 7))
	if err != nil {
		return invalidBroadcastDates()
	}
	items, err := h.service.List(c.Request().Context(), from, to, includeHidden)
	if errors.Is(err, ErrInvalidWindow) {
		return invalidBroadcastDates()
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"listings": items, "from": from.Format("2006-01-02"), "to": to.Format("2006-01-02"), "attribution": "French TV schedule provided by Footao"})
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
		return httpx.NewProblem(409, "Connector disabled", "Set FOOTAO_ENABLED=true and configure an identifiable FOOTAO_USER_AGENT before synchronizing.")
	}
	if err != nil {
		h.logger.Error("Footao synchronization failed", "error", err)
		return httpx.NewProblem(502, "Synchronization failed", "TV schedule data could not be synchronized. Try again later.")
	}
	return c.JSON(http.StatusOK, map[string]any{"run": run})
}
func (h *Handler) Correct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return invalidListing()
	}
	var input CorrectionInput
	if err = c.Bind(&input); err != nil || c.Validate(input) != nil {
		return httpx.NewProblem(400, "Invalid correction", "Check the date, label, kind, channels, and note.")
	}
	user, _ := auth.UserFromContext(c)
	if err = h.service.Correct(c.Request().Context(), id, user.ID, input); errors.Is(err, ErrListingMissing) {
		return invalidListing()
	} else if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
func (h *Handler) Clear(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return invalidListing()
	}
	user, _ := auth.UserFromContext(c)
	if err = h.service.ClearCorrection(c.Request().Context(), id, user.ID); errors.Is(err, ErrListingMissing) {
		return invalidListing()
	} else if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
func (h *Handler) Audit(c echo.Context) error {
	items, err := h.service.Audit(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"audit": items})
}
func parseBroadcastDate(value string, fallback time.Time) (time.Time, error) {
	if value == "" {
		return time.Date(fallback.Year(), fallback.Month(), fallback.Day(), 0, 0, 0, 0, fallback.Location()), nil
	}
	return time.Parse("2006-01-02", value)
}
func invalidBroadcastDates() error {
	return httpx.NewProblem(400, "Invalid date window", "Use YYYY-MM-DD dates spanning no more than 31 days.")
}
func invalidListing() error {
	return httpx.NewProblem(404, "TV listing unavailable", "The requested TV listing does not exist.")
}
