// Code generated by MockGen. DO NOT EDIT.
// Source: product_repository.go
//
// Generated by this command:
//
//	mockgen --source=product_repository.go --destination=mock/product_repository.go --package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	domain "github.com/Ranik23/avito-tech-spring/internal/models/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockProductRepository is a mock of ProductRepository interface.
type MockProductRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProductRepositoryMockRecorder
	isgomock struct{}
}

// MockProductRepositoryMockRecorder is the mock recorder for MockProductRepository.
type MockProductRepositoryMockRecorder struct {
	mock *MockProductRepository
}

// NewMockProductRepository creates a new mock instance.
func NewMockProductRepository(ctrl *gomock.Controller) *MockProductRepository {
	mock := &MockProductRepository{ctrl: ctrl}
	mock.recorder = &MockProductRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductRepository) EXPECT() *MockProductRepositoryMockRecorder {
	return m.recorder
}

// CreateProduct mocks base method.
func (m *MockProductRepository) CreateProduct(ctx context.Context, productType, receptionID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProduct", ctx, productType, receptionID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProduct indicates an expected call of CreateProduct.
func (mr *MockProductRepositoryMockRecorder) CreateProduct(ctx, productType, receptionID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProduct", reflect.TypeOf((*MockProductRepository)(nil).CreateProduct), ctx, productType, receptionID)
}

// DeleteProduct mocks base method.
func (m *MockProductRepository) DeleteProduct(ctx context.Context, productID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProduct", ctx, productID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProduct indicates an expected call of DeleteProduct.
func (mr *MockProductRepositoryMockRecorder) DeleteProduct(ctx, productID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProduct", reflect.TypeOf((*MockProductRepository)(nil).DeleteProduct), ctx, productID)
}

// FindTheLastProduct mocks base method.
func (m *MockProductRepository) FindTheLastProduct(ctx context.Context, pvzID string) (*domain.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTheLastProduct", ctx, pvzID)
	ret0, _ := ret[0].(*domain.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindTheLastProduct indicates an expected call of FindTheLastProduct.
func (mr *MockProductRepositoryMockRecorder) FindTheLastProduct(ctx, pvzID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTheLastProduct", reflect.TypeOf((*MockProductRepository)(nil).FindTheLastProduct), ctx, pvzID)
}

// GetProducts mocks base method.
func (m *MockProductRepository) GetProducts(ctx context.Context, receptionID string) ([]domain.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProducts", ctx, receptionID)
	ret0, _ := ret[0].([]domain.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProducts indicates an expected call of GetProducts.
func (mr *MockProductRepositoryMockRecorder) GetProducts(ctx, receptionID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProducts", reflect.TypeOf((*MockProductRepository)(nil).GetProducts), ctx, receptionID)
}
