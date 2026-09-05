package handler

import (
	"encoding/json"
	"net/http"

	"cafe-erp-system/backend/internal/domain"
	"cafe-erp-system/backend/pkg/response"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	res, err := h.authUsecase.Login(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Set refresh token in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		HttpOnly: true,
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
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(w, "Token refreshed", map[string]interface{}{
		"token": res.Token,
		"user":  res.User,
	})
}
