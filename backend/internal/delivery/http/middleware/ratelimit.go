package middleware

import (
	"net/http"
	"sync"
	"time"

	"cafe-erp-system/backend/pkg/response"
)

type client struct {
	tokens int
	last   time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

const (
	rateLimit     = 100 // reqs per window
	window        = time.Minute
	cleanupPeriod = 5 * time.Minute
)

func init() {
	go cleanupClients()
}

func cleanupClients() {
	for {
		time.Sleep(cleanupPeriod)
		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.last) > cleanupPeriod {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		mu.Lock()

		c, exists := clients[ip]
		if !exists || time.Since(c.last) > window {
			clients[ip] = &client{tokens: rateLimit - 1, last: time.Now()}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if c.tokens <= 0 {
			mu.Unlock()
			response.Error(w, http.StatusTooManyRequests, "Too many requests")
			return
		}

		c.tokens--
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
