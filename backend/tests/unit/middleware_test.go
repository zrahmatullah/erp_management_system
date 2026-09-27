package unit_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"cafe-erp-system/backend/internal/config"
	"cafe-erp-system/backend/internal/delivery/http/middleware"
)

func TestMiddleware_InitAuth(t *testing.T) {
	// Nil config
	middleware.InitAuth(nil)

	// Development config
	cfg := &config.Config{
		Server: config.ServerConfig{Env: "development"},
		JWT:    config.JWTConfig{Secret: "test-secret-key-1234567890123456"},
	}
	middleware.InitAuth(cfg)

	// Production with valid secret
	prodCfg := &config.Config{
		Server: config.ServerConfig{Env: "production"},
		JWT:    config.JWTConfig{Secret: "production-ultra-secure-key-32chars-long!!"},
	}
	assert.NotPanics(t, func() {
		middleware.InitAuth(prodCfg)
	})

	// Production with short/default secret should panic
	badProdCfg := &config.Config{
		Server: config.ServerConfig{Env: "production"},
		JWT:    config.JWTConfig{Secret: "short"},
	}
	assert.Panics(t, func() {
		middleware.InitAuth(badProdCfg)
	})
}

func TestMiddleware_JWTAuth(t *testing.T) {
	middleware.SetJWTSecretForTesting("test-secret-key-for-unit-testing-32chars!!")
	userID := uuid.New()
	branchID := uuid.New()

	validAccess, _, err := middleware.GenerateTokenPair(userID, "cashier@cafe.com", "Cashier", &branchID, []string{"pos:create"})
	assert.NoError(t, err)

	handlerCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		claims, err := middleware.GetUserFromContext(r.Context())
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, "cashier@cafe.com", claims.Email)
		assert.Equal(t, "Cashier", claims.Role)
		w.WriteHeader(http.StatusOK)
	})

	jwtMiddleware := middleware.JWTAuth(nextHandler)

	t.Run("Valid Token Presentation", func(t *testing.T) {
		handlerCalled = false
		req := httptest.NewRequest("GET", "/api/v1/orders", nil)
		req.Header.Set("Authorization", "Bearer "+validAccess)
		rec := httptest.NewRecorder()

		jwtMiddleware.ServeHTTP(rec, req)
		assert.True(t, handlerCalled)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		handlerCalled = false
		req := httptest.NewRequest("GET", "/api/v1/orders", nil)
		rec := httptest.NewRecorder()

		jwtMiddleware.ServeHTTP(rec, req)
		assert.False(t, handlerCalled)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Missing authorization header")
	})

	t.Run("Invalid Header Format - Not Bearer", func(t *testing.T) {
		handlerCalled = false
		req := httptest.NewRequest("GET", "/api/v1/orders", nil)
		req.Header.Set("Authorization", "Basic "+validAccess)
		rec := httptest.NewRecorder()

		jwtMiddleware.ServeHTTP(rec, req)
		assert.False(t, handlerCalled)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid authorization header format")
	})

	t.Run("Invalid Header Format - Single Word", func(t *testing.T) {
		handlerCalled = false
		req := httptest.NewRequest("GET", "/api/v1/orders", nil)
		req.Header.Set("Authorization", "InvalidHeader")
		rec := httptest.NewRecorder()

		jwtMiddleware.ServeHTTP(rec, req)
		assert.False(t, handlerCalled)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Malformed Token String", func(t *testing.T) {
		handlerCalled = false
		req := httptest.NewRequest("GET", "/api/v1/orders", nil)
		req.Header.Set("Authorization", "Bearer malformed.jwt.token")
		rec := httptest.NewRecorder()

		jwtMiddleware.ServeHTTP(rec, req)
		assert.False(t, handlerCalled)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid or expired token")
	})
}

func TestMiddleware_GetUserFromContext_Failure(t *testing.T) {
	claims, err := middleware.GetUserFromContext(context.Background())
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, "user claims not found in context", err.Error())
}

func TestMiddleware_RBAC_RequireRole(t *testing.T) {
	middleware.SetJWTSecretForTesting("test-secret-key-for-unit-testing-32chars!!")
	userID := uuid.New()
	cashierToken, _, _ := middleware.GenerateTokenPair(userID, "cashier@cafe.com", "Cashier", nil, nil)
	managerToken, _, _ := middleware.GenerateTokenPair(userID, "mgr@cafe.com", "Branch Manager", nil, nil)

	handlerCalled := false
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	roleMiddleware := middleware.RequireRole("Branch Manager", "Super Admin")(protectedHandler)
	fullPipeline := middleware.JWTAuth(roleMiddleware)

	t.Run("Authorized Role Allowed", func(t *testing.T) {
		handlerCalled = false
		req := httptest.NewRequest("GET", "/api/v1/reports", nil)
		req.Header.Set("Authorization", "Bearer "+managerToken)
		rec := httptest.NewRecorder()

		fullPipeline.ServeHTTP(rec, req)
		assert.True(t, handlerCalled)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Unauthorized Role Forbidden", func(t *testing.T) {
		handlerCalled = false
		req := httptest.NewRequest("GET", "/api/v1/reports", nil)
		req.Header.Set("Authorization", "Bearer "+cashierToken)
		rec := httptest.NewRecorder()

		fullPipeline.ServeHTTP(rec, req)
		assert.False(t, handlerCalled)
		assert.Equal(t, http.StatusForbidden, rec.Code)
		assert.Contains(t, rec.Body.String(), "Forbidden: insufficient role")
	})

	t.Run("Missing Claims in Context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reports", nil)
		rec := httptest.NewRecorder()

		roleMiddleware.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestMiddleware_RBAC_RequirePermission(t *testing.T) {
	middleware.SetJWTSecretForTesting("test-secret-key-for-unit-testing-32chars!!")
	userID := uuid.New()

	posToken, _, _ := middleware.GenerateTokenPair(userID, "pos@cafe.com", "Cashier", nil, []string{"pos:create", "pos:read"})
	inventoryToken, _, _ := middleware.GenerateTokenPair(userID, "inv@cafe.com", "Barista", nil, []string{"inventory:*"})
	superAdminToken, _, _ := middleware.GenerateTokenPair(userID, "admin@cafe.com", "Super Admin", nil, nil)

	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	permMiddleware := middleware.RequirePermission("inventory", "update")(protectedHandler)
	fullPipeline := middleware.JWTAuth(permMiddleware)

	t.Run("Module Wildcard Grants Access", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/inventory/adjust", nil)
		req.Header.Set("Authorization", "Bearer "+inventoryToken)
		rec := httptest.NewRecorder()

		fullPipeline.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Super Admin Bypasses Permission Check", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/inventory/adjust", nil)
		req.Header.Set("Authorization", "Bearer "+superAdminToken)
		rec := httptest.NewRecorder()

		fullPipeline.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Missing Permission Denied", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/inventory/adjust", nil)
		req.Header.Set("Authorization", "Bearer "+posToken)
		rec := httptest.NewRecorder()

		fullPipeline.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
		assert.Contains(t, rec.Body.String(), "Forbidden: missing permission inventory:update")
	})

	t.Run("Missing Claims in Context", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/inventory/adjust", nil)
		rec := httptest.NewRecorder()

		permMiddleware.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestMiddleware_CORS(t *testing.T) {
	os.Setenv("CORS_ALLOWED_ORIGINS", "https://erp.mycafe.com, https://admin.mycafe.com")
	defer os.Unsetenv("CORS_ALLOWED_ORIGINS")

	corsMiddleware := middleware.CORS()
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("Allowed Origin Request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/ping", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "http://localhost:5173", rec.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("Environment Allowed Origin Request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/ping", nil)
		req.Header.Set("Origin", "https://erp.mycafe.com")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "https://erp.mycafe.com", rec.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("Preflight OPTIONS Request", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/api/v1/orders", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Access-Control-Request-Headers", "Authorization, Content-Type")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "POST")
	})
}

func TestMiddleware_Recovery(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("simulated fatal runtime panic")
	})

	recoveryMiddleware := middleware.Recovery(panicHandler)
	req := httptest.NewRequest("GET", "/api/v1/panic", nil)
	rec := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		recoveryMiddleware.ServeHTTP(rec, req)
	})

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Internal Server Error")
}

func TestMiddleware_SecurityHeaders(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	secMiddleware := middleware.SecurityHeaders(nextHandler)
	req := httptest.NewRequest("GET", "/api/v1/status", nil)
	rec := httptest.NewRecorder()

	secMiddleware.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", rec.Header().Get("X-Frame-Options"))
	assert.Equal(t, "strict-origin-when-cross-origin", rec.Header().Get("Referrer-Policy"))
	assert.Equal(t, "0", rec.Header().Get("X-XSS-Protection"))
	assert.Contains(t, rec.Header().Get("Strict-Transport-Security"), "max-age=31536000")
	assert.Contains(t, rec.Header().Get("Content-Security-Policy"), "default-src 'self'")
}

func TestMiddleware_RequestSizeLimit(t *testing.T) {
	bodyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 1024)
		_, err := r.Body.Read(buf)
		if err != nil && err.Error() == "http: request body too large" {
			http.Error(w, "Payload Too Large", http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	limitMiddleware := middleware.RequestSizeLimit(10)(bodyHandler)

	t.Run("Payload Within Limit", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/data", strings.NewReader("small"))
		rec := httptest.NewRecorder()
		limitMiddleware.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Payload Exceeds Limit", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/data", strings.NewReader("this payload exceeds 10 bytes limit!"))
		rec := httptest.NewRecorder()
		limitMiddleware.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
	})
}

func TestMiddleware_RateLimit(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("GetClientIP Extracts Correct Header", func(t *testing.T) {
		req1 := httptest.NewRequest("GET", "/", nil)
		req1.Header.Set("X-Forwarded-For", "203.0.113.195, 70.41.3.18")
		assert.Equal(t, "203.0.113.195", middleware.GetClientIP(req1))

		req2 := httptest.NewRequest("GET", "/", nil)
		req2.Header.Set("X-Real-IP", "198.51.100.1")
		assert.Equal(t, "198.51.100.1", middleware.GetClientIP(req2))

		req3 := httptest.NewRequest("GET", "/", nil)
		req3.RemoteAddr = "192.0.2.1:54321"
		assert.Equal(t, "192.0.2.1", middleware.GetClientIP(req3))
	})

	t.Run("AuthRateLimit Allows Normal Traffic and Blocks Abusive Rate", func(t *testing.T) {
		limiter := middleware.AuthRateLimit(handler)
		ip := "10.20.30.40"

		allowedCount := 0
		blockedCount := 0
		for i := 0; i < 20; i++ {
			req := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
			req.Header.Set("X-Real-IP", ip)
			rec := httptest.NewRecorder()
			limiter.ServeHTTP(rec, req)

			if rec.Code == http.StatusOK {
				allowedCount++
			} else if rec.Code == http.StatusTooManyRequests {
				blockedCount++
			}
		}

		assert.Greater(t, allowedCount, 0)
		assert.Greater(t, blockedCount, 0)
	})

	t.Run("RateLimit General API Limiter", func(t *testing.T) {
		limiter := middleware.RateLimit(handler)
		req := httptest.NewRequest("GET", "/api/v1/products", nil)
		req.Header.Set("X-Real-IP", "172.16.0.50")
		rec := httptest.NewRecorder()
		limiter.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestMiddleware_RequestLogger(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	loggerMiddleware := middleware.RequestLogger(handler)
	req := httptest.NewRequest("GET", "/api/v1/ping", nil)
	rec := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		loggerMiddleware.ServeHTTP(rec, req)
	})
	assert.Equal(t, http.StatusOK, rec.Code)
}
