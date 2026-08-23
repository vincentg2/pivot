package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vincentg2/pivot/backend/internal/auth"
	"github.com/vincentg2/pivot/backend/internal/httpx"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

type updateRequest struct {
	Nickname string `json:"nickname" validate:"required,min=2,max=40"`
	Theme    string `json:"theme" validate:"required,oneof=system light dark"`
}

func (h *Handler) Update(c echo.Context) error {
	var request updateRequest
	if err := c.Bind(&request); err != nil || c.Validate(request) != nil {
		return httpx.NewProblem(400, "Invalid profile", "Nickname and theme are invalid.")
	}
	current, _ := auth.UserFromContext(c)
	if err := h.service.UpdateProfile(c.Request().Context(), current.ID, request.Nickname, request.Theme); err != nil {
		return err
	}
	current.Nickname, current.Theme = request.Nickname, request.Theme
	return c.JSON(http.StatusOK, map[string]any{"user": current})
}

func (h *Handler) Delete(c echo.Context) error {
	current, _ := auth.UserFromContext(c)
	if err := h.service.DeleteAccount(c.Request().Context(), current.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
