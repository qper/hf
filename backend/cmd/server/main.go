package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/qper/hf/internal/api"
	"github.com/qper/hf/internal/config"
	"github.com/qper/hf/internal/repository"
	"github.com/qper/hf/internal/service"
)

var version = "dev"

func main() {
	cfg := config.Load()
	cfg.Version = version

	srv := newServer(cfg)
	go func() {
		if err := srv.Start(cfg.Addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server start failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

func newServer(cfg ...config.Config) *echo.Echo {
	var appConfig config.Config
	if len(cfg) > 0 {
		appConfig = cfg[0]
	} else {
		appConfig = config.Load()
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogMethod: true,
		LogURI:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Printf("method=%s uri=%s status=%d", v.Method, v.URI, v.Status)
			return nil
		},
	}))
	e.Use(middleware.Recover())

	healthService := service.NewHealthService()
	repo := repository.NewRepository()
	_ = repo
	apiHandler := api.NewHandler(healthService, appConfig.Version)
	apiHandler.Register(e)

	fmt.Printf("listening on %s\n", appConfig.Addr)
	return e
}
