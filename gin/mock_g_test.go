// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gone-io/goner/g (interfaces: Cmux,Tracer)
//
// Generated by this command:
//
//	mockgen -package=gin -destination=mock_g_test.go github.com/gone-io/goner/g Cmux,Tracer
//

// Package gin is a generated GoMock package.
package gin

import (
	net "net"
	reflect "reflect"

	g "github.com/gone-io/goner/g"
	gomock "go.uber.org/mock/gomock"
)

// MockCmux is a mock of Cmux interface.
type MockCmux struct {
	ctrl     *gomock.Controller
	recorder *MockCmuxMockRecorder
	isgomock struct{}
}

// MockCmuxMockRecorder is the mock recorder for MockCmux.
type MockCmuxMockRecorder struct {
	mock *MockCmux
}

// NewMockCmux creates a new mock instance.
func NewMockCmux(ctrl *gomock.Controller) *MockCmux {
	mock := &MockCmux{ctrl: ctrl}
	mock.recorder = &MockCmuxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCmux) EXPECT() *MockCmuxMockRecorder {
	return m.recorder
}

// GetAddress mocks base method.
func (m *MockCmux) GetAddress() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddress")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAddress indicates an expected call of GetAddress.
func (mr *MockCmuxMockRecorder) GetAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddress", reflect.TypeOf((*MockCmux)(nil).GetAddress))
}

// MatchFor mocks base method.
func (m *MockCmux) MatchFor(arg0 g.ProtocolType) net.Listener {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MatchFor", arg0)
	ret0, _ := ret[0].(net.Listener)
	return ret0
}

// MatchFor indicates an expected call of MatchFor.
func (mr *MockCmuxMockRecorder) MatchFor(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MatchFor", reflect.TypeOf((*MockCmux)(nil).MatchFor), arg0)
}

// MockTracer is a mock of Tracer interface.
type MockTracer struct {
	ctrl     *gomock.Controller
	recorder *MockTracerMockRecorder
	isgomock struct{}
}

// MockTracerMockRecorder is the mock recorder for MockTracer.
type MockTracerMockRecorder struct {
	mock *MockTracer
}

// NewMockTracer creates a new mock instance.
func NewMockTracer(ctrl *gomock.Controller) *MockTracer {
	mock := &MockTracer{ctrl: ctrl}
	mock.recorder = &MockTracerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTracer) EXPECT() *MockTracerMockRecorder {
	return m.recorder
}

// GetTraceId mocks base method.
func (m *MockTracer) GetTraceId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTraceId")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTraceId indicates an expected call of GetTraceId.
func (mr *MockTracerMockRecorder) GetTraceId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTraceId", reflect.TypeOf((*MockTracer)(nil).GetTraceId))
}

// Go mocks base method.
func (m *MockTracer) Go(fn func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Go", fn)
}

// Go indicates an expected call of Go.
func (mr *MockTracerMockRecorder) Go(fn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Go", reflect.TypeOf((*MockTracer)(nil).Go), fn)
}

// SetTraceId mocks base method.
func (m *MockTracer) SetTraceId(traceId string, fn func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTraceId", traceId, fn)
}

// SetTraceId indicates an expected call of SetTraceId.
func (mr *MockTracerMockRecorder) SetTraceId(traceId, fn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTraceId", reflect.TypeOf((*MockTracer)(nil).SetTraceId), traceId, fn)
}
