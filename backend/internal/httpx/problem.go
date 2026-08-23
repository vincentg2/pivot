package httpx

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Problem struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func (p Problem) Error() string { return p.Title }

func NewProblem(status int, title, detail string) Problem {
	return Problem{Type: "about:blank", Title: title, Status: status, Detail: detail}
}

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	problem := NewProblem(http.StatusInternalServerError, "Internal Server Error", "An unexpected error occurred.")
	var typed Problem
	if errors.As(err, &typed) {
		problem = typed
	} else {
		var httpError *echo.HTTPError
		if errors.As(err, &httpError) {
			problem = NewProblem(httpError.Code, http.StatusText(httpError.Code), "The request could not be completed.")
		}
	}
	c.Response().Header().Set(echo.HeaderContentType, "application/problem+json")
	_ = c.JSON(problem.Status, problem)
}
