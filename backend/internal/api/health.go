package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/qper/hf/internal/domain"
	"github.com/qper/hf/internal/service"
)

type DBChecker interface {
	Check(ctx context.Context) error
}

type AuthService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.RegisterResponse, error)
}

type Handler struct {
	healthService *service.HealthService
	version       string
	authService   AuthService
	dbChecker     DBChecker
}

func NewHandler(healthService *service.HealthService, version string) *Handler {
	return &Handler{healthService: healthService, version: version}
}

func NewHandlerWithAuth(healthService *service.HealthService, version string, authService AuthService) *Handler {
	return &Handler{healthService: healthService, version: version, authService: authService}
}

func (h *Handler) Register(e *echo.Echo) {
	e.GET("/healthz", func(c echo.Context) error {
		h.Healthz(c.Response(), c.Request())
		return nil
	})
	e.GET("/readyz", func(c echo.Context) error {
		h.Readyz(h.dbChecker)(c.Response(), c.Request())
		return nil
	})
	e.POST("/api/v1/auth/register", h.RegisterUser)
}

func (h *Handler) RegisterUser(c echo.Context) error {
	var req domain.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	if h.authService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "auth service unavailable"})
	}

	resp, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrConflict):
			return c.JSON(http.StatusConflict, map[string]string{"error": "username or email already exists"})
		case errors.Is(err, service.ErrValidation):
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "invalid registration payload"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "registration failed"})
		}
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": h.version})
}

func (h *Handler) Readyz(db DBChecker) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "db_unavailable"})
			return
		}

		ctx := r.Context()
		if err := db.Check(ctx); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "db_unavailable"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}
