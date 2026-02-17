package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/thealamenthelumiere/pet-project-GO/internal/service"
)

type VerifyHandler struct {
	userService service.UserService
}

func NewVerifyHandler(userService service.UserService) *VerifyHandler {
	return &VerifyHandler{userService: userService}
}

func (h *VerifyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	refreshToken, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	tokenPair, err := h.userService.RefreshToken(r.Context(), refreshToken)
    if err != nil {
        switch {
        case errors.Is(err, service.ErrTokenExpired):
            http.Error(w, "Refresh token expired", http.StatusUnauthorized)
        case errors.Is(err, service.ErrInvalidToken):
            http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
        default:
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenPair)
}

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format, use Bearer scheme")
	}

	token := parts[1]
	if token == "" {
		return "", errors.New("empty bearer token")
	}

	return token, nil
}