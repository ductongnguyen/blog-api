// Code generated by MockGen. DO NOT EDIT.
// Source: cache.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/ductong169z/shorten-url/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockRedisRepository is a mock of RedisRepository interface.
type MockRedisRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRedisRepositoryMockRecorder
}

// MockRedisRepositoryMockRecorder is the mock recorder for MockRedisRepository.
type MockRedisRepositoryMockRecorder struct {
	mock *MockRedisRepository
}

// NewMockRedisRepository creates a new mock instance.
func NewMockRedisRepository(ctrl *gomock.Controller) *MockRedisRepository {
	mock := &MockRedisRepository{ctrl: ctrl}
	mock.recorder = &MockRedisRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedisRepository) EXPECT() *MockRedisRepositoryMockRecorder {
	return m.recorder
}

// GetUserByIDCtx mocks base method.
func (m *MockRedisRepository) GetUserByIDCtx(ctx context.Context, key string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByIDCtx", ctx, key)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByIDCtx indicates an expected call of GetUserByIDCtx.
func (mr *MockRedisRepositoryMockRecorder) GetUserByIDCtx(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByIDCtx", reflect.TypeOf((*MockRedisRepository)(nil).GetUserByIDCtx), ctx, key)
}

// SetUserByIDCtx mocks base method.
func (m *MockRedisRepository) SetUserByIDCtx(ctx context.Context, key string, user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUserByIDCtx", ctx, key, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUserByIDCtx indicates an expected call of SetUserByIDCtx.
func (mr *MockRedisRepositoryMockRecorder) SetUserByIDCtx(ctx, key, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUserByIDCtx", reflect.TypeOf((*MockRedisRepository)(nil).SetUserByIDCtx), ctx, key, user)
}
