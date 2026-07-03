package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qper/hf/internal/service"
)

type stubDBChecker struct {
	err error
}

func (s stubDBChecker) Check(ctx context.Context) error {
	return s.err
}

func TestHealthz(t *testing.T) {
	h := NewHandler(service.NewHealthService(), "1.0.0")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	h.Healthz(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestReadyz_DBDown(t *testing.T) {
	h := NewHandler(service.NewHealthService(), "1.0.0")

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	h.Readyz(stubDBChecker{err: context.DeadlineExceeded})(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
}
