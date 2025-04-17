// Code generated by MockGen. DO NOT EDIT.
// Source: pvz_repository.go
//
// Generated by this command:
//
//	mockgen --source=pvz_repository.go --destination=mock/pvz_repository.go --package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	domain "github.com/Ranik23/avito-tech-spring/internal/models/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockPvzRepository is a mock of PvzRepository interface.
type MockPvzRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPvzRepositoryMockRecorder
	isgomock struct{}
}

// MockPvzRepositoryMockRecorder is the mock recorder for MockPvzRepository.
type MockPvzRepositoryMockRecorder struct {
	mock *MockPvzRepository
}

// NewMockPvzRepository creates a new mock instance.
func NewMockPvzRepository(ctrl *gomock.Controller) *MockPvzRepository {
	mock := &MockPvzRepository{ctrl: ctrl}
	mock.recorder = &MockPvzRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPvzRepository) EXPECT() *MockPvzRepositoryMockRecorder {
	return m.recorder
}

// CreatePVZ mocks base method.
func (m *MockPvzRepository) CreatePVZ(ctx context.Context, city string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePVZ", ctx, city)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePVZ indicates an expected call of CreatePVZ.
func (mr *MockPvzRepositoryMockRecorder) CreatePVZ(ctx, city any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePVZ", reflect.TypeOf((*MockPvzRepository)(nil).CreatePVZ), ctx, city)
}

// GetPVZ mocks base method.
func (m *MockPvzRepository) GetPVZ(ctx context.Context, id string) (*domain.Pvz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPVZ", ctx, id)
	ret0, _ := ret[0].(*domain.Pvz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPVZ indicates an expected call of GetPVZ.
func (mr *MockPvzRepositoryMockRecorder) GetPVZ(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPVZ", reflect.TypeOf((*MockPvzRepository)(nil).GetPVZ), ctx, id)
}

// GetPVZS mocks base method.
func (m *MockPvzRepository) GetPVZS(ctx context.Context, offset, limit int) ([]domain.Pvz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPVZS", ctx, offset, limit)
	ret0, _ := ret[0].([]domain.Pvz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPVZS indicates an expected call of GetPVZS.
func (mr *MockPvzRepositoryMockRecorder) GetPVZS(ctx, offset, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPVZS", reflect.TypeOf((*MockPvzRepository)(nil).GetPVZS), ctx, offset, limit)
}
