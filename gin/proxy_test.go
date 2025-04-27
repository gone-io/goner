package gin

import (
	"errors"
	"github.com/gone-io/gone/mock/v2"
	"github.com/gone-io/gone/v2"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_proxy_GonerName(t *testing.T) {
	p := &proxy{}
	assert.Equal(t, IdGoneGinProxy, p.GonerName())
}

func Test_proxy_Proxy(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(controller)
	mockFuncInjector := mock.NewMockFuncInjector(controller)
	mockResponser := NewMockResponser(controller)
	mockInjector := NewMockHttInjector(controller)

	// 创建 proxy 实例
	p := &proxy{
		log:          mockLogger,
		funcInjector: mockFuncInjector,
		responser:    mockResponser,
		injector:     mockInjector,
		stat:         false,
	}

	// 测试 Proxy 方法
	handlers := []HandlerFunc{
		func(ctx *gin.Context) {},
		func(ctx *gin.Context) {},
	}

	result := p.Proxy(handlers...)
	assert.Equal(t, 2, len(result))
}

func Test_proxy_ProxyForMiddleware(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(controller)
	mockFuncInjector := mock.NewMockFuncInjector(controller)
	mockResponser := NewMockResponser(controller)
	mockInjector := NewMockHttInjector(controller)

	// 创建 proxy 实例
	p := &proxy{
		log:          mockLogger,
		funcInjector: mockFuncInjector,
		responser:    mockResponser,
		injector:     mockInjector,
		stat:         false,
	}

	// 测试 ProxyForMiddleware 方法
	handlers := []HandlerFunc{
		func(ctx *gin.Context) {},
		func(ctx *gin.Context) {},
	}

	result := p.ProxyForMiddleware(handlers...)
	assert.Equal(t, 2, len(result))
}

func Test_proxy_proxyOne_GoneContext(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(controller)
	mockFuncInjector := mock.NewMockFuncInjector(controller)
	mockResponser := NewMockResponser(controller)
	mockInjector := NewMockHttInjector(controller)

	// 设置 mockResponser 的期望行为
	mockResponser.EXPECT().ProcessResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	// 创建 proxy 实例
	p := &proxy{
		log:          mockLogger,
		funcInjector: mockFuncInjector,
		responser:    mockResponser,
		injector:     mockInjector,
		stat:         false,
	}

	// 创建测试用的 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试 func(*Context) (any, error) 类型
	handler1 := func(ctx *Context) (any, error) {
		return "test", nil
	}
	ginHandler1 := p.proxyOne(handler1, true)
	ginHandler1(c)

	// 测试 func(*Context) error 类型
	handler2 := func(ctx *Context) error {
		return nil
	}
	ginHandler2 := p.proxyOne(handler2, true)
	ginHandler2(c)

	// 测试 func(*Context) 类型
	handler3 := func(ctx *Context) {
		// 空函数
	}
	ginHandler3 := p.proxyOne(handler3, true)
	ginHandler3(c)
}

func Test_proxy_proxyOne_GinContext(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(controller)
	mockFuncInjector := mock.NewMockFuncInjector(controller)
	mockResponser := NewMockResponser(controller)
	mockInjector := NewMockHttInjector(controller)

	// 设置 mockResponser 的期望行为
	mockResponser.EXPECT().ProcessResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	// 创建 proxy 实例
	p := &proxy{
		log:          mockLogger,
		funcInjector: mockFuncInjector,
		responser:    mockResponser,
		injector:     mockInjector,
		stat:         false,
	}

	// 创建测试用的 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试 func(ctx *gin.Context) 类型
	handler1 := func(ctx *gin.Context) {
		// 空函数
	}
	ginHandler1 := p.proxyOne(handler1, true)
	ginHandler1(c)

	// 测试 func(ctx *gin.Context) (any, error) 类型
	handler2 := func(ctx *gin.Context) (any, error) {
		return "test", nil
	}
	ginHandler2 := p.proxyOne(handler2, true)
	ginHandler2(c)

	// 测试 func(ctx *gin.Context) error 类型
	handler3 := func(ctx *gin.Context) error {
		return nil
	}
	ginHandler3 := p.proxyOne(handler3, true)
	ginHandler3(c)
}

func Test_proxy_proxyOne_NoContext(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(controller)
	mockFuncInjector := mock.NewMockFuncInjector(controller)
	mockResponser := NewMockResponser(controller)
	mockInjector := NewMockHttInjector(controller)

	// 设置 mockResponser 的期望行为
	mockResponser.EXPECT().ProcessResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	// 创建 proxy 实例
	p := &proxy{
		log:          mockLogger,
		funcInjector: mockFuncInjector,
		responser:    mockResponser,
		injector:     mockInjector,
		stat:         false,
	}

	// 创建测试用的 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试 func() 类型
	handler1 := func() {
		// 空函数
	}
	ginHandler1 := p.proxyOne(handler1, true)
	ginHandler1(c)

	// 测试 func() (any, error) 类型
	handler2 := func() (any, error) {
		return "test", nil
	}
	ginHandler2 := p.proxyOne(handler2, true)
	ginHandler2(c)

	// 测试 func() error 类型
	handler3 := func() error {
		return nil
	}
	ginHandler3 := p.proxyOne(handler3, true)
	ginHandler3(c)
}

func Test_proxy_buildProxyFn(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(controller)
	mockFuncInjector := mock.NewMockFuncInjector(controller)
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	mockResponser := NewMockResponser(controller)
	mockInjector := NewMockHttInjector(controller)

	// 设置 mockInjector 的期望行为
	mockInjector.EXPECT().StartBindFuncs().AnyTimes()
	mockInjector.EXPECT().BindFuncs().Return(
		func(ctx *gin.Context, arg reflect.Value) (reflect.Value, error) {
			return arg, nil
		},
	).AnyTimes()

	// 测试自定义结构体参数的函数
	type TestStruct struct {
		Name string
	}

	// 设置 mockFuncInjector 的期望行为
	mockFuncInjector.EXPECT().InjectFuncParameters(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).DoAndReturn(func(fn any, onCtx any, onStruct any) ([]reflect.Value, error) {
		testStruct := TestStruct{}

		// 简单返回一个空的参数列表
		return []reflect.Value{reflect.ValueOf(testStruct)}, nil
	}).AnyTimes()

	// 设置 mockResponser 的期望行为
	mockResponser.EXPECT().ProcessResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockResponser.EXPECT().Failed(gomock.Any(), gomock.Any()).AnyTimes()

	// 创建 proxy 实例
	p := &proxy{
		log:          mockLogger,
		funcInjector: mockFuncInjector,
		responser:    mockResponser,
		injector:     mockInjector,
		stat:         true, // 启用统计功能
	}

	// 创建测试用的 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := func(param TestStruct) error {
		return nil
	}

	ginHandler := p.buildProxyFn(handler, "testFunc", true)
	ginHandler(c)
}

func Test_proxy_buildProxyFn_Error(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(controller)
	mockFuncInjector := mock.NewMockFuncInjector(controller)
	mockResponser := NewMockResponser(controller)
	mockInjector := NewMockHttInjector(controller)

	// 设置 mockInjector 的期望行为
	mockInjector.EXPECT().StartBindFuncs().AnyTimes()

	// 设置 mockFuncInjector 的期望行为，返回错误
	mockFuncInjector.EXPECT().InjectFuncParameters(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(nil, errors.New("injection error")).AnyTimes()

	// 创建 proxy 实例
	p := &proxy{
		log:          mockLogger,
		funcInjector: mockFuncInjector,
		responser:    mockResponser,
		injector:     mockInjector,
		stat:         false,
	}

	// 测试自定义结构体参数的函数
	type TestStruct struct {
		Name string
	}

	handler := func(param TestStruct) error {
		return nil
	}

	// 应该会 panic
	assert.Panics(t, func() {
		p.buildProxyFn(handler, "testFunc", true)
	})
}

func Test_proxy_buildProxyFn_BindError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟对象
	mockLogger := mock.NewMockLogger(controller)
	mockFuncInjector := mock.NewMockFuncInjector(controller)
	mockResponser := NewMockResponser(controller)
	mockInjector := NewMockHttInjector(controller)

	// 设置 mockInjector 的期望行为
	mockInjector.EXPECT().StartBindFuncs().AnyTimes()
	mockInjector.EXPECT().BindFuncs().Return(
		func(ctx *gin.Context, arg reflect.Value) (reflect.Value, error) {
			return arg, errors.New("bind error")
		},
	).AnyTimes()

	// 设置 mockFuncInjector 的期望行为
	mockFuncInjector.EXPECT().InjectFuncParameters(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).DoAndReturn(func(fn any, onCtx gone.FuncInjectHook, onStruct gone.FuncInjectHook) ([]reflect.Value, error) {
		onStruct(reflect.TypeOf(""), 0, false)
		return []reflect.Value{reflect.ValueOf("test")}, nil
	}).AnyTimes()

	// 设置 mockResponser 的期望行为
	mockResponser.EXPECT().Failed(gomock.Any(), gomock.Any()).Times(1)

	// 创建 proxy 实例
	p := &proxy{
		log:          mockLogger,
		funcInjector: mockFuncInjector,
		responser:    mockResponser,
		injector:     mockInjector,
		stat:         false,
	}

	// 创建测试用的 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试自定义结构体参数的函数
	type TestStruct struct {
		Name string
	}

	handler := func(param TestStruct) error {
		return nil
	}

	ginHandler := p.buildProxyFn(handler, "testFunc", true)
	ginHandler(c)
}

func Test_TimeStat(t *testing.T) {
	// 测试 TimeStat 函数
	start := time.Now()
	time.Sleep(10 * time.Millisecond) // 等待一小段时间

	// 使用自定义的日志函数
	logCalled := false
	logFunc := func(format string, args ...any) {
		logCalled = true
		assert.Contains(t, format, "%s executed %v times")
	}

	TimeStat("test_function", start, logFunc)
	assert.True(t, logCalled)

	// 测试不提供日志函数的情况
	TimeStat("test_function2", start)
}
