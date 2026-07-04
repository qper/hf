package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/qper/hf/internal/domain"
	"github.com/qper/hf/internal/service"
)

type stubAuthService struct {
	resp *domain.RegisterResponse
	err  error
}

func (s stubAuthService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.RegisterResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.resp, nil
}

func TestRegisterEndpointReturnsCreated(t *testing.T) {
	h := NewHandler(service.NewHealthService(), "1.0.0")
	h.authService = stubAuthService{resp: &domain.RegisterResponse{User: domain.User{ID: "u1", Username: "alice", Email: "alice@example.com"}, RecoveryCodes: []string{"AAAAAA", "BBBBBB"}}}

	e := echo.New()
	h.Register(e)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(`{"username":"alice","email":"alice@example.com","password":"StrongPass1"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestRegisterEndpointReturnsConflict(t *testing.T) {
	h := NewHandler(service.NewHealthService(), "1.0.0")
	h.authService = stubAuthService{err: service.ErrConflict}

	e := echo.New()
	h.Register(e)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(`{"username":"alice","email":"alice@example.com","password":"StrongPass1"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
}

func TestRegisterEndpointReturnsUnprocessableEntity(t *testing.T) {
	h := NewHandler(service.NewHealthService(), "1.0.0")
	h.authService = stubAuthService{err: service.ErrValidation}

	e := echo.New()
	h.Register(e)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(`{"username":"al","email":"alice@example.com","password":"short"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, rec.Code)
	}
}
