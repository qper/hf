package api

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// simple in-memory sliding window rate limiter per IP
type rateLimiter struct {
	mu     sync.Mutex
	store  map[string][]time.Time
	limit  int
	window time.Duration
}

func NewRateLimiter(limit int, window time.Duration) echo.MiddlewareFunc {
	rl := &rateLimiter{store: make(map[string][]time.Time), limit: limit, window: window}
	return rl.middleware
}

func (r *rateLimiter) getIP(c echo.Context) string {
	// prefer X-Forwarded-For
	xff := c.Request().Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(c.Request().RemoteAddr)
	if err != nil {
		return c.Request().RemoteAddr
	}
	return host
}

func (r *rateLimiter) middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := r.getIP(c)
		now := time.Now()

		r.mu.Lock()
		arr := r.store[ip]
		// prune
		cutoff := now.Add(-r.window)
		j := 0
		for i := range arr {
			if arr[i].After(cutoff) {
				arr[j] = arr[i]
				j++
			}
		}
		arr = arr[:j]

		if len(arr) >= r.limit {
			r.mu.Unlock()
			// too many
			c.Response().Header().Set("Retry-After", "900")
			return c.JSON(http.StatusTooManyRequests, map[string]string{"error": "rate limit exceeded"})
		}

		// record
		arr = append(arr, now)
		r.store[ip] = arr
		r.mu.Unlock()

		return next(c)
	}
}

// simple CORS middleware
func CORSMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h := c.Response().Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		h.Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if c.Request().Method == http.MethodOptions {
			return c.NoContent(http.StatusNoContent)
		}
		return next(c)
	}
}

// simple CSP middleware
func CSPMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")
		return next(c)
	}
}
