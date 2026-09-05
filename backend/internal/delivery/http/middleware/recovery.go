package middleware

import (
	"net/http"
	"runtime/debug"

	"cafe-erp-system/backend/pkg/logger"
	"cafe-erp-system/backend/pkg/response"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.Error().
					Interface("error", err).
					Bytes("stack", debug.Stack()).
					Msg("Panic recovered")

				response.Error(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
