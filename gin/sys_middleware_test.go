package gin

import (
	"bytes"
	"github.com/gone-io/gone/mock/v2"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_SysMiddleware_GonerName(t *testing.T) {
	m := &SysMiddleware{}
	assert.Equal(t, IdGoneGinSysMiddleware, m.GonerName())
}

func Test_SysMiddleware_Init(t *testing.T) {
	// Test with limit enabled
	m := &SysMiddleware{
		enableLimit: true,
		limit:       100,
		burst:       300,
	}

	err := m.Init()
	assert.Nil(t, err)
	assert.NotNil(t, m.limiter)

	// Test with limit disabled
	m2 := &SysMiddleware{
		enableLimit: false,
	}

	err = m2.Init()
	assert.Nil(t, err)
	assert.Nil(t, m2.limiter)
}

func Test_SysMiddleware_allow(t *testing.T) {
	// Test with limit enabled
	m := &SysMiddleware{
		enableLimit: true,
		limit:       100,
		burst:       300,
	}
	m.Init()

	assert.True(t, m.allow())

	// Test with limit disabled
	m2 := &SysMiddleware{
		enableLimit: false,
	}
	m2.Init()

	assert.True(t, m2.allow())
}

func Test_SysMiddleware_Process_Disabled(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create a test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// Create middleware with disable=true
	m := &SysMiddleware{
		disable: true,
	}

	// Process should just call Next() and return
	m.Process(c)

	// Since Next() is not actually doing anything in the test context,
	// we just verify that the function didn't panic
	assert.Equal(t, 200, w.Code) // Default status is 200
}

func Test_SysMiddleware_Process_HealthCheck(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create a test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/health", nil)

	// Create middleware with healthCheckUrl set
	m := &SysMiddleware{
		healthCheckUrl: "/health",
	}

	// Process should abort with 200 status
	m.Process(c)

	assert.Equal(t, 200, w.Code)
	assert.True(t, c.IsAborted())
}

func Test_SysMiddleware_Process_RateLimitExceeded(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockResponser := NewMockResponser(controller)

	// Create a test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// Set up expectations
	mockResponser.EXPECT().Failed(gomock.Any(), gomock.Any()).Do(func(ctx XContext, err error) {
		be, ok := err.(gone.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusTooManyRequests, be.Code())
		assert.Equal(t, TooManyRequests, be.Msg())
	})

	// Create middleware with rate limit exceeded
	m := &SysMiddleware{
		enableLimit: true,
		limit:       0, // Set limit to 0 to ensure it's exceeded
		burst:       0,
		resHandler:  mockResponser,
	}
	m.Init()

	// Process should call Failed with TooManyRequests error
	m.Process(c)
}

func Test_SysMiddleware_process(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := mock.NewMockLogger(controller)

	// Create a test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// Set up expectations for logging
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	// Create middleware
	m := &SysMiddleware{
		logger:          mockLogger,
		showRequestLog:  true,
		showResponseLog: true,
		showRequestTime: true,
	}

	// Save original testInProcess function and restore it after the test
	originalTestInProcess := testInProcess
	defer func() { testInProcess = originalTestInProcess }()

	// Set testInProcess to verify it's called
	testInProcessCalled := false
	testInProcess = func(context *gin.Context) {
		testInProcessCalled = true
	}

	// Call process
	m.process(c)

	// Verify testInProcess was called
	assert.True(t, testInProcessCalled)
}

func Test_SysMiddleware_recover(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := mock.NewMockLogger(controller)
	mockResponser := NewMockResponser(controller)

	// Create a test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// Set up expectations
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
	mockResponser.EXPECT().Failed(gomock.Any(), gomock.Any()).Times(1)

	// Create middleware
	m := &SysMiddleware{
		logger:     mockLogger,
		resHandler: mockResponser,
	}

	// Test recover with panic
	defer func() {
		// Verify that the context was aborted
		assert.True(t, c.IsAborted())
	}()
	defer m.recover(c)

	// Trigger a panic
	panic("test panic")
}

func Test_SysMiddleware_log(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create mock objects
	mockLogger := mock.NewMockLogger(controller)

	// Test cases for different log formats
	tests := []struct {
		name      string
		logFormat string
	}{
		{"console format", "console"},
		{"json format", "json"},
	}

	for _, tt := range tests {
		t := tt
		t.name = "Test_SysMiddleware_log_" + tt.name

		// Set up expectations
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)

		// Create middleware
		m := &SysMiddleware{
			logger:    mockLogger,
			logFormat: tt.logFormat,
		}

		// Call log
		m.log("test", map[string]any{
			"key1": "value1",
			"key2": 123,
		})
	}
}

func Test_cloneRequestBody(t *testing.T) {
	// Create a test context with a request body
	body := []byte(`{"test":"value"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))

	// Clone the request body
	clonedBody, err := cloneRequestBody(c)

	// Verify the cloned body matches the original
	assert.Nil(t, err)
	assert.Equal(t, body, clonedBody)

	// Verify the request body can still be read
	bodyBytes, err := c.GetRawData()
	assert.Nil(t, err)
	assert.Equal(t, body, bodyBytes)
}

func Test_CustomResponseWriter(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var response string
	// Create a test recorder
	w := NewMockResponseWriter(controller)
	w.EXPECT().Write(gomock.Any()).
		DoAndReturn(func(data []byte) (int, error) {
			response += string(data)
			return len(data), nil
		}).
		AnyTimes()
	w.EXPECT().WriteString(gomock.Any()).
		DoAndReturn(func(data string) (int, error) {
			response += data
			return len(data), nil
		}).
		AnyTimes()

	// Create a CustomResponseWriter
	crw := &CustomResponseWriter{
		ResponseWriter: w,
		body:           bytes.NewBufferString(""),
	}

	// Test Write method
	writeData := []byte("test data")
	n, err := crw.Write(writeData)
	assert.Nil(t, err)
	assert.Equal(t, len(writeData), n)
	assert.Equal(t, "test data", crw.body.String())
	assert.Equal(t, "test data", response)

	// Test WriteString method
	crw.body.Reset()
	response = ""

	writeString := "test string"
	n, err = crw.WriteString(writeString)
	assert.Nil(t, err)
	assert.Equal(t, len(writeString), n)
	assert.Equal(t, "test string", crw.body.String())
	assert.Equal(t, "test string", response)
}
