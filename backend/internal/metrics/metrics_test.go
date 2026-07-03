package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestMetricsHandlerExposesPrometheusMetrics(t *testing.T) {
	m := NewMetrics()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	m.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "habitflow_http_requests_total") {
		t.Fatalf("expected metrics output to include habitflow_http_requests_total, got %s", body)
	}
}

func TestMiddlewareRecordsRequestMetrics(t *testing.T) {
	m := NewMetrics()
	e := echo.New()
	e.Use(m.Middleware())
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusCreated, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	metricsReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metricsRec := httptest.NewRecorder()
	m.Handler().ServeHTTP(metricsRec, metricsReq)

	body := metricsRec.Body.String()
	if !strings.Contains(body, "habitflow_http_requests_total") {
		t.Fatalf("expected request metrics to be exposed, got %s", body)
	}
	if !strings.Contains(body, `method="GET"`) || !strings.Contains(body, `status_code="201"`) {
		t.Fatalf("expected request metrics to include labels, got %s", body)
	}
}
