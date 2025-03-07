// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go
//
// Generated by this command:
//
//	mockgen -package=gone_grpc -self_package=github.com/gone-io/gone/goner/grpc -source=interface.go -destination=mock_test.go
//

// Package gone_grpc is a generated GoMock package.
package gone_grpc

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
	isgomock struct{}
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Address mocks base method.
func (m *MockClient) Address() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Address")
	ret0, _ := ret[0].(string)
	return ret0
}

// Address indicates an expected call of Address.
func (mr *MockClientMockRecorder) Address() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Address", reflect.TypeOf((*MockClient)(nil).Address))
}

// Stub mocks base method.
func (m *MockClient) Stub(conn *grpc.ClientConn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stub", conn)
}

// Stub indicates an expected call of Stub.
func (mr *MockClientMockRecorder) Stub(conn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stub", reflect.TypeOf((*MockClient)(nil).Stub), conn)
}

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
	isgomock struct{}
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// RegisterGrpcServer mocks base method.
func (m *MockService) RegisterGrpcServer(server *grpc.Server) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterGrpcServer", server)
}

// RegisterGrpcServer indicates an expected call of RegisterGrpcServer.
func (mr *MockServiceMockRecorder) RegisterGrpcServer(server any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterGrpcServer", reflect.TypeOf((*MockService)(nil).RegisterGrpcServer), server)
}
