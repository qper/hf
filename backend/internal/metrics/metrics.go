package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	registry         *prometheus.Registry
	requestsTotal    *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
	requestsInFlight *prometheus.GaugeVec
}

func NewMetrics() *Metrics {
	registry := prometheus.NewRegistry()
	m := &Metrics{
		registry: registry,
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "habitflow_http_requests_total",
				Help: "Total number of HTTP requests processed by the API.",
			},
			[]string{"method", "status_code", "path"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "habitflow_http_request_duration_seconds",
				Help:    "Histogram of HTTP request latencies in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "status_code", "path"},
		),
		requestsInFlight: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "habitflow_http_requests_in_flight",
				Help: "Number of HTTP requests currently in flight.",
			},
			[]string{"method", "path"},
		),
	}
	registry.MustRegister(m.requestsTotal, m.requestDuration, m.requestsInFlight)
	m.requestsTotal.WithLabelValues("", "", "").Add(0)
	m.requestDuration.WithLabelValues("", "", "").Observe(0)
	m.requestsInFlight.WithLabelValues("", "").Add(0)
	return m
}

func (m *Metrics) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			method := c.Request().Method
			path := c.Path()
			if path == "" {
				path = c.Request().URL.Path
			}
			m.requestsInFlight.WithLabelValues(method, path).Inc()
			defer m.requestsInFlight.WithLabelValues(method, path).Dec()

			err := next(c)
			status := c.Response().Status
			if status == 0 {
				status = http.StatusOK
			}
			duration := time.Since(start).Seconds()
			m.requestsTotal.WithLabelValues(method, strconv.Itoa(status), path).Inc()
			m.requestDuration.WithLabelValues(method, strconv.Itoa(status), path).Observe(duration)
			return err
		}
	}
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
