package handler

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/thealamenthelumiere/pet-project-GO/internal/service"
)

type LoginHandler struct {
	userService service.UserService
}

func NewLoginHandler(userService service.UserService) *LoginHandler {
	return &LoginHandler{
		userService: userService,
	}
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем Basic Auth заголовок
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Authorization header required"))
		return
	}

	// Проверяем формат
	if !strings.HasPrefix(authHeader, "Basic ") {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid authorization format"))
		return
	}

	// Декодируем credentials
	encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid credentials encoding"))
		return
	}

	// Разбираем username:password
	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid credentials format"))
		return
	}

	username, password := credentials[0], credentials[1]

	// Валидируем учетные данные
	valid, err := h.userService.ValidateCredentials(username, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid credentials"))
		return
	}

	// Генерируем токен
	token, err := h.userService.GenerateToken(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to generate token"))
		return
	}

	// Возвращаем токен
	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success", "message": "Login successful"}`))
}
