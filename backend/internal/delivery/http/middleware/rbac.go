package middleware

import (
	"net/http"

	"cafe-erp-system/backend/pkg/response"
)

// HasRole determines whether the given userRole matches any of the allowed roles.
func HasRole(userRole string, allowedRoles ...string) bool {
	for _, role := range allowedRoles {
		if userRole == role {
			return true
		}
	}
	return false
}

// HasPermission determines whether the given role or list of permissions
// grants access to the specified module and action.
//
// Rules:
// 1. "Super Admin" role bypasses all checks (returns true).
// 2. Global wildcard "*:*" grants access to all modules and actions.
// 3. Module wildcard "<module>:*" grants access to any action in that module.
// 4. Exact match "<module>:<action>" grants access.
func HasPermission(role string, permissions []string, module, action string) bool {
	if role == "Super Admin" {
		return true
	}

	requiredPerm := module + ":" + action
	moduleWildcard := module + ":*"

	for _, perm := range permissions {
		if perm == requiredPerm || perm == moduleWildcard || perm == "*:*" {
			return true
		}
	}
	return false
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := GetUserFromContext(r.Context())
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			if !HasRole(claims.Role, roles...) {
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

			if !HasPermission(claims.Role, claims.Permissions, module, action) {
				response.Error(w, http.StatusForbidden, "Forbidden: missing permission "+module+":"+action)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
