package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/thealamenthelumiere/pet-project-GO/internal/service"
	mock_service "github.com/thealamenthelumiere/pet-project-GO/internal/service/mocks"
)

func TestVerifyHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		validRefreshToken   = "valid-refresh-token"
		expiredRefreshToken = "expired-refresh-token"
		invalidRefreshToken = "invalid-refresh-token"
		newAccessToken      = "new-access-token"
		newRefreshToken     = "new-refresh-token"
	)

	tests := []struct {
		name           string
		authHeader     string
		setupMock      func(*mock_service.MockUserService)
		expectedStatus int
		expectedBody   func(*testing.T, []byte)
	}{
		{
			name:       "successful token refresh",
			authHeader: "Bearer " + validRefreshToken,
			setupMock: func(m *mock_service.MockUserService) {
				m.EXPECT().
					RefreshToken(gomock.Any(), validRefreshToken).
					Return(&service.TokenPair{
						AccessToken:  newAccessToken,
						RefreshToken: newRefreshToken,
					}, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody: func(t *testing.T, body []byte) {
				var tp service.TokenPair
				err := json.Unmarshal(body, &tp)
				require.NoError(t, err)
				assert.Equal(t, newAccessToken, tp.AccessToken)
				assert.Equal(t, newRefreshToken, tp.RefreshToken)
			},
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			setupMock:      func(m *mock_service.MockUserService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "authorization header required")
			},
		},
		{
			name:           "invalid auth scheme (not Bearer)",
			authHeader:     "Basic xxx",
			setupMock:      func(m *mock_service.MockUserService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "invalid authorization header format, use Bearer scheme")
			},
		},
		{
			name:           "malformed token (empty after Bearer)",
			authHeader:     "Bearer ",
			setupMock:      func(m *mock_service.MockUserService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "empty bearer token")
			},
		},
		{
			name:       "token expired",
			authHeader: "Bearer " + expiredRefreshToken,
			setupMock: func(m *mock_service.MockUserService) {
				m.EXPECT().
					RefreshToken(gomock.Any(), expiredRefreshToken).
					Return(nil, service.ErrTokenExpired). 
					Times(1)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "Refresh token expired")
			},
		},
		{
			name:       "invalid token",
			authHeader: "Bearer " + invalidRefreshToken,
			setupMock: func(m *mock_service.MockUserService) {
				m.EXPECT().
					RefreshToken(gomock.Any(), invalidRefreshToken).
					Return(nil, service.ErrInvalidToken). 
					Times(1)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "Invalid refresh token")
			},
		},
		{
			name:       "internal server error",
			authHeader: "Bearer " + validRefreshToken,
			setupMock: func(m *mock_service.MockUserService) {
				m.EXPECT().
					RefreshToken(gomock.Any(), validRefreshToken).
					Return(nil, errors.New("unexpected db error")).
					Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "Internal server error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := mock_service.NewMockUserService(ctrl)
			tt.setupMock(mockUserService)

			handler := NewVerifyHandler(mockUserService)

			req := httptest.NewRequest(http.MethodPost, "/verify", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handler.Handle(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != nil {
				tt.expectedBody(t, rr.Body.Bytes())
			}
		})
	}
}