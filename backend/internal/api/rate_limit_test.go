package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/qper/hf/internal/service"
)

func TestRateLimitEnforcedPerIP(t *testing.T) {
	h := NewHandler(service.NewHealthService(), "1.0.0")
	// stub service to always return unauthorized
	h.authService = stubAuthService{loginErr: service.ErrUnauthorized}

	e := echo.New()
	h.Register(e)

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.RemoteAddr = "1.2.3.4:1234"
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 on attempt %d, got %d", i+1, rec.Code)
		}
	}

	// 6th attempt should be 429
	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.RemoteAddr = "1.2.3.4:1234"
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 on 6th attempt, got %d", rec.Code)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Fatalf("expected Retry-After header in 429 response")
	}
}
