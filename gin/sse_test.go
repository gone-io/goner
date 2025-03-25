package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewSSE(t *testing.T) {
	// Create a test recorder
	hw := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(hw)
	w := c.Writer

	// Create a new SSE instance
	sse := NewSSE(w)

	// Verify the SSE instance is created correctly
	assert.NotNil(t, sse)
	assert.IsType(t, &Sse{}, sse)
	assert.Equal(t, w, sse.(*Sse).Writer)
}

func TestSse_Start(t *testing.T) {
	// Create a test recorder
	hw := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(hw)
	w := c.Writer

	// Create a new SSE instance
	sse := NewSSE(w)

	// Call Start method
	sse.Start()

	// Verify headers are set correctly
	headers := w.Header()
	assert.Equal(t, "text/event-stream; charset=utf-8", headers.Get("Content-Type"))
	assert.Equal(t, "no-cache", headers.Get("Cache-Control"))
	assert.Equal(t, "keep-alive", headers.Get("Connection"))
	assert.Equal(t, "no", headers.Get("X-Accel-Buffering"))
}

func TestSse_Write(t *testing.T) {
	// Create a test recorder
	hw := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(hw)
	w := c.Writer

	// Create a new SSE instance
	sse := NewSSE(w)

	// Call Start method
	sse.Start()

	// Test writing different types of data
	tests := []struct {
		name     string
		data     any
		expected string
	}{
		{"string", "test message", "event: data\ndata: \"test message\"\n\n"},
		{"number", 123, "event: data\ndata: 123\n\n"},
		{"boolean", true, "event: data\ndata: true\n\n"},
		{"object", map[string]string{"key": "value"}, "event: data\ndata: {\"key\":\"value\"}\n\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the recorder body
			hw.Body.Reset()

			// Write the data
			err := sse.Write(tt.data)
			assert.Nil(t, err)

			// Verify the written data
			assert.Equal(t, tt.expected, hw.Body.String())
		})
	}
}

func TestSse_End(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Create a mock ResponseWriter
	mockWriter := NewMockResponseWriter(controller)

	// Set up expectations
	mockWriter.EXPECT().Header().Return(httptest.NewRecorder().Header()).AnyTimes()
	mockWriter.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
	mockWriter.EXPECT().WriteString(gomock.Any()).DoAndReturn(func(data string) (int, error) {
		return len(data), nil
	}).AnyTimes()
	mockWriter.EXPECT().Flush().Times(1)
	mockWriter.EXPECT().CloseNotify().Times(1)

	// Create a new SSE instance with the mock writer
	sse := &Sse{Writer: mockWriter}

	// Call End method
	err := sse.End()

	// Verify the result
	assert.Nil(t, err)
}

func TestSse_Write_Error(t *testing.T) {
	// Create a test case with invalid JSON data
	// This is a bit tricky to test since most Go values can be marshaled to JSON
	// One approach is to use a custom type that can't be marshaled

	// Create a test recorder
	hw := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(hw)
	w := c.Writer

	// Create a new SSE instance
	sse := NewSSE(w)

	// Call Start method
	sse.Start()

	// Create a mock ResponseWriter that returns an error on Write
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockWriter := NewMockResponseWriter(controller)
	mockWriter.EXPECT().Header().Return(w.Header()).AnyTimes()
	mockWriter.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
	mockWriter.EXPECT().WriteString(gomock.Any()).Return(0, assert.AnError).AnyTimes()
	mockWriter.EXPECT().Flush().AnyTimes()

	// Replace the writer with our mock
	sse = &Sse{Writer: mockWriter}

	// Write should return an error
	err := sse.Write("test")
	assert.Error(t, err)
}

func TestSse_End_Error(t *testing.T) {
	// Create a mock ResponseWriter that returns an error on Write
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockWriter := NewMockResponseWriter(controller)
	mockWriter.EXPECT().Header().Return(httptest.NewRecorder().Header()).AnyTimes()
	mockWriter.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
	mockWriter.EXPECT().WriteString(gomock.Any()).Return(0, assert.AnError).AnyTimes()

	// Create a new SSE instance with the mock writer
	sse := &Sse{Writer: mockWriter}

	// End should return an error
	err := sse.End()
	assert.Error(t, err)
}

func TestSse_Integration(t *testing.T) {
	// Create a test recorder
	controller := gomock.NewController(t)
	defer controller.Finish()

	var output string

	w := NewMockResponseWriter(controller)
	w.EXPECT().Header().Return(http.Header{}).AnyTimes()
	w.EXPECT().Flush().AnyTimes()
	w.EXPECT().CloseNotify().AnyTimes()
	w.EXPECT().WriteString(gomock.Any()).DoAndReturn(func(str string) (int, error) {
		output += str
		return len(str), nil
	}).AnyTimes()

	// Create a new SSE instance
	sse := NewSSE(w)

	// Test a complete SSE session
	sse.Start()
	sse.Write("message 1")
	sse.Write("message 2")
	sse.End()

	// Verify the output
	assert.True(t, strings.Contains(output, "event: data\ndata: \"message 1\"\n\n"))
	assert.True(t, strings.Contains(output, "event: data\ndata: \"message 2\"\n\n"))
	assert.True(t, strings.Contains(output, "event: done\ndata: [DONE]\n\n"))
}
