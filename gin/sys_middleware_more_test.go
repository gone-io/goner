package gin

import (
	"bytes"
	mock "github.com/gone-io/gone/v2"
	gMock "github.com/gone-io/goner/g/mock"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 测试 SysMiddleware.Process 方法 - 处理 traceId
func Test_SysMiddleware_Process_TraceId(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟对象
	mockTracer := gMock.NewMockTracer(controller)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// 测试场景1: 请求头中有 traceId
	testTraceId := "test-trace-id"
	c.Request.Header.Set("X-Trace-Id", testTraceId)

	// 设置期望
	mockTracer.EXPECT().SetTraceId(testTraceId, gomock.Any()).Do(func(_ string, callback func()) {
		callback()
	})

	// 创建中间件
	m := &SysMiddleware{
		tracer:      mockTracer,
		tracerIdKey: "X-Trace-Id",
	}

	// 调用 Process
	m.Process(c)

	// 测试场景2: 请求头中没有 traceId
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest("GET", "/test", nil)

	// 设置期望 - 应该生成新的 traceId
	mockTracer.EXPECT().SetTraceId(gomock.Any(), gomock.Any()).Do(func(traceId string, callback func()) {
		assert.NotEmpty(t, traceId)
		callback()
	})

	// 调用 Process
	m.Process(c2)

	// 测试场景3: 没有 tracer
	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)
	c3.Request = httptest.NewRequest("GET", "/test", nil)

	// 创建没有 tracer 的中间件
	m2 := &SysMiddleware{
		tracerIdKey: "X-Trace-Id",
	}

	// 调用 Process
	m2.Process(c3)

	// 验证请求上下文中有 traceId
	traceId := c3.Request.Context().Value("X-Trace-Id")
	assert.NotNil(t, traceId)
}

// 测试 SysMiddleware.requestLog 方法 - 处理不同的请求日志场景
func Test_SysMiddleware_requestLog_Scenarios(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(controller)

	// 测试场景1: 完整的请求日志
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"test":"data"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("User-Agent", "test-agent")
	c.Request.Header.Set("Referer", "http://test-referer")
	c.Request.Header.Set("X-Request-Id", "test-request-id")

	// 设置期望
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)

	// 创建中间件
	m := &SysMiddleware{
		logger:                     mockLogger,
		showRequestLog:             true,
		logRequestId:               true,
		logRemoteIp:                true,
		logUserAgent:               true,
		logReferer:                 true,
		logRequestBody:             true,
		requestBodyLogContentTypes: "application/json",
		requestIdKey:               "X-Request-Id",
	}

	// 调用 requestLog
	m.requestLog(c)

	// 测试场景2: 最小化请求日志
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest("GET", "/test", nil)

	// 设置期望
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)

	// 创建中间件
	m2 := &SysMiddleware{
		logger:         mockLogger,
		showRequestLog: true,
		logRequestId:   false,
		logRemoteIp:    false,
		logUserAgent:   false,
		logReferer:     false,
		logRequestBody: false,
	}

	// 调用 requestLog
	m2.requestLog(c2)

	// 测试场景3: 请求体日志长度限制
	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)
	c3.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"test":"data with very long content"}`))
	c3.Request.Header.Set("Content-Type", "application/json")

	// 设置期望
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)

	// 创建中间件
	m3 := &SysMiddleware{
		logger:                     mockLogger,
		showRequestLog:             true,
		logRequestBody:             true,
		requestBodyLogContentTypes: "application/json",
		logDataMaxLength:           10,
	}

	// 调用 requestLog
	m3.requestLog(c3)

	// 测试场景4: 代理后的 IP
	w4 := httptest.NewRecorder()
	c4, _ := gin.CreateTestContext(w4)
	c4.Request = httptest.NewRequest("GET", "/test", nil)
	c4.Request.Header.Set("X-Forwarded-For", "192.168.1.1")

	// 设置期望
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)

	// 创建中间件
	m4 := &SysMiddleware{
		logger:         mockLogger,
		showRequestLog: true,
		logRemoteIp:    true,
		isAfterProxy:   true,
	}

	// 调用 requestLog
	m4.requestLog(c4)

	// 测试场景5: 请求体克隆错误
	w5 := httptest.NewRecorder()
	c5, _ := gin.CreateTestContext(w5)
	// 创建一个已经被读取过的请求，模拟 GetRawData 错误
	c5.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"test":"data"}`))
	c5.Request.Header.Set("Content-Type", "application/json")
	_, _ = c5.GetRawData() // 先读取一次，使下次读取失败
	c5.Request.Body = nil

	// 设置期望
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)

	// 创建中间件
	m5 := &SysMiddleware{
		logger:                     mockLogger,
		showRequestLog:             true,
		logRequestBody:             true,
		requestBodyLogContentTypes: "application/json",
	}

	// 调用 requestLog
	m5.requestLog(c5)
}

// 测试 SysMiddleware.responseLog 方法 - 处理不同的响应日志场景
func Test_SysMiddleware_responseLog_Scenarios(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(controller)

	// 测试场景1: 完整的响应日志
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Writer.Header().Set("Content-Type", "application/json")

	// 设置期望
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)

	// 创建中间件
	m := &SysMiddleware{
		logger:                      mockLogger,
		showResponseLog:             true,
		responseBodyLogContentTypes: "application/json",
	}

	// 模拟 next 函数，写入响应内容
	nextFunc := func() {
		c.JSON(200, map[string]interface{}{"result": "success"})
	}

	// 调用 responseLog
	m.responseLog(c, nextFunc)

	// 测试场景2: 响应体日志长度限制
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest("GET", "/test", nil)
	c2.Writer.Header().Set("Content-Type", "application/json")

	// 设置期望
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)

	// 创建中间件
	m2 := &SysMiddleware{
		logger:                      mockLogger,
		showResponseLog:             true,
		responseBodyLogContentTypes: "application/json",
		logDataMaxLength:            10,
	}

	// 模拟 next 函数，写入长响应内容
	nextFunc2 := func() {
		c2.JSON(200, map[string]interface{}{"result": "success with very long content that should be truncated"})
	}

	// 调用 responseLog
	m2.responseLog(c2, nextFunc2)

	// 测试场景3: 不记录响应日志
	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)
	c3.Request = httptest.NewRequest("GET", "/test", nil)

	// 创建中间件
	m3 := &SysMiddleware{
		logger:          mockLogger,
		showResponseLog: false,
	}

	// 模拟 next 函数
	nextCalled := false
	nextFunc3 := func() {
		nextCalled = true
	}

	// 调用 responseLog
	m3.responseLog(c3, nextFunc3)

	// 验证 next 被调用
	assert.True(t, nextCalled)
}

// 测试 SysMiddleware.stat 方法
func Test_SysMiddleware_stat(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(controller)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// 测试场景1: 显示请求时间
	// 设置期望
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)

	// 创建中间件
	m := &SysMiddleware{
		logger:          mockLogger,
		showRequestTime: true,
	}

	// 调用 stat
	m.stat(c, time.Now().Add(-100*time.Millisecond))

	// 测试场景2: 不显示请求时间
	// 创建中间件
	m2 := &SysMiddleware{
		logger:          mockLogger,
		showRequestTime: false,
	}

	// 调用 stat
	m2.stat(c, time.Now())
}
