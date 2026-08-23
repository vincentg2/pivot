package invitation

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/vincentg2/pivot/backend/internal/auth"
	"github.com/vincentg2/pivot/backend/internal/httpx"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

type createRequest struct {
	Label     string     `json:"label" validate:"max=100"`
	ExpiresAt *time.Time `json:"expiresAt"`
	MaxUses   int        `json:"maxUses" validate:"required,min=1,max=100"`
}

func (h *Handler) Create(c echo.Context) error {
	var request createRequest
	if err := c.Bind(&request); err != nil || c.Validate(request) != nil {
		return httpx.NewProblem(400, "Invalid invitation", "Provide a usage limit between 1 and 100 and an optional expiry.")
	}
	if request.ExpiresAt != nil && request.ExpiresAt.Before(time.Now()) {
		return httpx.NewProblem(400, "Invalid expiry", "The expiry must be in the future.")
	}
	admin, _ := auth.UserFromContext(c)
	item, code, err := h.service.Create(c.Request().Context(), request.Label, request.ExpiresAt, request.MaxUses, admin.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]any{"invitation": item, "code": code, "warning": "This code is shown only once."})
}

func (h *Handler) List(c echo.Context) error {
	items, err := h.service.List(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"invitations": items})
}

func (h *Handler) Revoke(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpx.NewProblem(400, "Invalid invitation ID", "Use a valid invitation identifier.")
	}
	err = h.service.Revoke(c.Request().Context(), id)
	if errors.Is(err, ErrNotFound) {
		return httpx.NewProblem(404, "Invitation not found", "The invitation does not exist or is already revoked.")
	}
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
