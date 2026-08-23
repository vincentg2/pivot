package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vincentg2/pivot/backend/internal/httpx"
)

const ContextUserKey = "authenticated-user"

type Handler struct {
	service    *Service
	cookieName string
	secure     bool
	ttl        time.Duration
}

func NewHandler(service *Service, cookieName string, secure bool, ttl time.Duration) *Handler {
	return &Handler{service: service, cookieName: cookieName, secure: secure, ttl: ttl}
}

type registerRequest struct {
	Email          string `json:"email" validate:"required,email,max=254"`
	Password       string `json:"password" validate:"required,min=12,max=128"`
	Nickname       string `json:"nickname" validate:"required,min=2,max=40"`
	InvitationCode string `json:"invitationCode" validate:"required,min=8,max=100"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,max=128"`
}

func (h *Handler) Register(c echo.Context) error {
	var request registerRequest
	if err := c.Bind(&request); err != nil || c.Validate(request) != nil {
		return httpx.NewProblem(400, "Invalid registration", "Check the email, nickname, password, and invitation code.")
	}
	user, err := h.service.Register(c.Request().Context(), request.Email, request.Password, request.Nickname, request.InvitationCode)
	if errors.Is(err, ErrInvitationInvalid) {
		return httpx.NewProblem(422, "Invalid invitation", "This invitation is invalid, expired, revoked, or fully used.")
	}
	if errors.Is(err, ErrEmailTaken) {
		return httpx.NewProblem(409, "Email already registered", "Sign in or ask an administrator for help.")
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]any{"user": user})
}

func (h *Handler) Login(c echo.Context) error {
	var request loginRequest
	if err := c.Bind(&request); err != nil || c.Validate(request) != nil {
		return httpx.NewProblem(400, "Invalid login", "Enter a valid email and password.")
	}
	user, token, err := h.service.Login(c.Request().Context(), request.Email, request.Password, c.Request().UserAgent(), c.RealIP())
	if errors.Is(err, ErrInvalidCredentials) {
		return httpx.NewProblem(401, "Invalid credentials", "The email or password is incorrect.")
	}
	if err != nil {
		return err
	}
	h.setCookie(c, token, int(h.ttl.Seconds()))
	return c.JSON(http.StatusOK, map[string]any{"user": user})
}

func (h *Handler) Logout(c echo.Context) error {
	cookie, _ := c.Cookie(h.cookieName)
	if cookie != nil {
		_ = h.service.Logout(c.Request().Context(), cookie.Value)
	}
	h.setCookie(c, "", -1)
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Me(c echo.Context) error {
	user, ok := UserFromContext(c)
	if !ok {
		return httpx.NewProblem(401, "Authentication required", "Sign in to continue.")
	}
	return c.JSON(http.StatusOK, map[string]any{"user": user})
}

func (h *Handler) RequireSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(h.cookieName)
		if err != nil {
			return httpx.NewProblem(401, "Authentication required", "Sign in to continue.")
		}
		user, err := h.service.Authenticate(c.Request().Context(), cookie.Value)
		if err != nil {
			return httpx.NewProblem(401, "Authentication required", "Your session is invalid or expired.")
		}
		c.Set(ContextUserKey, user)
		return next(c)
	}
}

func (h *Handler) RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return h.RequireSession(func(c echo.Context) error {
		user, _ := UserFromContext(c)
		if user.Role != "admin" {
			return httpx.NewProblem(403, "Administrator required", "You do not have permission to perform this action.")
		}
		return next(c)
	})
}

func UserFromContext(c echo.Context) (User, bool) {
	user, ok := c.Get(ContextUserKey).(User)
	return user, ok
}

func (h *Handler) setCookie(c echo.Context, value string, maxAge int) {
	c.SetCookie(&http.Cookie{Name: h.cookieName, Value: value, Path: "/", HttpOnly: true, Secure: h.secure, SameSite: http.SameSiteLaxMode, MaxAge: maxAge})
}
