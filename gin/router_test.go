package gin

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_router_GonerName(t *testing.T) {
	r := &router{}
	assert.Equal(t, IdGoneGinRouter, r.GonerName())
}

func Test_router_Init(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := NewMockLogger(controller)
	mockProxy := NewMockHandleProxyToGin(controller)

	// Set up expectations
	mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

	// Create router instance
	r := &router{
		log:              mockLogger,
		HandleProxyToGin: mockProxy,
		mode:             "test",
	}

	// Test Init method
	err := r.Init()
	assert.Nil(t, err)
	assert.NotNil(t, r.Engine)
}

func Test_router_GetGinRouter(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := NewMockLogger(controller)
	mockProxy := NewMockHandleProxyToGin(controller)

	// Create router instance and initialize it
	r := &router{
		log:              mockLogger,
		HandleProxyToGin: mockProxy,
		mode:             "test",
	}
	r.Init()

	// Test GetGinRouter method
	router := r.GetGinRouter()
	assert.NotNil(t, router)
	assert.Equal(t, r.Engine, router)
}

func Test_router_getR(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := NewMockLogger(controller)
	mockProxy := NewMockHandleProxyToGin(controller)

	// Create router instance and initialize it
	r := &router{
		log:              mockLogger,
		HandleProxyToGin: mockProxy,
		mode:             "test",
	}
	r.Init()

	// Test getR method when r is nil
	router := r.getR()
	assert.NotNil(t, router)
	assert.Equal(t, r.Engine, router)

	// Test getR method when r is not nil
	r.r = r.Engine.Group("/test")
	router = r.getR()
	assert.NotNil(t, router)
	assert.Equal(t, r.r, router)
}

func Test_router_Use(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := NewMockLogger(controller)
	mockProxy := NewMockHandleProxyToGin(controller)

	// Set up expectations
	middlewareHandlers := []gin.HandlerFunc{func(c *gin.Context) {}}
	mockProxy.EXPECT().ProxyForMiddleware(gomock.Any()).Return(middlewareHandlers).AnyTimes()

	// Create router instance and initialize it
	r := &router{
		log:              mockLogger,
		HandleProxyToGin: mockProxy,
		mode:             "test",
	}
	r.Init()

	// Test Use method
	result := r.Use(func(c *gin.Context) {})
	assert.Equal(t, r, result)
}

func Test_router_Group(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := NewMockLogger(controller)
	mockProxy := NewMockHandleProxyToGin(controller)

	// Set up expectations
	middlewareHandlers := []gin.HandlerFunc{func(c *gin.Context) {}}
	mockProxy.EXPECT().ProxyForMiddleware(gomock.Any()).Return(middlewareHandlers).AnyTimes()

	// Create router instance and initialize it
	r := &router{
		log:              mockLogger,
		HandleProxyToGin: mockProxy,
		mode:             "test",
	}
	r.Init()

	// Test Group method
	group := r.Group("/test", func(c *gin.Context) {})
	assert.NotNil(t, group)
	assert.IsType(t, &router{}, group)
}

func Test_router_Handle(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := NewMockLogger(controller)
	mockProxy := NewMockHandleProxyToGin(controller)

	// Set up expectations
	handlers := []gin.HandlerFunc{func(c *gin.Context) {}}
	mockProxy.EXPECT().Proxy(gomock.Any()).Return(handlers).AnyTimes()

	// Create router instance and initialize it
	r := &router{
		log:              mockLogger,
		HandleProxyToGin: mockProxy,
		mode:             "test",
	}
	r.Init()

	// Test Handle method
	result := r.Handle(http.MethodGet, "/test", func(c *gin.Context) {})
	assert.Equal(t, r, result)
}

func Test_router_HTTP_Methods(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := NewMockLogger(controller)
	mockProxy := NewMockHandleProxyToGin(controller)

	// Set up expectations
	handlers := []gin.HandlerFunc{func(c *gin.Context) {}}
	mockProxy.EXPECT().Proxy(gomock.Any()).Return(handlers).AnyTimes()

	// Create router instance and initialize it
	r := &router{
		log:              mockLogger,
		HandleProxyToGin: mockProxy,
		mode:             "test",
	}
	r.Init()

	// Test HTTP method functions
	tests := []struct {
		name   string
		method func(string, ...HandlerFunc) IRoutes
	}{
		{"GET", r.GET},
		{"POST", r.POST},
		{"DELETE", r.DELETE},
		{"PATCH", r.PATCH},
		{"PUT", r.PUT},
		{"OPTIONS", r.OPTIONS},
		{"HEAD", r.HEAD},
	}

	for _, tt := range tests {
		t := tt
		t.name = "Test_router_" + tt.name
		t.method("/test", func(c *gin.Context) {})
	}
}

func Test_router_Any(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := NewMockLogger(controller)
	mockProxy := NewMockHandleProxyToGin(controller)

	// Set up expectations
	handlers := []gin.HandlerFunc{func(c *gin.Context) {}}
	mockProxy.EXPECT().Proxy(gomock.Any()).Return(handlers).AnyTimes()

	// Create router instance and initialize it
	r := &router{
		log:              mockLogger,
		HandleProxyToGin: mockProxy,
		mode:             "test",
	}
	r.Init()

	// Test Any method
	result := r.Any("/test", func(c *gin.Context) {})
	assert.Equal(t, r, result)
}

func Test_logWriter_Write(t *testing.T) {
	// Create a logWriter with a custom write function
	called := false
	expectedBytes := []byte("test log message")
	writer := logWriter{
		write: func(p []byte) (n int, err error) {
			called = true
			assert.Equal(t, expectedBytes, p)
			return len(p), nil
		},
	}

	// Test Write method
	n, err := writer.Write(expectedBytes)
	assert.True(t, called)
	assert.Nil(t, err)
	assert.Equal(t, len(expectedBytes), n)
}

func Test_router_getMiddlewaresFunc(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock middleware
	mockMiddleware := NewMockMiddleware(controller)
	mockMiddleware.EXPECT().Process(gomock.Any()).AnyTimes()

	// Create router instance with middleware
	r := &router{
		middlewares: []Middleware{mockMiddleware},
	}

	// Test getMiddlewaresFunc method
	funcs := r.getMiddlewaresFunc()
	assert.Equal(t, 1, len(funcs))
}
