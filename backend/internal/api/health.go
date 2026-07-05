package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/qper/hf/internal/domain"
	"github.com/qper/hf/internal/service"
)

type DBChecker interface {
	Check(ctx context.Context) error
}

type AuthService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.RegisterResponse, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*domain.RefreshResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, refreshToken string) error
	RecoverWithRecoveryCode(ctx context.Context, req domain.RecoverRequest) (*domain.RecoverResponse, error)
	GetRecoveryCodeCount(ctx context.Context, userID string) (int, error)
	RegenerateRecoveryCodes(ctx context.Context, userID, password string) (*domain.RecoveryCodeRegenerationResponse, error)
}

type dbChecker struct {
	db *sql.DB
}

func (d dbChecker) Check(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func NewDBChecker(db *sql.DB) DBChecker {
	return dbChecker{db: db}
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

func (h *Handler) WithDBChecker(dbChecker DBChecker) *Handler {
	h.dbChecker = dbChecker
	return h
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
	// auth routes with rate limiter, CORS and CSP
	authGroup := e.Group("/auth")
	authGroup.Use(NewRateLimiter(5, 15*time.Minute))
	authGroup.Use(CORSMiddleware)
	authGroup.Use(CSPMiddleware)
	authGroup.POST("/login", h.LoginUser)
	authGroup.POST("/refresh", h.RefreshUser)
	authGroup.POST("/logout", h.LogoutUser)
	authGroup.POST("/logout-all", h.LogoutAllUser)

	// registration and recovery are public endpoints under /api/v1/auth
	e.POST("/api/v1/auth/register", h.RegisterUser)
	e.POST("/api/v1/auth/recover", h.RecoverUser, CORSMiddleware, CSPMiddleware)

	apiGroup := e.Group("/api/v1")
	apiGroup.Use(CORSMiddleware, CSPMiddleware, JWTMiddleware())
	apiGroup.GET("/habits", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	apiGroup.GET("/me/recovery-codes", h.GetMyRecoveryCodes)
	apiGroup.POST("/me/recovery-codes", h.RegenerateMyRecoveryCodes)
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

func (h *Handler) LoginUser(c echo.Context) error {
	var req domain.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	if h.authService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "auth service unavailable"})
	}

	resp, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnauthorized):
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "login failed"})
		}
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) RefreshUser(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing refresh token"})
	}

	resp, err := h.authService.Refresh(c.Request().Context(), cookie.Value)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) LogoutUser(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		c.SetCookie(expireCookie("refresh_token"))
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
	if err := h.authService.Logout(c.Request().Context(), cookie.Value); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
	}
	c.SetCookie(expireCookie("refresh_token"))
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) LogoutAllUser(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		c.SetCookie(expireCookie("refresh_token"))
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
	if err := h.authService.LogoutAll(c.Request().Context(), cookie.Value); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
	}
	c.SetCookie(expireCookie("refresh_token"))
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) RecoverUser(c echo.Context) error {
	var req domain.RecoverRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	if h.authService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "auth service unavailable"})
	}

	resp, err := h.authService.RecoverWithRecoveryCode(c.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnauthorized):
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials or recovery code"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "recovery failed"})
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetMyRecoveryCodes(c echo.Context) error {
	userID, ok := userIDFromContext(c)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing user context"})
	}

	if h.authService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "auth service unavailable"})
	}

	remaining, err := h.authService.GetRecoveryCodeCount(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not load recovery code status"})
	}

	return c.JSON(http.StatusOK, domain.RecoveryCodeStatusResponse{Remaining: remaining})
}

func (h *Handler) RegenerateMyRecoveryCodes(c echo.Context) error {
	var req domain.RecoveryCodeRegenerationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	userID, ok := userIDFromContext(c)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing user context"})
	}

	if h.authService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "auth service unavailable"})
	}

	resp, err := h.authService.RegenerateRecoveryCodes(c.Request().Context(), userID, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnauthorized):
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid password"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not regenerate recovery codes"})
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func userIDFromContext(c echo.Context) (string, bool) {
	userID, ok := c.Get("user_id").(string)
	return userID, ok
}

func expireCookie(name string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
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
