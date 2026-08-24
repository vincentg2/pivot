package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestLoginRateLimiterAllowsConfiguredBurst(t *testing.T) {
	const requestsPerMinute = 10
	store := newLoginRateLimiter(requestsPerMinute)

	for request := 1; request <= requestsPerMinute; request++ {
		allowed, err := store.Allow("127.0.0.1")
		if err != nil {
			t.Fatalf("request %d: unexpected error: %v", request, err)
		}
		if !allowed {
			t.Fatalf("request %d: expected initial burst to be allowed", request)
		}
	}

	allowed, err := store.Allow("127.0.0.1")
	if err != nil {
		t.Fatalf("request after burst: unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("request after burst: expected rate limiter to reject it")
	}
}

func TestFrontendServesSPAWithoutInterceptingAPI(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "index.html"), []byte("<main>Pivot</main>"), 0o600); err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	e.GET("/api/v1/ping", func(c echo.Context) error { return c.String(http.StatusOK, "api") })
	e.GET("/health", func(c echo.Context) error { return c.String(http.StatusOK, "ok") })
	useFrontend(e, root)

	for path, expected := range map[string]string{"/clubs/example": "Pivot", "/api/v1/ping": "api", "/health": "ok"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		e.ServeHTTP(response, request)
		if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), expected) {
			t.Fatalf("%s: got status %d body %q", path, response.Code, response.Body.String())
		}
	}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/missing", nil)
	response := httptest.NewRecorder()
	e.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound || strings.Contains(response.Body.String(), "Pivot") {
		t.Fatalf("missing API route was handled by SPA: status %d body %q", response.Code, response.Body.String())
	}
}
