package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"cafe-erp-system/backend/pkg/response"
)

type client struct {
	tokens int
	last   time.Time
}

type rateLimiter struct {
	clients       map[string]*client
	mu            sync.Mutex
	rateLimit     int
	window        time.Duration
	cleanupPeriod time.Duration
}

func newRateLimiter(limit int, win, cleanup time.Duration) *rateLimiter {
	rl := &rateLimiter{
		clients:       make(map[string]*client),
		rateLimit:     limit,
		window:        win,
		cleanupPeriod: cleanup,
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *rateLimiter) cleanupLoop() {
	for {
		time.Sleep(rl.cleanupPeriod)
		rl.mu.Lock()
		for ip, c := range rl.clients {
			if time.Since(c.last) > rl.cleanupPeriod {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// GetClientIP extracts real client IP address ignoring dynamic client ports and respecting proxy headers.
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		ip := strings.TrimSpace(xrip)
		if ip != "" {
			return ip
		}
	}

	// Fall back to RemoteAddr stripped of port
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}

	return r.RemoteAddr
}

var (
	// Global API rate limiter (120 reqs/min)
	globalLimiter = newRateLimiter(120, time.Minute, 5*time.Minute)

	// Specialized Auth Login limiter to prevent brute-force attacks (10 reqs/min)
	authLimiter = newRateLimiter(10, time.Minute, 5*time.Minute)
)

// RateLimit is the general API rate limiter middleware.
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetClientIP(r)

		globalLimiter.mu.Lock()
		c, exists := globalLimiter.clients[ip]
		if !exists || time.Since(c.last) > globalLimiter.window {
			globalLimiter.clients[ip] = &client{tokens: globalLimiter.rateLimit - 1, last: time.Now()}
			globalLimiter.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if c.tokens <= 0 {
			globalLimiter.mu.Unlock()
			w.Header().Set("Retry-After", "60")
			response.Error(w, http.StatusTooManyRequests, "Terlalu banyak permintaan. Silakan tunggu beberapa saat.")
			return
		}

		c.tokens--
		globalLimiter.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

// AuthRateLimit restricts authentication attempts to defend against brute force attacks.
func AuthRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetClientIP(r)

		authLimiter.mu.Lock()
		c, exists := authLimiter.clients[ip]
		if !exists || time.Since(c.last) > authLimiter.window {
			authLimiter.clients[ip] = &client{tokens: authLimiter.rateLimit - 1, last: time.Now()}
			authLimiter.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if c.tokens <= 0 {
			authLimiter.mu.Unlock()
			w.Header().Set("Retry-After", "60")
			response.Error(w, http.StatusTooManyRequests, "Batas percobaan login tercapai. Demi keamanan akun, coba lagi dalam 1 menit.")
			return
		}

		c.tokens--
		authLimiter.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
