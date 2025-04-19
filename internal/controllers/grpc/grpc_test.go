//go:build unit

package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/Ranik23/avito-tech-spring/api/proto/gen/pvz_v1"
	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/service/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetPVZList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock.NewMockService(ctrl)

	server := NewPVZServer(mockService)

	tests := []struct {
		name             string
		mockExpect       func()
		request          *pvz_v1.GetPVZListRequest
		expectedResponse *pvz_v1.GetPVZListResponse
		expectedError    error
	}{
		{
			name: "success",
			mockExpect: func() {
				mockService.EXPECT().
					GetPVZList(gomock.Any()).
					Return([]domain.Pvz{{ID: "1", City: "Moscow"}}, nil)
			},
			request: &pvz_v1.GetPVZListRequest{},
			expectedResponse: &pvz_v1.GetPVZListResponse{
				Pvzs: []*pvz_v1.PVZ{
					{
						Id:   "1",
						City: "Moscow",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "service error",
			mockExpect: func() {
				mockService.EXPECT().
					GetPVZList(gomock.Any()).
					Return(nil, errors.New("internal service error"))
			},
			request:          &pvz_v1.GetPVZListRequest{},
			expectedResponse: nil,
			expectedError:    errors.New("internal service error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			resp, err := server.GetPVZList(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				for i := 0; i < len(tt.expectedResponse.Pvzs); i++ {
					assert.Equal(t, tt.expectedResponse.Pvzs[i].City, resp.Pvzs[i].City)
					assert.Equal(t, tt.expectedResponse.Pvzs[i].Id, resp.Pvzs[i].Id)
				}
			}
		})
	}
}
