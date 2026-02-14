package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/thealamenthelumiere/pet-project-GO/internal/service"
)

type LoginHandler struct {
	userService service.UserService
}

func NewLoginHandler(userService service.UserService) *LoginHandler {
	return &LoginHandler{userService: userService}
}

// Handle обрабатывает POST /login с Basic Auth.
func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем username и password из заголовка Authorization: Basic ...
	username, password, err := extractBasicAuth(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Вызываем сервис — он проверяет credentials и возвращает пару токенов
	tokenPair, err := h.userService.Login(username, password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Успех: возвращаем access_token и refresh_token в JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenPair)
}

// extractBasicAuth извлекает логин и пароль из заголовка Authorization: Basic <base64>.
// Возвращает username, password и ошибку с понятным текстом для клиента.
func extractBasicAuth(r *http.Request) (string, string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", "", errors.New("authorization header required")
	}

	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", "", errors.New("invalid authorization format, use Basic scheme")
	}

	encoded := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", errors.New("invalid base64 encoding")
	}

	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 {
		return "", "", errors.New("invalid credentials format, expected username:password")
	}

	return creds[0], creds[1], nil
}