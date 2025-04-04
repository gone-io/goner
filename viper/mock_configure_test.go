// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gone-io/gone/v2 (interfaces: Configure)
//
// Generated by this command:
//
//	mockgen -destination=mock_configure_test.go -package=viper github.com/gone-io/gone/v2 Configure
//

// Package viper is a generated GoMock package.
package viper

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockConfigure is a mock of Configure interface.
type MockConfigure struct {
	ctrl     *gomock.Controller
	recorder *MockConfigureMockRecorder
	isgomock struct{}
}

// MockConfigureMockRecorder is the mock recorder for MockConfigure.
type MockConfigureMockRecorder struct {
	mock *MockConfigure
}

// NewMockConfigure creates a new mock instance.
func NewMockConfigure(ctrl *gomock.Controller) *MockConfigure {
	mock := &MockConfigure{ctrl: ctrl}
	mock.recorder = &MockConfigureMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigure) EXPECT() *MockConfigureMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockConfigure) Get(key string, v any, defaultVal string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key, v, defaultVal)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockConfigureMockRecorder) Get(key, v, defaultVal any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockConfigure)(nil).Get), key, v, defaultVal)
}
