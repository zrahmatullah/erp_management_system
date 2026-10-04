package middleware

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"cafe-erp-system/backend/pkg/response"
)

// HeavyRouteLimiter restricts the number of concurrent in-flight requests to resource-heavy endpoints
// (such as reporting, batch payroll, stock opname audits, and reconciliation) to prevent
// PostgreSQL connection pool starvation and CPU exhaustion.
type HeavyRouteLimiter struct {
	sem     chan struct{}
	timeout time.Duration
}

// Global default limiter with sensible capacity for medium/large deployments
var DefaultHeavyLimiter *HeavyRouteLimiter

func init() {
	capacity := 10
	if envCap := os.Getenv("HEAVY_MAX_CONCURRENCY"); envCap != "" {
		if c, err := strconv.Atoi(envCap); err == nil && c > 0 {
			capacity = c
		}
	}
	DefaultHeavyLimiter = NewHeavyRouteLimiter(capacity, 2*time.Second)
}

// NewHeavyRouteLimiter creates a new semaphore limiter with the specified maximum concurrency and wait timeout.
func NewHeavyRouteLimiter(maxConcurrent int, acquireTimeout time.Duration) *HeavyRouteLimiter {
	if maxConcurrent <= 0 {
		maxConcurrent = 10
	}
	return &HeavyRouteLimiter{
		sem:     make(chan struct{}, maxConcurrent),
		timeout: acquireTimeout,
	}
}

// Limit is a Chi/net/http middleware that protects downstream handlers.
func (l *HeavyRouteLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if l.timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, l.timeout)
			defer cancel()
		}

		select {
		case l.sem <- struct{}{}:
			defer func() { <-l.sem }()
			next.ServeHTTP(w, r)
		case <-ctx.Done():
			response.Error(
				w,
				http.StatusTooManyRequests,
				"Sistem sedang memproses antrean laporan/kalkulasi intensif. Silakan coba beberapa saat lagi.",
			)
		}
	})
}

// HeavyLimitMiddleware is a shortcut handler function using the DefaultHeavyLimiter.
func HeavyLimitMiddleware(next http.Handler) http.Handler {
	return DefaultHeavyLimiter.Limit(next)
}
