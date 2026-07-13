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

type HabitService interface {
	Create(ctx context.Context, userID string, req domain.CreateHabitRequest) (*domain.Habit, error)
	List(ctx context.Context, userID string, categoryID *string, archived *bool) ([]domain.Habit, error)
	GetByID(ctx context.Context, userID string, habitID string) (*domain.Habit, error)
	Update(ctx context.Context, userID string, habitID string, req domain.UpdateHabitRequest) (*domain.Habit, error)
	Delete(ctx context.Context, userID string, habitID string) error
	Archive(ctx context.Context, userID string, habitID string, archived bool) (*domain.Habit, error)
	Reorder(ctx context.Context, userID string, ids []string) ([]domain.Habit, error)
}

type BoardService interface {
	GetBoard(ctx context.Context, userID string, date string, userTZ *time.Location) (*domain.Board, error)
}

type EntryService interface {
	Create(ctx context.Context, userID string, req domain.CreateEntryRequest) (*domain.Entry, error)
	Update(ctx context.Context, userID string, entryID string, req domain.UpdateEntryRequest) (*domain.Entry, error)
	Delete(ctx context.Context, userID string, entryID string) (*domain.Entry, error)
}

type CategoryService interface {
	Create(ctx context.Context, userID string, req domain.CreateCategoryRequest) (*domain.Category, error)
	List(ctx context.Context, userID string) ([]domain.Category, error)
	GetByID(ctx context.Context, userID string, categoryID string) (*domain.Category, error)
	Update(ctx context.Context, userID string, categoryID string, req domain.UpdateCategoryRequest) (*domain.Category, error)
	Delete(ctx context.Context, userID string, categoryID string) error
	Reorder(ctx context.Context, userID string, ids []string) ([]domain.Category, error)
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
	healthService   *service.HealthService
	version         string
	authService     AuthService
	habitService    HabitService
	boardService    BoardService
	entryService    EntryService
	categoryService CategoryService
	dbChecker       DBChecker
}

func NewHandler(healthService *service.HealthService, version string) *Handler {
	return &Handler{healthService: healthService, version: version}
}

func NewHandlerWithAuth(healthService *service.HealthService, version string, authService AuthService) *Handler {
	return &Handler{healthService: healthService, version: version, authService: authService}
}

func NewHandlerWithHabit(healthService *service.HealthService, version string, authService AuthService, habitService HabitService) *Handler {
	return &Handler{healthService: healthService, version: version, authService: authService, habitService: habitService}
}

func NewHandlerWithCategory(healthService *service.HealthService, version string, authService AuthService, categoryService CategoryService) *Handler {
	return &Handler{healthService: healthService, version: version, authService: authService, categoryService: categoryService}
}

func NewHandlerWithServices(healthService *service.HealthService, version string, authService AuthService, habitService HabitService, categoryService CategoryService, boardService BoardService, entryService EntryService) *Handler {
	return &Handler{healthService: healthService, version: version, authService: authService, habitService: habitService, categoryService: categoryService, boardService: boardService, entryService: entryService}
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
	apiGroup.GET("/habits", h.ListHabits)
	apiGroup.POST("/habits", h.CreateHabit)
	apiGroup.GET("/habits/:id", h.GetHabit)
	apiGroup.GET("/board/:date", h.GetBoard)
	apiGroup.POST("/entries", h.CreateEntry)
	apiGroup.PUT("/entries/:id", h.UpdateEntry)
	apiGroup.DELETE("/entries/:id", h.DeleteEntry)
	apiGroup.PUT("/habits/:id", h.UpdateHabit)
	apiGroup.DELETE("/habits/:id", h.DeleteHabit)
	apiGroup.PATCH("/habits/:id/archive", h.ArchiveHabit)
	apiGroup.PATCH("/habits/reorder", h.ReorderHabits)
	apiGroup.GET("/categories", h.ListCategories)
	apiGroup.POST("/categories", h.CreateCategory)
	apiGroup.GET("/categories/:id", h.GetCategory)
	apiGroup.PUT("/categories/:id", h.UpdateCategory)
	apiGroup.DELETE("/categories/:id", h.DeleteCategory)
	apiGroup.PATCH("/categories/reorder", h.ReorderCategories)
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

func (h *Handler) CreateHabit(c echo.Context) error {
	if h.habitService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "habit service unavailable"})
	}

	var req domain.CreateHabitRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	habit, err := h.habitService.Create(c.Request().Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrHabitValidation):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid habit payload"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create habit"})
		}
	}
	return c.JSON(http.StatusCreated, habit)
}

func (h *Handler) ListHabits(c echo.Context) error {
	if h.habitService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "habit service unavailable"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	var categoryID *string
	if value := c.QueryParam("category"); strings.TrimSpace(value) != "" {
		categoryID = &value
	}
	var archived *bool
	if value := c.QueryParam("archived"); strings.TrimSpace(value) != "" {
		parsed := value == "true"
		archived = &parsed
	}

	habits, err := h.habitService.List(c.Request().Context(), userID, categoryID, archived)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not list habits"})
	}
	return c.JSON(http.StatusOK, habits)
}

func (h *Handler) GetHabit(c echo.Context) error {
	if h.habitService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "habit service unavailable"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	habit, err := h.habitService.GetByID(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrHabitNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "habit not found"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not fetch habit"})
		}
	}
	return c.JSON(http.StatusOK, habit)
}

func (h *Handler) UpdateHabit(c echo.Context) error {
	if h.habitService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "habit service unavailable"})
	}

	var req domain.UpdateHabitRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	habit, err := h.habitService.Update(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrHabitValidation):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid habit payload"})
		case errors.Is(err, service.ErrHabitNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "habit not found"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not update habit"})
		}
	}
	return c.JSON(http.StatusOK, habit)
}

func (h *Handler) DeleteHabit(c echo.Context) error {
	if h.habitService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "habit service unavailable"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	err := h.habitService.Delete(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrHabitNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "habit not found"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not delete habit"})
		}
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) ArchiveHabit(c echo.Context) error {
	if h.habitService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "habit service unavailable"})
	}

	var req domain.ArchiveHabitRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	habit, err := h.habitService.Archive(c.Request().Context(), userID, c.Param("id"), req.Archived)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrHabitValidation):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid habit payload"})
		case errors.Is(err, service.ErrHabitNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "habit not found"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not archive habit"})
		}
	}
	return c.JSON(http.StatusOK, habit)
}

func (h *Handler) GetBoard(c echo.Context) error {
	if h.boardService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "board service unavailable"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	date := c.Param("date")
	if strings.TrimSpace(date) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid date"})
	}

	userTZ := time.Now().Location()
	board, err := h.boardService.GetBoard(c.Request().Context(), userID, date, userTZ)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBoardFutureDate):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "future date is not allowed"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not load board"})
		}
	}
	return c.JSON(http.StatusOK, board)
}

func (h *Handler) CreateEntry(c echo.Context) error {
	if h.entryService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "entry service unavailable"})
	}

	var req domain.CreateEntryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	entry, err := h.entryService.Create(c.Request().Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEntryForbidden):
			return c.JSON(http.StatusForbidden, map[string]string{"error": "entry date is out of edit window"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create entry"})
		}
	}
	return c.JSON(http.StatusCreated, entry)
}

func (h *Handler) UpdateEntry(c echo.Context) error {
	if h.entryService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "entry service unavailable"})
	}

	var req domain.UpdateEntryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	entry, err := h.entryService.Update(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not update entry"})
	}
	if entry == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "entry not found"})
	}
	return c.JSON(http.StatusOK, entry)
}

func (h *Handler) DeleteEntry(c echo.Context) error {
	if h.entryService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "entry service unavailable"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	entry, err := h.entryService.Delete(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not delete entry"})
	}
	if entry == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "entry not found"})
	}
	return c.JSON(http.StatusOK, entry)
}

func (h *Handler) ReorderHabits(c echo.Context) error {
	if h.habitService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "habit service unavailable"})
	}

	var req domain.ReorderHabitsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	habits, err := h.habitService.Reorder(c.Request().Context(), userID, req.IDs)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrHabitValidation):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid habit payload"})
		case errors.Is(err, service.ErrHabitForbidden):
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not reorder habits"})
		}
	}
	return c.JSON(http.StatusOK, habits)
}

func (h *Handler) CreateCategory(c echo.Context) error {
	if h.categoryService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "category service unavailable"})
	}

	var req domain.CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	category, err := h.categoryService.Create(c.Request().Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCategoryValidation):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category payload"})
		case errors.Is(err, service.ErrCategoryConflict):
			return c.JSON(http.StatusConflict, map[string]string{"error": "category already exists"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create category"})
		}
	}
	return c.JSON(http.StatusCreated, category)
}

func (h *Handler) ListCategories(c echo.Context) error {
	if h.categoryService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "category service unavailable"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	categories, err := h.categoryService.List(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not list categories"})
	}
	return c.JSON(http.StatusOK, categories)
}

func (h *Handler) GetCategory(c echo.Context) error {
	if h.categoryService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "category service unavailable"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	category, err := h.categoryService.GetByID(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCategoryNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "category not found"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not fetch category"})
		}
	}
	return c.JSON(http.StatusOK, category)
}

func (h *Handler) UpdateCategory(c echo.Context) error {
	if h.categoryService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "category service unavailable"})
	}

	var req domain.UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	category, err := h.categoryService.Update(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCategoryValidation):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category payload"})
		case errors.Is(err, service.ErrCategoryNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "category not found"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not update category"})
		}
	}
	return c.JSON(http.StatusOK, category)
}

func (h *Handler) DeleteCategory(c echo.Context) error {
	if h.categoryService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "category service unavailable"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	err := h.categoryService.Delete(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCategoryNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "category not found"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not delete category"})
		}
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) ReorderCategories(c echo.Context) error {
	if h.categoryService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "category service unavailable"})
	}

	var req domain.ReorderCategoriesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	userID, ok := c.Get(ContextUserID).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	categories, err := h.categoryService.Reorder(c.Request().Context(), userID, req.IDs)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCategoryValidation):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category payload"})
		case errors.Is(err, service.ErrCategoryForbidden):
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not reorder categories"})
		}
	}
	return c.JSON(http.StatusOK, categories)
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
