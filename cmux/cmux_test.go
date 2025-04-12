package cmux

import (
	"net"
	"sync"
	"testing"
	"time"

	mock "github.com/gone-io/gone/mock/v2"
	"github.com/gone-io/goner/g"
	gMock "github.com/gone-io/goner/g/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// TestServer_Init 测试服务器初始化
func TestServer_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(ctrl)
	mockRegistry := gMock.NewMockServiceRegistry(ctrl)
	mockTracer := gMock.NewMockTracer(ctrl)
	mockListener := NewMockListener(ctrl)

	// 创建server实例
	s := &server{
		logger:   mockLogger,
		registry: mockRegistry,
		tracer:   mockTracer,
		network:  "tcp",
		address:  "localhost:8080",
		host:     "localhost",
		port:     8080,
		lock:     sync.Mutex{},
		listen: func(network, address string) (net.Listener, error) {
			return mockListener, nil
		},
	}

	// 测试场景1: 正常初始化
	t.Run("Normal initialization", func(t *testing.T) {
		// 设置模拟行为
		mockListener.EXPECT().Addr().Return(&net.TCPAddr{Port: 8080}).AnyTimes()
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

		// 执行测试
		err := s.Init()

		// 验证结果
		assert.NoError(t, err)
		assert.NotNil(t, s.cMux)
	})
}

// TestServer_MatchFor 测试协议匹配
func TestServer_MatchFor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(ctrl)
	mockListener := NewMockListener(ctrl)

	// 创建server实例
	s := &server{
		logger:   mockLogger,
		metadata: make(g.Metadata),
		listen: func(network, address string) (net.Listener, error) {
			return mockListener, nil
		},
	}

	// 初始化服务器
	mockListener.EXPECT().Addr().Return(&net.TCPAddr{Port: 8080}).AnyTimes()
	s.Init()

	// 测试场景1: 匹配GRPC协议
	t.Run("Match GRPC protocol", func(t *testing.T) {
		listener := s.MatchFor(g.GRPC)
		assert.NotNil(t, listener)
		assert.Equal(t, "true", s.metadata["grpc"])
	})

	// 测试场景2: 匹配HTTP1协议
	t.Run("Match HTTP1 protocol", func(t *testing.T) {
		listener := s.MatchFor(g.HTTP1)
		assert.NotNil(t, listener)
		assert.Equal(t, "true", s.metadata["http1"])
	})

	// 测试场景3: 匹配不支持的协议
	t.Run("Match unsupported protocol", func(t *testing.T) {
		assert.Panics(t, func() {
			s.MatchFor(g.ProtocolType(999))
		})
	})
}

// TestServer_Start 测试服务器启动
func TestServer_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(ctrl)
	mockRegistry := gMock.NewMockServiceRegistry(ctrl)
	mockTracer := gMock.NewMockTracer(ctrl)
	mockListener := NewMockListener(ctrl)
	mockConn := NewMockConn(ctrl)

	// 创建server实例
	s := &server{
		logger:           mockLogger,
		registry:         mockRegistry,
		tracer:           mockTracer,
		network:          "tcp",
		address:          "localhost:8080",
		host:             "localhost",
		port:             8080,
		serviceName:      "test-service",
		serviceUseSubNet: "0.0.0.0/0",
		metadata:         make(g.Metadata),
		listen: func(network, address string) (net.Listener, error) {
			return mockListener, nil
		},
	}

	// 初始化服务器
	mockListener.EXPECT().Addr().Return(&net.TCPAddr{Port: 8080}).AnyTimes()
	mockListener.EXPECT().Accept().Return(mockConn, nil).AnyTimes()
	mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	mockConn.EXPECT().Close().Return(nil).AnyTimes()
	s.Init()

	// 测试场景1: 正常启动服务器
	t.Run("Normal server start", func(t *testing.T) {
		// 设置模拟行为
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockRegistry.EXPECT().Register(gomock.Any()).Return(nil)
		mockTracer.EXPECT().Go(gomock.Any()).Do(func(fn func()) {
			go fn()
		})

		// 执行测试
		err := s.Start()
		time.Sleep(50 * time.Millisecond) // 等待goroutine启动

		// 验证结果
		assert.NoError(t, err)
		assert.False(t, s.stopFlag)
	})
}

// TestServer_Stop 测试服务器停止
func TestServer_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(ctrl)
	mockRegistry := gMock.NewMockServiceRegistry(ctrl)
	mockListener := NewMockListener(ctrl)

	// 创建server实例
	s := &server{
		logger:   mockLogger,
		registry: mockRegistry,
		network:  "tcp",
		address:  "localhost:8080",
		host:     "localhost",
		port:     8080,
		metadata: make(g.Metadata),
		listen: func(network, address string) (net.Listener, error) {
			return mockListener, nil
		},
	}

	// 初始化服务器
	mockListener.EXPECT().Addr().Return(&net.TCPAddr{Port: 8080}).AnyTimes()
	s.Init()

	// 测试场景1: 正常停止服务器
	t.Run("Normal server stop", func(t *testing.T) {
		// 设置模拟行为
		mockLogger.EXPECT().Warnf(gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		// 执行测试
		err := s.Stop()

		// 验证结果
		assert.NoError(t, err)
		assert.True(t, s.stopFlag)
	})
}
