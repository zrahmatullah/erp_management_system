package middleware

import (
	"net/http"
)

// SecurityHeaders applies OWASP recommended security headers to every HTTP response.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent MIME-sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent Clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Control Referrer information leakage
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Restrict access to sensitive browser device features
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")

		// Disable buggy legacy XSS filter in favor of CSP
		w.Header().Set("X-XSS-Protection", "0")

		// Enforce HTTPS transport security (1 year with subdomains)
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content Security Policy (CSP) mitigating XSS, data injection, and malicious frames
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https: blob:; " +
			"font-src 'self' data: https:; " +
			"connect-src 'self' http://localhost:* ws://localhost:* https:; " +
			"frame-ancestors 'none'; " +
			"object-src 'none'; " +
			"base-uri 'self';"
		w.Header().Set("Content-Security-Policy", csp)

		next.ServeHTTP(w, r)
	})
}

// RequestSizeLimit restricts the maximum readable body payload to prevent DoS / Memory exhaustion.
func RequestSizeLimit(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}
