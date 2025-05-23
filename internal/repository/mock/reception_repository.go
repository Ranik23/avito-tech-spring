// Code generated by MockGen. DO NOT EDIT.
// Source: /home/anton/avito-tech-spring/internal/repository/reception_repository.go
//
// Generated by this command:
//
//	mockgen --source=/home/anton/avito-tech-spring/internal/repository/reception_repository.go --destination=/home/anton/avito-tech-spring/internal/repository/mock/reception_repository.go --package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	time "time"

	domain "github.com/Ranik23/avito-tech-spring/internal/models/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockReceptionRepository is a mock of ReceptionRepository interface.
type MockReceptionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockReceptionRepositoryMockRecorder
	isgomock struct{}
}

// MockReceptionRepositoryMockRecorder is the mock recorder for MockReceptionRepository.
type MockReceptionRepositoryMockRecorder struct {
	mock *MockReceptionRepository
}

// NewMockReceptionRepository creates a new mock instance.
func NewMockReceptionRepository(ctrl *gomock.Controller) *MockReceptionRepository {
	mock := &MockReceptionRepository{ctrl: ctrl}
	mock.recorder = &MockReceptionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReceptionRepository) EXPECT() *MockReceptionRepositoryMockRecorder {
	return m.recorder
}

// CreateReception mocks base method.
func (m *MockReceptionRepository) CreateReception(ctx context.Context, pvzID string) (*domain.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReception", ctx, pvzID)
	ret0, _ := ret[0].(*domain.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateReception indicates an expected call of CreateReception.
func (mr *MockReceptionRepositoryMockRecorder) CreateReception(ctx, pvzID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReception", reflect.TypeOf((*MockReceptionRepository)(nil).CreateReception), ctx, pvzID)
}

// FindOpen mocks base method.
func (m *MockReceptionRepository) FindOpen(ctx context.Context, pvzID string) (*domain.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOpen", ctx, pvzID)
	ret0, _ := ret[0].(*domain.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOpen indicates an expected call of FindOpen.
func (mr *MockReceptionRepositoryMockRecorder) FindOpen(ctx, pvzID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOpen", reflect.TypeOf((*MockReceptionRepository)(nil).FindOpen), ctx, pvzID)
}

// GetReceptionsFiltered mocks base method.
func (m *MockReceptionRepository) GetReceptionsFiltered(ctx context.Context, pvzID string, startTime, endTime time.Time) ([]*domain.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReceptionsFiltered", ctx, pvzID, startTime, endTime)
	ret0, _ := ret[0].([]*domain.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReceptionsFiltered indicates an expected call of GetReceptionsFiltered.
func (mr *MockReceptionRepositoryMockRecorder) GetReceptionsFiltered(ctx, pvzID, startTime, endTime any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReceptionsFiltered", reflect.TypeOf((*MockReceptionRepository)(nil).GetReceptionsFiltered), ctx, pvzID, startTime, endTime)
}

// UpdateReceptionStatus mocks base method.
func (m *MockReceptionRepository) UpdateReceptionStatus(ctx context.Context, receptionID, newStatus string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateReceptionStatus", ctx, receptionID, newStatus)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateReceptionStatus indicates an expected call of UpdateReceptionStatus.
func (mr *MockReceptionRepositoryMockRecorder) UpdateReceptionStatus(ctx, receptionID, newStatus any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateReceptionStatus", reflect.TypeOf((*MockReceptionRepository)(nil).UpdateReceptionStatus), ctx, receptionID, newStatus)
}
