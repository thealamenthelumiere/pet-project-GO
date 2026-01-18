package handler

import (
	"net/http"
	"strings"

	"github.com/thealamenthelumiere/pet-project-GO/internal/service"
)

type VerifyHandler struct {
	userService service.UserService
}

func NewVerifyHandler(userService service.UserService) *VerifyHandler {
	return &VerifyHandler{
		userService: userService,
	}
}

func (h *VerifyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод (можно использовать POST или GET, обычно POST для обновления токена)
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем Bearer токен
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Authorization header required"))
		return
	}

	// Проверяем формат Bearer
	if !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid authorization format. Use Bearer scheme"))
		return
	}

	// Извлекаем токен
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Обновляем токен
	newToken, err := h.userService.RefreshToken(token)
	if err != nil {
		// Проверяем тип ошибки для более точного ответа
		if err.Error() == "token expired" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Token expired"))
			return
		}
		if err.Error() == "invalid token" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid token"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	// Возвращаем новый токен
	w.Header().Set("Authorization", "Bearer "+newToken)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success", "message": "Token refreshed"}`))
}
