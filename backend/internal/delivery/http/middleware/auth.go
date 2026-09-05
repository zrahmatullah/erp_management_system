package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"cafe-erp-system/backend/internal/config"
	"cafe-erp-system/backend/pkg/response"
)

type contextKey string

const UserContextKey contextKey = "user_claims"

type Claims struct {
	UserID      uuid.UUID `json:"user_id"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	BranchID    *uuid.UUID `json:"branch_id,omitempty"`
	Permissions []string  `json:"permissions"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

func InitAuth(cfg *config.Config) {
	jwtSecret = []byte(cfg.JWT.Secret)
}

func GenerateTokenPair(userID uuid.UUID, email, role string, branchID *uuid.UUID, permissions []string) (string, string, error) {
	// Access Token - 15 mins
	accessClaims := &Claims{
		UserID:      userID,
		Email:       email,
		Role:        role,
		BranchID:    branchID,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token - 7 days
	refreshClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		claims, err := ValidateToken(parts[1])
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(UserContextKey).(*Claims)
	if !ok {
		return nil, errors.New("user claims not found in context")
	}
	return claims, nil
}
