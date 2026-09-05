package middleware

import (
	"github.com/go-chi/cors"
	"net/http"
)

func CORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "http://127.0.0.1:5173", "http://localhost:8080", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "Content-Disposition"},
		AllowCredentials: false,
		MaxAge:           300, 
	})
}
