package unit_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"cafe-erp-system/backend/internal/config"
	"cafe-erp-system/backend/internal/delivery/http/middleware"
	"cafe-erp-system/backend/pkg/crypto"
)

// TestAuth_PasswordSecurity tests bcrypt password hashing, salting, and verification.
func TestAuth_PasswordSecurity(t *testing.T) {
	password := "SecretP@ssw0rd!2026"

	// 1. Hash password
	hashed, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("crypto.HashPassword failed: %v", err)
	}

	if hashed == password {
		t.Fatalf("password hash must not match plaintext password")
	}

	// 2. Validate correct password
	if !crypto.CheckPassword(password, hashed) {
		t.Errorf("crypto.CheckPassword returned false for correct password")
	}

	// 3. Reject incorrect password
	if crypto.CheckPassword("WrongPassword123", hashed) {
		t.Errorf("crypto.CheckPassword returned true for incorrect password")
	}

	// 4. Reject empty password against non-empty hash
	if crypto.CheckPassword("", hashed) {
		t.Errorf("crypto.CheckPassword returned true for empty password")
	}
}

// TestAuth_JWTTokenLifecycle tests token generation, claims preservation, and signature verification.
func TestAuth_JWTTokenLifecycle(t *testing.T) {
	testSecret := "super-secure-cafe-erp-testing-secret-key-32chars"
	middleware.SetJWTSecretForTesting(testSecret)

	userID := uuid.New()
	branchID := uuid.New()
	email := "cashier@cafe-erp.local"
	role := "Cashier"
	permissions := []string{"pos:read", "pos:create", "pos:pay"}

	// 1. Generate Token Pair
	accessToken, refreshToken, err := middleware.GenerateTokenPair(userID, email, role, &branchID, permissions)
	if err != nil {
		t.Fatalf("GenerateTokenPair failed: %v", err)
	}
	if accessToken == "" || refreshToken == "" {
		t.Fatalf("tokens must not be empty string")
	}

	// 2. Validate Access Token Claims
	claims, err := middleware.ValidateToken(accessToken)
	if err != nil {
		t.Fatalf("ValidateToken failed on valid token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("claims.UserID = %v; want %v", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("claims.Email = %v; want %v", claims.Email, email)
	}
	if claims.Role != role {
		t.Errorf("claims.Role = %v; want %v", claims.Role, role)
	}
	if claims.BranchID == nil || *claims.BranchID != branchID {
		t.Errorf("claims.BranchID = %v; want %v", claims.BranchID, branchID)
	}
	if len(claims.Permissions) != 3 {
		t.Errorf("claims.Permissions length = %d; want 3", len(claims.Permissions))
	}

	// 3. Reject Tampered Signature
	tamperedToken := accessToken[:len(accessToken)-6] + "xxxxxx"
	_, err = middleware.ValidateToken(tamperedToken)
	if err == nil {
		t.Errorf("ValidateToken must reject token with tampered signature")
	}

	// 4. Reject Expired Token
	expiredClaims := &middleware.Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	expiredTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenStr, _ := expiredTokenObj.SignedString([]byte(testSecret))

	_, err = middleware.ValidateToken(expiredTokenStr)
	if err == nil {
		t.Errorf("ValidateToken must reject expired token")
	}
}

// TestRBAC_PermissionMatrix tests fine-grained permission rules and wildcard matching.
func TestRBAC_PermissionMatrix(t *testing.T) {
	tests := []struct {
		name           string
		role           string
		permissions    []string
		targetModule   string
		targetAction   string
		expectedAccess bool
	}{
		{
			name:           "Super Admin Unconditional Bypass",
			role:           "Super Admin",
			permissions:    nil, // No explicit permissions needed
			targetModule:   "finance",
			targetAction:   "delete_journal",
			expectedAccess: true,
		},
		{
			name:           "Global Wildcard Access (*:*)",
			role:           "General Manager",
			permissions:    []string{"*:*"},
			targetModule:   "master",
			targetAction:   "delete_user",
			expectedAccess: true,
		},
		{
			name:           "Module Wildcard Access (pos:*)",
			role:           "Cashier Lead",
			permissions:    []string{"pos:*"},
			targetModule:   "pos",
			targetAction:   "void_order",
			expectedAccess: true,
		},
		{
			name:           "Module Wildcard Does Not Grant Other Modules",
			role:           "Cashier",
			permissions:    []string{"pos:*"},
			targetModule:   "inventory",
			targetAction:   "write",
			expectedAccess: false,
		},
		{
			name:           "Exact Match Permission Allowed",
			role:           "Warehouse Staff",
			permissions:    []string{"inventory:read", "inventory:stock_opname"},
			targetModule:   "inventory",
			targetAction:   "stock_opname",
			expectedAccess: true,
		},
		{
			name:           "Unauthorized Action in Same Module",
			role:           "Warehouse Staff",
			permissions:    []string{"inventory:read"},
			targetModule:   "inventory",
			targetAction:   "delete_item",
			expectedAccess: false,
		},
		{
			name:           "Empty Permissions List Denied",
			role:           "Guest Staff",
			permissions:    []string{},
			targetModule:   "pos",
			targetAction:   "read",
			expectedAccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasAccess := middleware.HasPermission(tt.role, tt.permissions, tt.targetModule, tt.targetAction)
			if hasAccess != tt.expectedAccess {
				t.Errorf("HasPermission(%s, %v, %s, %s) = %v; want %v",
					tt.role, tt.permissions, tt.targetModule, tt.targetAction, hasAccess, tt.expectedAccess)
			}
		})
	}
}

// TestRBAC_RoleValidation tests role-based restrictions.
func TestRBAC_RoleValidation(t *testing.T) {
	tests := []struct {
		name         string
		userRole     string
		allowedRoles []string
		want         bool
	}{
		{
			name:         "Role Matches Directly",
			userRole:     "Store Manager",
			allowedRoles: []string{"Super Admin", "Store Manager"},
			want:         true,
		},
		{
			name:         "Role Missing in Allowed Roles",
			userRole:     "Barista",
			allowedRoles: []string{"Super Admin", "Store Manager"},
			want:         false,
		},
		{
			name:         "Empty Allowed Roles List",
			userRole:     "Super Admin",
			allowedRoles: []string{},
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched := middleware.HasRole(tt.userRole, tt.allowedRoles...)
			if matched != tt.want {
				t.Errorf("HasRole(%s, %v) = %v; want %v", tt.userRole, tt.allowedRoles, matched, tt.want)
			}
		})
	}
}

// TestAuth_ProductionConfigGuard tests that insecure JWT configuration triggers panic in production.
func TestAuth_ProductionConfigGuard(t *testing.T) {
	// 1. Should panic when JWT_SECRET is empty in production
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("InitAuth must panic when JWT secret is empty in production")
		}
	}()

	insecureCfg := &config.Config{
		Server: config.ServerConfig{Env: "production"},
		JWT:    config.JWTConfig{Secret: ""},
	}
	middleware.InitAuth(insecureCfg)
}
