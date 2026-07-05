package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/qper/hf/internal/service"
)

func TestProtectedAPIRequiresJWT(t *testing.T) {
	h := NewHandler(service.NewHealthService(), "1.0.0")
	e := echo.New()
	h.Register(e)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/habits", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing token, got %d", rec.Code)
	}
}
