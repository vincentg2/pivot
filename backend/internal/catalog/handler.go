package catalog

import (
	"errors"
	"log/slog"
	"net/http"

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

func (h *Handler) Competitions(c echo.Context) error {
	items, err := h.service.ListCompetitions(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"competitions": items, "attribution": "Football data provided by football-data.org"})
}

func (h *Handler) Clubs(c echo.Context) error {
	items, err := h.service.ListClubs(c.Request().Context(), c.QueryParam("competition"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"clubs": items, "attribution": "Football data provided by football-data.org"})
}

func (h *Handler) Club(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpx.NewProblem(400, "Invalid club ID", "Use a valid club identifier.")
	}
	item, err := h.service.GetClub(c.Request().Context(), id)
	if errors.Is(err, ErrClubNotFound) {
		return httpx.NewProblem(404, "Club not found", "This club is not in the local catalog.")
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"club": item, "attribution": "Football data provided by football-data.org"})
}

func (h *Handler) Favorites(c echo.Context) error {
	current, _ := auth.UserFromContext(c)
	items, err := h.service.ListFavorites(c.Request().Context(), current.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"favorites": items})
}

type favoritesRequest struct {
	ClubIDs []uuid.UUID `json:"clubIds" validate:"max=5,unique"`
}

func (h *Handler) ReplaceFavorites(c echo.Context) error {
	var request favoritesRequest
	if err := c.Bind(&request); err != nil || c.Validate(request) != nil {
		return httpx.NewProblem(400, "Invalid favorites", "Choose up to five distinct clubs.")
	}
	current, _ := auth.UserFromContext(c)
	items, err := h.service.ReplaceFavorites(c.Request().Context(), current.ID, request.ClubIDs)
	if errors.Is(err, ErrTooManyFavorites) {
		return httpx.NewProblem(422, "Favorite limit reached", "You can follow up to five clubs.")
	}
	if errors.Is(err, ErrClubNotFound) {
		return httpx.NewProblem(422, "Unknown club", "One of the selected clubs is not in the catalog.")
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"favorites": items})
}

func (h *Handler) CollectionStatus(c echo.Context) error {
	enabled, latest, err := h.service.CollectionStatus(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"enabled": enabled, "latestRun": latest})
}

func (h *Handler) Sync(c echo.Context) error {
	run, err := h.service.Sync(c.Request().Context())
	if errors.Is(err, ErrConnectorDisabled) {
		return httpx.NewProblem(409, "Connector disabled", "Set FOOTBALL_DATA_API_KEY on this installation before synchronizing.")
	}
	if err != nil {
		h.logger.Error("catalog synchronization failed", "error", err)
		return httpx.NewProblem(502, "Synchronization failed", "football-data.org could not be synchronized. Try again later.")
	}
	return c.JSON(http.StatusOK, map[string]any{"run": run})
}
