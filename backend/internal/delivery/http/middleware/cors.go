package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/cors"
)

// CORS configures strict, explicit Cross-Origin Resource Sharing adhering to OWASP guidelines.
func CORS() func(http.Handler) http.Handler {
	allowedOrigins := []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:3000",
		"http://localhost:8080",
	}

	// Support extra origins specified in production env (e.g. https://erp.yourdomain.com)
	if extraOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); extraOrigins != "" {
		for _, o := range strings.Split(extraOrigins, ",") {
			trimmed := strings.TrimSpace(o)
			if trimmed != "" && trimmed != "*" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "Content-Disposition", "Retry-After", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
