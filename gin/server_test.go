package gin

import (
	"github.com/gone-io/gone/mock/v2"
	gMock "github.com/gone-io/goner/g/mock"
	"net/http"
	"testing"
	"time"

	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_server_GonerName(t *testing.T) {
	s := &server{}
	assert.Equal(t, IdGoneGin, s.GonerName())
}

func Test_NewGinServer(t *testing.T) {
	goner, option := NewGinServer()
	assert.NotNil(t, goner)
	assert.NotNil(t, option)
	assert.IsType(t, &server{}, goner)
}

func Test_server_initListener_WithCMux(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockCMux := gMock.NewMockCmux(controller)
	mockListener := NewMockListener(controller)

	// Set up expectations
	mockCMux.EXPECT().MatchFor(g.HTTP1).Return(mockListener)
	mockCMux.EXPECT().GetAddress().Return("127.0.0.1:8080")

	// Create server instance
	s := &server{
		cMuxServer: mockCMux,
	}

	// Test initListener method
	err := s.initListener()
	assert.Nil(t, err)
	assert.Equal(t, mockListener, s.listener)
	assert.Equal(t, "127.0.0.1:8080", s.getAddress())
}

func Test_server_initListener_WithoutCMux(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	addr := NewMockAddr(controller)
	addr.EXPECT().String().Return("127.0.0.1:8080")
	mockListener := NewMockListener(controller)
	mockListener.EXPECT().Addr().Return(addr)

	// Create server instance with custom createListener function
	s := &server{
		host: "127.0.0.1",
		port: 8080,
		createListener: func(s *server) error {
			s.listener = mockListener
			return nil
		},
	}

	// Test initListener method
	err := s.initListener()
	assert.Nil(t, err)
	assert.Equal(t, mockListener, s.listener)
	assert.Equal(t, "127.0.0.1:8080", s.getAddress())
}

func Test_server_mount(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := mock.NewMockLogger(controller)
	mockController := NewMockController(controller)

	// Set up expectations
	mockLogger.EXPECT().Warnf(gomock.Any()).AnyTimes()
	mockController.EXPECT().Mount().Return(nil)

	// Create server instance
	s := &server{
		logger:      mockLogger,
		controllers: []Controller{mockController},
	}

	// Test mount method
	err := s.mount()
	assert.Nil(t, err)
}

func Test_server_mount_Error(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := mock.NewMockLogger(controller)
	mockController := NewMockController(controller)

	// Set up expectations
	mockLogger.EXPECT().Warnf(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	mockController.EXPECT().Mount().Return(assert.AnError)

	// Create server instance
	s := &server{
		logger:      mockLogger,
		controllers: []Controller{mockController},
	}

	// Test mount method with error
	err := s.mount()
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func Test_server_Start(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := mock.NewMockLogger(controller)
	mockHandler := NewMockHandler(controller)

	// Set up expectations
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warnf(gomock.Any()).AnyTimes()

	// Create server instance
	s := &server{
		logger:         mockLogger,
		httpHandler:    mockHandler,
		controllers:    []Controller{},
		createListener: createListener,
		host:           "127.0.0.1",
		port:           0,
	}

	// Test Start method
	err := s.Start()
	assert.Nil(t, err)
	assert.False(t, s.stopFlag)
	assert.NotNil(t, s.httpServer)
	assert.Equal(t, mockHandler, s.httpServer.Handler)
	time.Sleep(time.Millisecond * 100)
}

func Test_server_Stop(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := mock.NewMockLogger(controller)

	// Set up expectations
	mockLogger.EXPECT().Warnf(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

	// Create a test HTTP server
	testServer := &http.Server{}

	// Create server instance
	s := &server{
		logger:            mockLogger,
		httpServer:        testServer,
		maxWaitBeforeStop: 100 * time.Millisecond,
	}

	// Test Stop method
	err := s.Stop()
	assert.Nil(t, err)
	assert.True(t, s.stopFlag)
}

func Test_server_Stop_NilServer(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := mock.NewMockLogger(controller)

	// Set up expectations
	mockLogger.EXPECT().Warnf(gomock.Any()).AnyTimes()

	// Create server instance with nil httpServer
	s := &server{
		logger:     mockLogger,
		httpServer: nil,
	}

	// Test Stop method with nil server
	err := s.Stop()
	assert.Nil(t, err)
}

func Test_server_processServeError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := mock.NewMockLogger(controller)

	// Set up expectations for normal error
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

	// Create server instance
	s := &server{
		logger: mockLogger,
	}

	// Test processServeError method with stopFlag=false
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		assert.Equal(t, assert.AnError, r)
	}()

	s.processServeError(assert.AnError)
}

func Test_server_processServeError_Stopping(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := mock.NewMockLogger(controller)

	// Set up expectations for warning
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

	// Create server instance with stopFlag=true
	s := &server{
		logger:   mockLogger,
		stopFlag: true,
	}

	// Test processServeError method with stopFlag=true
	s.processServeError(assert.AnError)
	// No panic should occur
}

func Test_createListener(t *testing.T) {

	// Create a server instance
	s := &server{
		host: "127.0.0.1", // Use port 0 to get a random available port
		port: 0,
	}

	// Test createListener function
	err := createListener(s)
	assert.Nil(t, err)
	assert.NotNil(t, s.listener)

	// Clean up
	s.listener.Close()
}
