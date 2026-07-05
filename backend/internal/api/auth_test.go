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
	resp      *domain.RegisterResponse
	loginResp *domain.LoginResponse
	loginErr  error
	err       error
}

func (s stubAuthService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.RegisterResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.resp, nil
}

func (s stubAuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	if s.loginErr != nil {
		return nil, s.loginErr
	}
	return s.loginResp, nil
}

func (s stubAuthService) Refresh(ctx context.Context, refreshToken string) (*domain.RefreshResponse, error) {
	return nil, nil
}

func (s stubAuthService) Logout(ctx context.Context, refreshToken string) error {
	return nil
}

func (s stubAuthService) LogoutAll(ctx context.Context, refreshToken string) error {
	return nil
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

func TestLoginEndpointReturnsTokensAndCookie(t *testing.T) {
	h := NewHandler(service.NewHealthService(), "1.0.0")
	h.authService = stubAuthService{loginResp: &domain.LoginResponse{AccessToken: "token", RefreshToken: "refresh", TokenType: "Bearer", ExpiresIn: 900}}

	e := echo.New()
	h.Register(e)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"alice","password":"StrongPass1"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Header().Get("Set-Cookie") == "" {
		t.Fatalf("expected refresh token cookie to be set")
	}
	if !strings.Contains(rec.Header().Get("Set-Cookie"), "HttpOnly") || !strings.Contains(rec.Header().Get("Set-Cookie"), "Secure") || !strings.Contains(rec.Header().Get("Set-Cookie"), "SameSite=Strict") {
		t.Fatalf("expected refresh cookie attributes to be set, got %q", rec.Header().Get("Set-Cookie"))
	}
	if !strings.Contains(rec.Body.String(), "\"access_token\":\"token\"") {
		t.Fatalf("expected access token in JSON body")
	}
}
