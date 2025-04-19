//go:build unit

package http

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ranik23/avito-tech-spring/internal/service"
	"github.com/Ranik23/avito-tech-spring/internal/service/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDummyLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock.NewMockAuthService(ctrl)
	mockPVZService := mock.NewMockPVZService(ctrl)

	svc := service.NewService(mockAuthService, mockPVZService)
	logger := slog.Default()
	controller := NewAuthController(svc, logger)

	tests := []struct {
		name           string
		requestBody    string
		mockExpect     func()
		expectedStatus int
		expectedToken  string
	}{
		{
			name:        "Success",
			requestBody: `{"role":"client"}`,
			mockExpect: func() {
				mockAuthService.EXPECT().
					DummyLogin(gomock.Any(), "client").
					Return("valid-token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedToken:  "valid-token",
		},
		{
			name:        "Invalid JSON",
			requestBody: `{"role":123}`, 
			mockExpect:  func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Service Error",
			requestBody: `{"role":"moderator"}`,
			mockExpect: func() {
				mockAuthService.EXPECT().
					DummyLogin(gomock.Any(), "moderator").
					Return("", errors.New("token gen failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/dummy-login", bytes.NewBufferString(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			controller.DummyLogin(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, w.Body.String(), tt.expectedToken)
			}
		})
	}
}


func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock.NewMockAuthService(ctrl)
	mockPVZService := mock.NewMockPVZService(ctrl)

	svc := service.NewService(mockAuthService, mockPVZService)
	logger := slog.Default()
	controller := NewAuthController(svc, logger)

	tests := []struct {
		name           string
		requestBody    string
		mockExpect     func()
		expectedStatus int
		expectedToken  string
	}{
		{
			name:        "Success",
			requestBody: `{"email":"test@example.com", "password":"1234"}`,
			mockExpect: func() {
				mockAuthService.EXPECT().
					Login(gomock.Any(), "test@example.com", "1234").
					Return("jwt-token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedToken:  "jwt-token",
		},
		{
			name:        "Invalid JSON",
			requestBody: `{"email":123}`,
			mockExpect:  func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Invalid Credentials",
			requestBody: `{"email":"wrong@example.com", "password":"wrong"}`,
			mockExpect: func() {
				mockAuthService.EXPECT().
					Login(gomock.Any(), "wrong@example.com", "wrong").
					Return("", service.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:        "Internal Error",
			requestBody: `{"email":"fail@example.com", "password":"1234"}`,
			mockExpect: func() {
				mockAuthService.EXPECT().
					Login(gomock.Any(), "fail@example.com", "1234").
					Return("", errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			controller.Login(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, w.Body.String(), tt.expectedToken)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock.NewMockAuthService(ctrl)
	mockPVZService := mock.NewMockPVZService(ctrl)

	svc := service.NewService(mockAuthService, mockPVZService)
	logger := slog.Default()
	controller := NewAuthController(svc, logger)

	tests := []struct {
		name           string
		requestBody    string
		mockExpect     func()
		expectedStatus int
		expectedID     string
	}{
		{
			name:        "Success",
			requestBody: `{"email":"user@example.com", "password":"secret", "role":"client"}`,
			mockExpect: func() {
				mockAuthService.EXPECT().
					Register(gomock.Any(), "user@example.com", "secret", "client").
					Return("generated-id", nil)
			},
			expectedStatus: http.StatusCreated,
			expectedID:     "generated-id",
		},
		{
			name:        "Bad Request",
			requestBody: `{"email":123}`,
			mockExpect:  func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Internal Error",
			requestBody: `{"email":"fail@example.com", "password":"1234", "role":"admin"}`,
			mockExpect: func() {
				mockAuthService.EXPECT().
					Register(gomock.Any(), "fail@example.com", "1234", "admin").
					Return("", errors.New("creation failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			controller.Register(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusCreated {
				assert.Contains(t, w.Body.String(), tt.expectedID)
			}
		})
	}
}

