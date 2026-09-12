package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"cafe-erp-system/backend/internal/delivery/http/middleware"
	"cafe-erp-system/backend/internal/domain"
	"cafe-erp-system/backend/pkg/logger"
	"cafe-erp-system/backend/pkg/response"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func isSecureRequest(r *http.Request) bool {
	return r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" || os.Getenv("ENV") == "production"
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	ip := middleware.GetClientIP(r)

	res, err := h.authUsecase.Login(r.Context(), req)
	if err != nil {
		logger.Log.Warn().
			Str("event", "AUTH_LOGIN_FAILED").
			Str("email", req.Email).
			Str("ip", ip).
			Str("user_agent", r.UserAgent()).
			Msg("Failed user login attempt")

		response.Error(w, http.StatusUnauthorized, "Email atau kata sandi tidak valid")
		return
	}

	logger.Log.Info().
		Str("event", "AUTH_LOGIN_SUCCESS").
		Str("user_id", res.User.ID.String()).
		Str("username", res.User.Username).
		Str("email", res.User.Email).
		Str("ip", ip).
		Msg("User authenticated successfully")

	// Set refresh token in secure cookie conforming to OWASP ASVS V3
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		HttpOnly: true,
		Secure:   isSecureRequest(r),
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60, // 7 days
	})

	response.Success(w, "Login successful", map[string]interface{}{
		"token": res.Token,
		"user":  res.User,
	})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Refresh token missing")
		return
	}

	res, err := h.authUsecase.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		logger.Log.Warn().
			Str("event", "AUTH_REFRESH_FAILED").
			Str("ip", middleware.GetClientIP(r)).
			Msg("Invalid refresh token presentation")

		response.Error(w, http.StatusUnauthorized, "Sesi Anda telah berakhir, silakan login kembali")
		return
	}

	// Update refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		HttpOnly: true,
		Secure:   isSecureRequest(r),
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
	})

	response.Success(w, "Token refreshed", map[string]interface{}{
		"token": res.Token,
		"user":  res.User,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Invalidate refresh token cookie immediately
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   isSecureRequest(r),
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	})

	logger.Log.Info().
		Str("event", "AUTH_LOGOUT").
		Str("ip", middleware.GetClientIP(r)).
		Msg("User logged out and session cookie revoked")

	response.Success(w, "Logout successful", nil)
}
