// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/server/repository/app_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	model "go-starter-kit/internal/server/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAppRepository is a mock of AppRepository interface.
type MockAppRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAppRepositoryMockRecorder
}

// MockAppRepositoryMockRecorder is the mock recorder for MockAppRepository.
type MockAppRepositoryMockRecorder struct {
	mock *MockAppRepository
}

// NewMockAppRepository creates a new mock instance.
func NewMockAppRepository(ctrl *gomock.Controller) *MockAppRepository {
	mock := &MockAppRepository{ctrl: ctrl}
	mock.recorder = &MockAppRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppRepository) EXPECT() *MockAppRepositoryMockRecorder {
	return m.recorder
}

// CreateApp mocks base method.
func (m *MockAppRepository) CreateApp(ctx context.Context, app model.App) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateApp", ctx, app)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateApp indicates an expected call of CreateApp.
func (mr *MockAppRepositoryMockRecorder) CreateApp(ctx, app interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateApp", reflect.TypeOf((*MockAppRepository)(nil).CreateApp), ctx, app)
}

// GetByAppID mocks base method.
func (m *MockAppRepository) GetByAppID(ctx context.Context, appID string) (*model.App, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByAppID", ctx, appID)
	ret0, _ := ret[0].(*model.App)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByAppID indicates an expected call of GetByAppID.
func (mr *MockAppRepositoryMockRecorder) GetByAppID(ctx, appID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByAppID", reflect.TypeOf((*MockAppRepository)(nil).GetByAppID), ctx, appID)
}