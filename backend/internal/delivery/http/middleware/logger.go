package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"cafe-erp-system/backend/pkg/logger"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		
		next.ServeHTTP(ww, r)

		logger.Log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.Status()).
			Dur("duration", time.Since(start)).
			Str("ip", r.RemoteAddr).
			Str("req_id", middleware.GetReqID(r.Context())).
			Msg("HTTP Request")
	})
}
