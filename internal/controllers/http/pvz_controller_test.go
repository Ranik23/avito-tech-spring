//go:build unit

package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/service"
	"github.com/Ranik23/avito-tech-spring/internal/service/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCloseLastReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock.NewMockAuthService(ctrl)
	mockPVZService := mock.NewMockPVZService(ctrl)
	svc := service.NewService(mockAuthService, mockPVZService)
	controller := NewPVZController(svc)

	tests := []struct {
		name           string
		role           string
		pvzId          string
		mockExpect     func()
		expectedStatus int
	}{
		{
			name:   "Success",
			role:   "employee",
			pvzId:  "123",
			mockExpect: func() {
				mockPVZService.EXPECT().CloseReception(gomock.Any(), "123").Return(&domain.Reception{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing pvzId",
			role:           "employee",
			pvzId:          "",
			mockExpect:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Forbidden role",
			role:   "client",
			pvzId:  "123",
			mockExpect: func() {},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:   "Service Error",
			role:   "employee",
			pvzId:  "123",
			mockExpect: func() {
				mockPVZService.EXPECT().CloseReception(gomock.Any(), "123").Return(&domain.Reception{}, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "All Receptions Closed",
			role:   "employee",
			pvzId:  "123",
			mockExpect: func() {
				mockPVZService.EXPECT().CloseReception(gomock.Any(), "123").Return(&domain.Reception{}, service.ErrAllReceptionsClosed)
			},
			expectedStatus: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.mockExpect()

			c.Params = gin.Params{{Key: "pvzId", Value: tt.pvzId}}
			c.Set("role", tt.role)

			controller.CloseLastReception(c)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock.NewMockAuthService(ctrl)
	mockPVZService := mock.NewMockPVZService(ctrl)
	svc := service.NewService(mockAuthService, mockPVZService)
	controller := NewPVZController(svc)

	tests := []struct {
		name           string
		body           string
		role           string
		mockExpect     func()
		expectedStatus int
	}{
		{
			name: "success",
			body: `{"pvzId":"123","type":"electronics"}`,
			role: "employee",
			mockExpect: func() {
				mockPVZService.EXPECT().
					AddProduct(gomock.Any(), "123", "electronics").
					Return(&domain.Product{}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json",
			body:           `{"pvzId":123}`,
			role:           "employee",
			mockExpect:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: `{"pvzId":"123","type":"fail"}`,
			role: "employee",
			mockExpect: func() {
				mockPVZService.EXPECT().
					AddProduct(gomock.Any(), "123", "fail").
					Return(nil, errors.New("error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "unauthorized role",
			body:           `{"pvzId":"123","type":"x"}`,
			role:           "client",
			mockExpect:     func() {},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/add-product", bytes.NewBufferString(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set("role", tt.role)

			controller.AddProduct(c)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCreateReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock.NewMockAuthService(ctrl)
	mockPVZService := mock.NewMockPVZService(ctrl)
	svc := service.NewService(mockAuthService, mockPVZService)
	controller := NewPVZController(svc)

	tests := []struct {
		name           string
		body           string
		role           string
		mockExpect     func()
		expectedStatus int
	}{
		{
			name: "success",
			body: `{"pvzId":"321"}`,
			role: "employee",
			mockExpect: func() {
				mockPVZService.EXPECT().
					StartReception(gomock.Any(), "321").
					Return(&domain.Reception{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid json",
			body:           `{"pvzId":123}`,
			role:           "employee",
			mockExpect:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: `{"pvzId":"fail"}`,
			role: "employee",
			mockExpect: func() {
				mockPVZService.EXPECT().
					StartReception(gomock.Any(), "fail").
					Return(nil, errors.New("err"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "unauthorized",
			body:           `{"pvzId":"x"}`,
			role:           "client",
			mockExpect:     func() {},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/create-reception", bytes.NewBufferString(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set("role", tt.role)

			controller.CreateReception(c)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestDeleteLastProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock.NewMockAuthService(ctrl)
	mockPVZService := mock.NewMockPVZService(ctrl)
	svc := service.NewService(mockAuthService, mockPVZService)
	controller := NewPVZController(svc)

	tests := []struct {
		name           string
		pvzID          string
		role           string
		mockExpect     func()
		expectedStatus int
	}{
		{
			name:  "success",
			pvzID: "123",
			role:  "employee",
			mockExpect: func() {
				mockPVZService.EXPECT().
					DeleteLastProduct(gomock.Any(), "123").
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "no pvzID",
			pvzID:          "",
			role:           "employee",
			mockExpect:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "service error",
			pvzID: "fail",
			role:  "employee",
			mockExpect: func() {
				mockPVZService.EXPECT().
					DeleteLastProduct(gomock.Any(), "fail").
					Return(errors.New("boom"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "unauthorized",
			pvzID:          "123",
			role:           "client",
			mockExpect:     func() {},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			url := "/pvz/" + tt.pvzID
			c.Params = gin.Params{gin.Param{Key: "pvzId", Value: tt.pvzID}}
			c.Request = httptest.NewRequest(http.MethodDelete, url, nil)
			c.Set("role", tt.role)

			controller.DeleteLastProduct(c)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCreatePvz(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZService := mock.NewMockPVZService(ctrl)
	mockAuthService := mock.NewMockAuthService(ctrl)

	svc := service.NewService(mockAuthService, mockPVZService)
	controller := NewPVZController(svc).(*pvzController)

	tests := []struct {
		name           string
		role           string
		requestBody    string
		mockExpect     func()
		expectedStatus int
		expectedID     string
		expectedCity   string
	}{
		{
			name:        "success",
			role:        "moderator",
			requestBody: `{"city": "Moscow"}`,
			mockExpect: func() {
				mockPVZService.EXPECT().
					CreatePVZ(gomock.Any(), "Moscow").
					Return(&domain.Pvz{ID: "pvz1", City: "Moscow"}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedID:     "pvz1",
			expectedCity:   "Moscow",
		},
		{
			name:        "invalid JSON",
			role:        "moderator",
			requestBody: `{"city": 123}`,
			mockExpect:  func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "unauthorized",
			role:        "employee",
			requestBody: `{"city": "Moscow"}`,
			mockExpect:  func() {},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:        "service error",
			role:        "moderator",
			requestBody: `{"city": "Moscow"}`,
			mockExpect: func() {
				mockPVZService.EXPECT().
					CreatePVZ(gomock.Any(), "Moscow").
					Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/pvz", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("role", tt.role)

			controller.CreatePvz(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				assert.Contains(t, w.Body.String(), `"id":"`+tt.expectedID+`"`)
				assert.Contains(t, w.Body.String(), `"city":"`+tt.expectedCity+`"`)
			}
		})
	}
}
