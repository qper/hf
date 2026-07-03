package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/qper/hf/internal/api"
	"github.com/qper/hf/internal/config"
	"github.com/qper/hf/internal/repository"
	"github.com/qper/hf/internal/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var version = "dev"

func main() {
	cfg := config.Load()
	cfg.Version = version

	logger := newLogger(cfg.LogLevel)
	defer func() {
		if err := logger.Sync(); err != nil {
			_ = err
		}
	}()

	logger.Debug("configuration loaded", zap.String("log_level", cfg.LogLevel), zap.Int("server_port", cfg.Server.Port), zap.String("addr", cfg.Addr))

	srv := newServer(cfg, logger)
	go func() {
		if err := srv.Start(cfg.Addr); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server start failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
	}
}

func newServer(cfg ...interface{}) *echo.Echo {
	var appConfig config.Config
	var logger *zap.Logger
	for _, item := range cfg {
		switch v := item.(type) {
		case config.Config:
			appConfig = v
		case *zap.Logger:
			logger = v
		}
	}
	if appConfig.Addr == "" {
		appConfig = config.Load()
	}
	if logger == nil {
		logger = newLogger(appConfig.LogLevel)
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogMethod: true,
		LogURI:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if c.Path() == "/healthz" || c.Path() == "/readyz" {
				return nil
			}
			traceID := c.Request().Header.Get("X-Trace-Id")
			if traceID == "" {
				traceID = "-"
			}
			logger.Info("request",
				zap.String("method", v.Method),
				zap.String("path", v.URI),
				zap.Int("status", v.Status),
				zap.Int64("duration_ms", v.Latency.Milliseconds()),
				zap.String("trace_id", traceID),
			)
			return nil
		},
	}))
	e.Use(middleware.Recover())

	healthService := service.NewHealthService()
	repo := repository.NewRepository()
	_ = repo
	apiHandler := api.NewHandler(healthService, appConfig.Version)
	apiHandler.Register(e)

	logger.Info("server listening", zap.String("addr", appConfig.Addr))
	return e
}

func newLogger(level string) *zap.Logger {
	var lvl zap.AtomicLevel
	if err := lvl.UnmarshalText([]byte(strings.ToLower(level))); err != nil {
		lvl.SetLevel(zap.InfoLevel)
	} else {
		lvl.SetLevel(lvl.Level())
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	cfg.Encoding = "json"
	cfg.DisableStacktrace = true
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := cfg.Build()
	return logger
}
