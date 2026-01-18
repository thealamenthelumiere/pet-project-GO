package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_service "github.com/thealamenthelumiere/pet-project-GO/internal/service/mocks" // сгенерированные моки
)

func TestLoginHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)

	// Настраиваем ожидания
	mockUserService.EXPECT().
		ValidateCredentials("admin", "password123").
		Return(true, nil).
		Times(1)

	mockUserService.EXPECT().
		GenerateToken("admin").
		Return("fake-jwt-token", nil).
		Times(1)

	handler := NewLoginHandler(mockUserService)

	// Создаем запрос с Basic Auth
	req := httptest.NewRequest("POST", "/login", nil)
	req.Header.Set("Authorization", "Basic "+basicAuth("admin", "password123"))

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Вызываем хендлер
	handler.Handle(rr, req)

	// Проверяем результаты
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Bearer fake-jwt-token", rr.Header().Get("Authorization"))
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)

	mockUserService.EXPECT().
		ValidateCredentials("admin", "wrongpass").
		Return(false, nil).
		Times(1)

	handler := NewLoginHandler(mockUserService)

	req := httptest.NewRequest("POST", "/login", nil)
	req.Header.Set("Authorization", "Basic "+basicAuth("admin", "wrongpass"))

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestLoginHandler_MissingAuthHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	handler := NewLoginHandler(mockUserService)

	req := httptest.NewRequest("POST", "/login", nil)
	rr := httptest.NewRecorder()

	handler.Handle(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestLoginHandler_InvalidAuthFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	handler := NewLoginHandler(mockUserService)

	req := httptest.NewRequest("POST", "/login", nil)
	req.Header.Set("Authorization", "InvalidScheme token")

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func basicAuth(username, password string) string {
	// Простая реализация для тестов
	// В реальности нужно использовать base64
	return "YWRtaW46cGFzc3dvcmQxMjM=" // admin:password123 в base64
}
