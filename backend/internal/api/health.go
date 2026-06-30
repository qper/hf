package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/smirnofflab/habitflow/internal/service"
)

type Handler struct {
	healthService *service.HealthService
	version       string
}

func NewHandler(healthService *service.HealthService, version string) *Handler {
	return &Handler{healthService: healthService, version: version}
}

func (h *Handler) Register(e *echo.Echo) {
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": h.healthService.Status(), "version": h.version})
	})
}
