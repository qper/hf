package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/qper/hf/internal/api"
	"github.com/qper/hf/internal/config"
	"github.com/qper/hf/internal/metrics"
	"github.com/qper/hf/internal/migrations"
	"github.com/qper/hf/internal/repository"
	"github.com/qper/hf/internal/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--stop" {
		stopServer()
		return
	}

	cfg := config.Load()
	cfg.Version = version

	logger := newLogger(cfg.LogLevel)
	defer func() {
		if err := logger.Sync(); err != nil {
			_ = err
		}
	}()

	logger.Debug("configuration loaded", zap.String("log_level", cfg.LogLevel), zap.Int("server_port", cfg.Server.Port), zap.String("addr", cfg.Addr))

	if err := stopExistingServer(logger); err != nil {
		logger.Warn("failed to stop existing server instance", zap.Error(err))
	}
	if err := stopPortListeners(8080, 9090); err != nil {
		logger.Warn("failed to stop port listeners", zap.Error(err))
	}
	if err := writePIDFile(); err != nil {
		logger.Error("failed to write pid file", zap.Error(err))
		os.Exit(1)
	}
	defer removePIDFile()

	result, err := migrations.RunWithDSN(cfg.DB.DSN)
	if err != nil {
		logger.Error("database migrations failed", zap.Error(err))
		os.Exit(1)
	}
	if result.NoChange {
		logger.Info("database migrations completed", zap.Bool("no_change", true), zap.Uint64("version", uint64(result.Version)))
	} else {
		logger.Info("database migrations completed", zap.Int("applied", result.Applied), zap.Uint64("version", uint64(result.Version)))
	}

	appMetrics := metrics.NewMetrics()
	srv := newServer(cfg, logger, appMetrics)
	metricsServer := newMetricsServer(appMetrics)
	go func() {
		if err := srv.Start(cfg.Addr); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server start failed", zap.Error(err))
		}
	}()
	go func() {
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("metrics server start failed", zap.Error(err))
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
	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.Error("metrics server shutdown failed", zap.Error(err))
	}
}

func newServer(cfg ...interface{}) *echo.Echo {
	var appConfig config.Config
	var logger *zap.Logger
	var appMetrics *metrics.Metrics
	for _, item := range cfg {
		switch v := item.(type) {
		case config.Config:
			appConfig = v
		case *zap.Logger:
			logger = v
		case *metrics.Metrics:
			appMetrics = v
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
	if appMetrics == nil {
		appMetrics = metrics.NewMetrics()
	}
	e.Use(appMetrics.Middleware())
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

	var authService api.AuthService
	var db *sql.DB
	if openedDB, err := sql.Open("postgres", appConfig.DB.DSN); err != nil {
		logger.Warn("failed to open auth database connection", zap.Error(err))
	} else {
		db = openedDB
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err != nil {
			logger.Warn("auth database unavailable, registration endpoint disabled", zap.Error(err))
			_ = db.Close()
			db = nil
		} else {
			authRepo := repository.NewAuthRepository(db)
			authService = service.NewAuthService(authRepo)
		}
	}

	apiHandler := api.NewHandlerWithAuth(healthService, appConfig.Version, authService)
	if db != nil {
		habitRepo := repository.NewHabitRepository(db)
		habitService := service.NewHabitService(habitRepo)
		apiHandler = api.NewHandlerWithHabit(healthService, appConfig.Version, authService, habitService)
		apiHandler = apiHandler.WithDBChecker(api.NewDBChecker(db))
	}
	apiHandler.Register(e)

	logger.Info("server listening", zap.String("addr", appConfig.Addr))
	return e
}

func newMetricsServer(appMetrics *metrics.Metrics) *http.Server {
	if appMetrics == nil {
		appMetrics = metrics.NewMetrics()
	}
	return &http.Server{
		Addr: ":9090",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/metrics" {
				http.NotFound(w, r)
				return
			}
			appMetrics.Handler().ServeHTTP(w, r)
		}),
	}
}

func stopServer() {
	pidFile := pidFilePath()
	data, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "no running server pid file found at %s\n", pidFile)
		os.Exit(0)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid pid in %s: %v\n", pidFile, err)
		os.Exit(1)
	}

	if err := stopProcessByPID(pid); err != nil {
		fmt.Fprintf(os.Stderr, "failed to stop process %d: %v\n", pid, err)
		os.Exit(1)
	}
	_ = stopPortListeners(8080, 9090)
	_ = os.Remove(pidFile)
	fmt.Printf("stopped server pid %d\n", pid)
}

func stopExistingServer(logger *zap.Logger) error {
	pidFile := pidFilePath()
	data, err := os.ReadFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		_ = os.Remove(pidFile)
		return nil
	}

	if err := stopProcessByPID(pid); err != nil {
		return err
	}
	_ = os.Remove(pidFile)
	if logger != nil {
		logger.Info("stopped previous server instance", zap.Int("pid", pid))
	}
	return nil
}

func stopPortListeners(ports ...int) error {
	for _, port := range ports {
		cmd := exec.Command("lsof", "-nP", "-t", fmt.Sprintf("-iTCP:%d", port), "-sTCP:LISTEN")
		output, err := cmd.CombinedOutput()
		if err != nil {
			continue
		}

		for _, pidLine := range strings.Split(strings.TrimSpace(string(output)), "\n") {
			pidLine = strings.TrimSpace(pidLine)
			if pidLine == "" {
				continue
			}
			pid, err := strconv.Atoi(pidLine)
			if err != nil {
				continue
			}
			if err := stopProcessByPID(pid); err != nil {
				continue
			}
		}
	}
	return nil
}

func stopProcessByPID(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	if err := process.Signal(syscall.Signal(0)); err != nil {
		return nil
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		return err
	}

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if err := process.Signal(syscall.Signal(0)); err != nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	_ = process.Signal(syscall.SIGKILL)
	return nil
}

func writePIDFile() error {
	pidFile := pidFilePath()
	pid := os.Getpid()
	return os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0o600)
}

func removePIDFile() {
	_ = os.Remove(pidFilePath())
}

func pidFilePath() string {
	return filepath.Join(os.TempDir(), "hf-server.pid")
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
