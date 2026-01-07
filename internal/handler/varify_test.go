package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/thealamenthelumiere/pet-project-GO/internal/service"
	mock_service "github.com/thealamenthelumiere/pet-project-GO/internal/service/mocks"
)

func TestVerifyHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)

	mockUserService.EXPECT().
		RefreshToken("valid-token").
		Return("new-valid-token", nil).
		Times(1)

	handler := NewVerifyHandler(mockUserService)

	req := httptest.NewRequest("GET", "/verify", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Bearer new-valid-token", rr.Header().Get("Authorization"))
}

func TestVerifyHandler_ExpiredToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)

	mockUserService.EXPECT().
		RefreshToken("expired-token").
		Return("", service.NewTokenError("token expired")).
		Times(1)

	handler := NewVerifyHandler(mockUserService)

	req := httptest.NewRequest("GET", "/verify", nil)
	req.Header.Set("Authorization", "Bearer expired-token")

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestVerifyHandler_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)

	mockUserService.EXPECT().
		RefreshToken("invalid-token").
		Return("", service.NewTokenError("invalid token")).
		Times(1)

	handler := NewVerifyHandler(mockUserService)

	req := httptest.NewRequest("GET", "/verify", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestVerifyHandler_MissingAuthHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	handler := NewVerifyHandler(mockUserService)

	req := httptest.NewRequest("GET", "/verify", nil)
	rr := httptest.NewRecorder()

	handler.Handle(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestVerifyHandler_InvalidAuthScheme(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	handler := NewVerifyHandler(mockUserService)

	req := httptest.NewRequest("GET", "/verify", nil)
	req.Header.Set("Authorization", "Basic username:password")

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
