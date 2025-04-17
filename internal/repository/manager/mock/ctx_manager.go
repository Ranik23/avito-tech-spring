// Code generated by MockGen. DO NOT EDIT.
// Source: ctx_manager.go
//
// Generated by this command:
//
//	mockgen --source=ctx_manager.go --destination=mock/ctx_manager.go --package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	manager "github.com/Ranik23/avito-tech-spring/internal/repository/manager"
	gomock "go.uber.org/mock/gomock"
)

// MockCtxManager is a mock of CtxManager interface.
type MockCtxManager struct {
	ctrl     *gomock.Controller
	recorder *MockCtxManagerMockRecorder
	isgomock struct{}
}

// MockCtxManagerMockRecorder is the mock recorder for MockCtxManager.
type MockCtxManagerMockRecorder struct {
	mock *MockCtxManager
}

// NewMockCtxManager creates a new mock instance.
func NewMockCtxManager(ctrl *gomock.Controller) *MockCtxManager {
	mock := &MockCtxManager{ctrl: ctrl}
	mock.recorder = &MockCtxManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCtxManager) EXPECT() *MockCtxManagerMockRecorder {
	return m.recorder
}

// ByKey mocks base method.
func (m *MockCtxManager) ByKey(arg0 context.Context, arg1 manager.CtxKey) manager.Transaction {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByKey", arg0, arg1)
	ret0, _ := ret[0].(manager.Transaction)
	return ret0
}

// ByKey indicates an expected call of ByKey.
func (mr *MockCtxManagerMockRecorder) ByKey(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByKey", reflect.TypeOf((*MockCtxManager)(nil).ByKey), arg0, arg1)
}

// CtxKey mocks base method.
func (m *MockCtxManager) CtxKey() manager.CtxKey {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CtxKey")
	ret0, _ := ret[0].(manager.CtxKey)
	return ret0
}

// CtxKey indicates an expected call of CtxKey.
func (mr *MockCtxManagerMockRecorder) CtxKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CtxKey", reflect.TypeOf((*MockCtxManager)(nil).CtxKey))
}

// Default mocks base method.
func (m *MockCtxManager) Default(arg0 context.Context) manager.Transaction {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Default", arg0)
	ret0, _ := ret[0].(manager.Transaction)
	return ret0
}

// Default indicates an expected call of Default.
func (mr *MockCtxManagerMockRecorder) Default(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Default", reflect.TypeOf((*MockCtxManager)(nil).Default), arg0)
}
