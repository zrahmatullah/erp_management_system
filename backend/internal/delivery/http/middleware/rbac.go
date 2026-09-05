package middleware

import (
	"net/http"

	"cafe-erp-system/backend/pkg/response"
)

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := GetUserFromContext(r.Context())
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			hasRole := false
			for _, role := range roles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				response.Error(w, http.StatusForbidden, "Forbidden: insufficient role")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequirePermission(module, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := GetUserFromContext(r.Context())
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			// Admin bypasses permission check
			if claims.Role == "Super Admin" {
				next.ServeHTTP(w, r)
				return
			}

			requiredPerm := module + ":" + action
			hasPerm := false
			for _, perm := range claims.Permissions {
				if perm == requiredPerm || perm == module+":*" || perm == "*:*" {
					hasPerm = true
					break
				}
			}

			if !hasPerm {
				response.Error(w, http.StatusForbidden, "Forbidden: missing permission "+requiredPerm)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
