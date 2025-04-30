package grpc

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	gMock "github.com/gone-io/goner/g/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"net"
	"testing"
)

type mockService struct {
	registered bool
}

func (m *mockService) RegisterGrpcServer(*grpc.Server) {
	m.registered = true
}

func Test_server_cmux_integration(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	cMuxServer := gMock.NewMockCmux(controller)
	listener := NewMockListener(controller)

	mockAddr := NewMockAddr(controller)
	mockAddr.EXPECT().String().Return("127.0.0.1:8080").AnyTimes()

	conn := NewMockConn(controller)
	conn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()

	listener.EXPECT().Addr().Return(mockAddr).AnyTimes()
	listener.EXPECT().Accept().Return(conn, nil).AnyTimes()
	listener.EXPECT().Close().Return(nil).AnyTimes()
	cMuxServer.EXPECT().MatchFor(g.GRPC).Return(listener)

	s := &server{
		logger:       gone.GetDefaultLogger(),
		grpcServices: []Service{&mockService{}},
		cMuxServer:   cMuxServer,
	}

	s.Init()
	assert.NotNil(t, s.listener)
}

func Test_server_service_discovery(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	registry := gMock.NewMockServiceRegistry(controller)
	registry.EXPECT().Register(gomock.Any()).Return(nil)
	registry.EXPECT().Deregister(gomock.Any()).Return(nil)

	listener := NewMockListener(controller)
	addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
	listener.EXPECT().Addr().Return(addr).AnyTimes()

	s := &server{
		logger:           gone.GetDefaultLogger(),
		grpcServices:     []Service{&mockService{}},
		registry:         registry,
		serviceName:      "test-service",
		serviceUseSubNet: "127.0.0.1/24",
		createListener: func(host string, port int) net.Listener {
			return listener
		},
		getLocalIps: func() []net.IP {
			return []net.IP{net.ParseIP("127.0.0.1")}
		},
	}

	s.Init()

	undo := s.regService()
	assert.NotNil(t, undo)

	err := undo()
	assert.NoError(t, err)
}

func Test_server_service_discovery_invalid_subnet(t *testing.T) {
	s := &server{
		logger:           gone.GetDefaultLogger(),
		grpcServices:     []Service{&mockService{}},
		registry:         gMock.NewMockServiceRegistry(gomock.NewController(t)),
		serviceName:      "test-service",
		serviceUseSubNet: "invalid-subnet",
		getLocalIps: func() []net.IP {
			return []net.IP{net.ParseIP("127.0.0.1")}
		},
	}

	assert.Panics(t, func() {
		s.regService()
	})
}

func Test_server_service_discovery_empty_service_name(t *testing.T) {
	s := &server{
		logger:           gone.GetDefaultLogger(),
		grpcServices:     []Service{&mockService{}},
		registry:         gMock.NewMockServiceRegistry(gomock.NewController(t)),
		serviceUseSubNet: "127.0.0.1/24",
	}

	assert.Panics(t, func() {
		s.regService()
	})
}
