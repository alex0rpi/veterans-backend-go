package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"veterans-go-chi-server/internal/auth"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	jwtManager    *auth.JWTManager
	adminUsername string
	adminPasswordHash string
}

func NewAuthHandler(jwtManager *auth.JWTManager, adminUsername, adminPasswordHash string) *AuthHandler {
	return &AuthHandler{
		jwtManager:    jwtManager,
		adminUsername: adminUsername,
		adminPasswordHash: adminPasswordHash,
	}
}

// Login to generate tokens (access)
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req auth.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	passwordErr := bcrypt.CompareHashAndPassword(
		[]byte(h.adminPasswordHash),
		[]byte(req.Password),
	)

	if req.Username != h.adminUsername || passwordErr != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generar tokens
	accessToken, err := h.jwtManager.GenerateToken(auth.TokenTypeAccess, 15*time.Minute)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(auth.LoginTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   900, // 15 minutos
	})
}
