package installation

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vincentg2/pivot/backend/internal/httpx"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) Status(c echo.Context) error {
	required, configured, err := h.service.Status(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"setupRequired": required, "setupTokenConfigured": configured})
}

type installRequest struct {
	Token    string `json:"token" validate:"required,max=300"`
	Email    string `json:"email" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,min=12,max=128"`
	Nickname string `json:"nickname" validate:"required,min=2,max=40"`
}

func (h *Handler) Install(c echo.Context) error {
	var request installRequest
	if err := c.Bind(&request); err != nil || c.Validate(request) != nil {
		return httpx.NewProblem(400, "Invalid setup details", "Check the setup token, email, nickname, and password.")
	}
	admin, err := h.service.Install(c.Request().Context(), request.Token, request.Email, request.Password, request.Nickname)
	if errors.Is(err, ErrAlreadyInstalled) {
		return httpx.NewProblem(409, "Already installed", "This Pivot installation already has an administrator.")
	}
	if errors.Is(err, ErrSetupDisabled) {
		return httpx.NewProblem(503, "Setup unavailable", "Configure SETUP_TOKEN on the server and restart Pivot.")
	}
	if errors.Is(err, ErrInvalidToken) {
		return httpx.NewProblem(403, "Invalid setup token", "The setup token is incorrect.")
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]any{"admin": admin})
}
