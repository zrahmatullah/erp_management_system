package unit_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"cafe-erp-system/backend/internal/delivery/http/handler"
	"cafe-erp-system/backend/internal/domain"
)

// MockAuthUsecase implements domain.AuthUsecase
type MockAuthUsecase struct {
	loginFunc        func(ctx context.Context, req domain.LoginRequest) (domain.LoginResponse, error)
	refreshTokenFunc func(ctx context.Context, token string) (domain.LoginResponse, error)
}

func (m *MockAuthUsecase) Login(ctx context.Context, req domain.LoginRequest) (domain.LoginResponse, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, req)
	}
	return domain.LoginResponse{}, errors.New("unimplemented")
}

func (m *MockAuthUsecase) Register(ctx context.Context, req domain.RegisterRequest) (domain.User, error) {
	return domain.User{}, nil
}

func (m *MockAuthUsecase) RefreshToken(ctx context.Context, token string) (domain.LoginResponse, error) {
	if m.refreshTokenFunc != nil {
		return m.refreshTokenFunc(ctx, token)
	}
	return domain.LoginResponse{}, errors.New("unimplemented")
}

func (m *MockAuthUsecase) ForgotPassword(ctx context.Context, email string) error {
	return nil
}

func (m *MockAuthUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	return nil
}

func (m *MockAuthUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (domain.LoginResponse, error) {
	return domain.LoginResponse{
		User: domain.User{
			BaseEntity: domain.BaseEntity{ID: userID},
			Email:      "admin@cafe.com",
			Username:   "admin",
		},
		Role: "Super Admin",
	}, nil
}


func TestAuthHandler_Login(t *testing.T) {
	mockUsecase := &MockAuthUsecase{}
	h := handler.NewAuthHandler(mockUsecase)

	t.Run("Happy Path - Successful Login sets Secure Cookie and Returns 200", func(t *testing.T) {
		userID := uuid.New()
		mockUsecase.loginFunc = func(ctx context.Context, req domain.LoginRequest) (domain.LoginResponse, error) {
			return domain.LoginResponse{
				Token:        "access-token-xyz",
				RefreshToken: "refresh-token-abc",
				User: domain.User{
					BaseEntity: domain.BaseEntity{ID: userID},
					Email:      "barista@cafe.com",
					Username:   "barista",
				},
				Role: "Barista",
			}, nil
		}

		payload, _ := json.Marshal(domain.LoginRequest{
			Email:    "barista@cafe.com",
			Password: "CorrectPassword123!",
		})

		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(payload))
		rec := httptest.NewRecorder()

		h.Login(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Login successful")
		assert.Contains(t, rec.Body.String(), "access-token-xyz")

		// Check Set-Cookie
		cookies := rec.Result().Cookies()
		assert.NotEmpty(t, cookies)
		var refreshCookie *http.Cookie
		for _, c := range cookies {
			if c.Name == "refresh_token" {
				refreshCookie = c
				break
			}
		}
		assert.NotNil(t, refreshCookie)
		assert.Equal(t, "refresh-token-abc", refreshCookie.Value)
		assert.True(t, refreshCookie.HttpOnly)
	})

	t.Run("Invalid Payload Syntax Returns 400", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader([]byte("{invalid-json")))
		rec := httptest.NewRecorder()

		h.Login(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid request payload")
	})

	t.Run("Invalid Credentials Returns 401", func(t *testing.T) {
		mockUsecase.loginFunc = func(ctx context.Context, req domain.LoginRequest) (domain.LoginResponse, error) {
			return domain.LoginResponse{}, errors.New("invalid credentials")
		}

		payload, _ := json.Marshal(domain.LoginRequest{
			Email:    "bad@cafe.com",
			Password: "Wrong",
		})

		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(payload))
		rec := httptest.NewRecorder()

		h.Login(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Email atau kata sandi tidak valid")
	})
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	mockUsecase := &MockAuthUsecase{}
	h := handler.NewAuthHandler(mockUsecase)

	t.Run("Happy Path - Refresh Token rotates cookie and returns new access token", func(t *testing.T) {
		mockUsecase.refreshTokenFunc = func(ctx context.Context, token string) (domain.LoginResponse, error) {
			return domain.LoginResponse{
				Token:        "new-access-token",
				RefreshToken: "rotated-refresh-token",
			}, nil
		}

		req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "existing-refresh-token"})
		rec := httptest.NewRecorder()

		h.RefreshToken(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Token refreshed")
		assert.Contains(t, rec.Body.String(), "new-access-token")
	})

	t.Run("Missing Refresh Token Cookie Returns 401", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
		rec := httptest.NewRecorder()

		h.RefreshToken(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Refresh token missing")
	})

	t.Run("Invalid Refresh Token Presentation Returns 401", func(t *testing.T) {
		mockUsecase.refreshTokenFunc = func(ctx context.Context, token string) (domain.LoginResponse, error) {
			return domain.LoginResponse{}, errors.New("token expired")
		}

		req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "expired-refresh-token"})
		rec := httptest.NewRecorder()

		h.RefreshToken(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Sesi Anda telah berakhir")
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	h := handler.NewAuthHandler(&MockAuthUsecase{})

	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	rec := httptest.NewRecorder()

	h.Logout(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Logout successful")

	cookies := rec.Result().Cookies()
	var refreshCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "refresh_token" {
			refreshCookie = c
			break
		}
	}
	assert.NotNil(t, refreshCookie)
	assert.Equal(t, -1, refreshCookie.MaxAge)
}

func TestAuthHandler_ForgotPassword(t *testing.T) {
	h := handler.NewAuthHandler(&MockAuthUsecase{})

	t.Run("Valid Email", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"email":"admin@cafe.com"}`))
		req := httptest.NewRequest("POST", "/api/v1/auth/forgot-password", body)
		rec := httptest.NewRecorder()
		h.ForgotPassword(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Instruksi reset password")
	})

	t.Run("Empty Email", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"email":""}`))
		req := httptest.NewRequest("POST", "/api/v1/auth/forgot-password", body)
		rec := httptest.NewRecorder()
		h.ForgotPassword(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthHandler_ResetPassword(t *testing.T) {
	h := handler.NewAuthHandler(&MockAuthUsecase{})

	t.Run("Valid Request", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"token":"valid-token","new_password":"NewSecret123!"}`))
		req := httptest.NewRequest("POST", "/api/v1/auth/reset-password", body)
		rec := httptest.NewRecorder()
		h.ResetPassword(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Password reset successful")
	})

	t.Run("Missing Parameters", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"token":""}`))
		req := httptest.NewRequest("POST", "/api/v1/auth/reset-password", body)
		rec := httptest.NewRecorder()
		h.ResetPassword(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

