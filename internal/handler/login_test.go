package handler

import (
	"encoding/base64"
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


func encodeBasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func TestLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	
	const (
		validUsername   = "admin"
		validPassword   = "password123"
		invalidPassword = "wrongpass"
		accessToken     = "fake-access-token"
		refreshToken    = "fake-refresh-token"
	)

	
	tests := []struct {
		name           string
		authHeader     string                     
		setupMock      func(*mock_service.MockUserService) 
		expectedStatus int                        
		expectedBody   func(*testing.T, []byte)   
	}{
		{
			name:       "successful login",
			authHeader: "Basic " + encodeBasicAuth(validUsername, validPassword),
			setupMock: func(m *mock_service.MockUserService) {
				m.EXPECT().
					Login(validUsername, validPassword).
					Return(&service.TokenPair{
						AccessToken:  accessToken,
						RefreshToken: refreshToken,
					}, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody: func(t *testing.T, body []byte) {
				var tp service.TokenPair
				err := json.Unmarshal(body, &tp)
				require.NoError(t, err)
				assert.Equal(t, accessToken, tp.AccessToken)
				assert.Equal(t, refreshToken, tp.RefreshToken)
			},
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			setupMock:      func(m *mock_service.MockUserService) {}, // нет вызовов
			expectedStatus: http.StatusUnauthorized,
			expectedBody: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "authorization header required")
			},
		},
		{
			name:           "invalid auth scheme",
			authHeader:     "Bearer token",
			setupMock:      func(m *mock_service.MockUserService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "invalid authorization format, use Basic scheme")
			},
		},
		{
			name:           "malformed base64",
			authHeader:     "Basic not-base64!",
			setupMock:      func(m *mock_service.MockUserService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "invalid base64 encoding")
			},
		},
		{
			name:           "invalid credentials format (no colon)",
			authHeader:     "Basic " + base64.StdEncoding.EncodeToString([]byte("username-without-password")),
			setupMock:      func(m *mock_service.MockUserService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "invalid credentials format, expected username:password")
			},
		},
		{
			name:       "invalid credentials (wrong password)",
			authHeader: "Basic " + encodeBasicAuth(validUsername, invalidPassword),
			setupMock: func(m *mock_service.MockUserService) {
				m.EXPECT().
					Login(validUsername, invalidPassword).
					Return(nil, service.ErrInvalidCredentials).
					Times(1)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "Invalid username or password")
			},
		},
		{
			name:       "user not found (same as invalid credentials for security)",
			authHeader: "Basic " + encodeBasicAuth("unknown", "pass"),
			setupMock: func(m *mock_service.MockUserService) {
				m.EXPECT().
					Login("unknown", "pass").
					Return(nil, service.ErrInvalidCredentials).
					Times(1)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "Invalid username or password")
			},
		},
		{
			name:       "internal server error",
			authHeader: "Basic " + encodeBasicAuth(validUsername, validPassword),
			setupMock: func(m *mock_service.MockUserService) {
				m.EXPECT().
					Login(validUsername, validPassword).
					Return(nil, errors.New("database connection failed")).
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

			handler := NewLoginHandler(mockUserService)

			req := httptest.NewRequest(http.MethodPost, "/login", nil)
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