package tracer

import (
	"fmt"
	"sync"
	"testing"

	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
)

// 模拟Logger实现
type mockLogger struct {
	gone.Flag
	warningMessages []string
}

func (m *mockLogger) Warnf(format string, args ...any) {
	m.warningMessages = append(m.warningMessages, format)
}

func (m *mockLogger) Infof(format string, args ...any)  {}
func (m *mockLogger) Debugf(format string, args ...any) {}
func (m *mockLogger) Errorf(format string, args ...any) {}

// 测试GetTraceId和SetTraceId的基本功能
func TestGetAndSetTraceId(t *testing.T) {
	// 初始状态下应该没有traceId
	traceId := GetTraceId()
	assert.Equal(t, "", traceId, "初始状态下traceId应为空")

	// 设置一个指定的traceId
	expectedId := "test-trace-id"
	var actualId string
	SetTraceId(expectedId, func() {
		actualId = GetTraceId()
	})
	assert.Equal(t, expectedId, actualId, "设置的traceId应该能够被获取")

	// 设置完成后，外部应该无法获取到traceId
	traceId = GetTraceId()
	assert.Equal(t, "", traceId, "SetTraceId执行完后，外部traceId应为空")
}

// 测试自动生成traceId
func TestAutoGenerateTraceId(t *testing.T) {
	var traceId string
	SetTraceId("", func() {
		traceId = GetTraceId()
	})

	assert.NotEqual(t, "", traceId, "空traceId参数应该触发自动生成")
	assert.Len(t, traceId, 36, "自动生成的traceId应该是一个UUID格式")
}

// 测试重复设置traceId的情况
func TestSetTraceIdTwice(t *testing.T) {
	mockLog := &mockLogger{}

	SetTraceId("first-id", func() {
		// 在已有traceId的情况下再次设置
		SetTraceId("second-id", func() {
			// 不应该被设置为second-id
			assert.Equal(t, "first-id", GetTraceId(), "不应该覆盖已有的traceId")
		}, mockLog.Warnf)

		// 应该有警告日志
		assert.Len(t, mockLog.warningMessages, 1, "应该记录警告日志")
		assert.Contains(t, mockLog.warningMessages[0], "SetTraceId not success", "警告日志内容不正确")
	})
}

// 测试tracer结构体的实现
func TestTracerImplementation(t *testing.T) {
	tracerImpl := &tracer{
		logger: gone.GetDefaultLogger(),
	}

	// 测试GonerName方法
	assert.Equal(t, "gone-tracer", tracerImpl.GonerName(), "GonerName返回值不正确")

	// 测试GetTraceId方法
	SetTraceId("test-id", func() {
		assert.Equal(t, "test-id", tracerImpl.GetTraceId(), "tracer.GetTraceId返回值不正确")
	})

	// 测试SetTraceId方法
	var actualId string
	tracerImpl.SetTraceId("tracer-impl-id", func() {
		actualId = GetTraceId()
	})
	assert.Equal(t, "tracer-impl-id", actualId, "tracer.SetTraceId设置的traceId不正确")
}

// 测试Go方法在有traceId的情况
func TestGoWithTraceId(t *testing.T) {
	tracerImpl := &tracer{
		logger: gone.GetDefaultLogger(),
	}
	var wg sync.WaitGroup
	var childTraceId string

	wg.Add(1)
	SetTraceId("parent-trace-id", func() {
		parentTraceId := GetTraceId()
		assert.Equal(t, "parent-trace-id", parentTraceId, "父goroutine的traceId不正确")

		tracerImpl.Go(func() {
			defer wg.Done()
			childTraceId = GetTraceId()
		})
	})

	wg.Wait()
	assert.Equal(t, "parent-trace-id", childTraceId, "子goroutine应该继承父goroutine的traceId")
}

// 测试Go方法在没有traceId的情况
func TestGoWithoutTraceId(t *testing.T) {
	tracerImpl := &tracer{
		logger: gone.GetDefaultLogger(),
	}
	var wg sync.WaitGroup
	var childTraceId string

	wg.Add(1)
	// 不设置traceId
	assert.Equal(t, "", GetTraceId(), "初始状态下traceId应为空")

	tracerImpl.Go(func() {
		defer wg.Done()
		childTraceId = GetTraceId()
	})

	wg.Wait()
	assert.Equal(t, "", childTraceId, "没有父traceId时，子goroutine的traceId也应为空")
}

// 测试并发情况下的traceId传递
func TestConcurrentTraceId(t *testing.T) {
	tracerImpl := &tracer{logger: gone.GetDefaultLogger()}
	var wg sync.WaitGroup
	results := make(map[string]string)
	var mu sync.Mutex

	// 创建多个goroutine，每个都有自己的traceId
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			traceIdKey := fmt.Sprintf("trace-%d", index)
			SetTraceId(traceIdKey, func() {
				// 在每个goroutine中再创建一个子goroutine
				var childWg sync.WaitGroup
				childWg.Add(1)

				tracerImpl.Go(func() {
					defer childWg.Done()
					// 记录子goroutine获取到的traceId
					mu.Lock()
					results[traceIdKey] = GetTraceId()
					mu.Unlock()
				})

				childWg.Wait()
			})
			wg.Done()
		}(i)
	}

	wg.Wait()

	// 验证每个子goroutine都正确获取到了父goroutine的traceId
	assert.Len(t, results, 10, "应该有10个结果")
	for i := 0; i < 10; i++ {
		traceIdKey := fmt.Sprintf("trace-%d", i)
		assert.Equal(t, traceIdKey, results[traceIdKey], "子goroutine的traceId与父goroutine不匹配")
	}
}

// 测试Load函数
func TestLoad(t *testing.T) {
	// 创建一个模拟的Loader
	loader := &mockLoader{}

	// 调用Load函数
	err := Load(loader)

	// 验证结果
	assert.NoError(t, err, "Load函数不应返回错误")
	assert.True(t, loader.loadCalled, "应该调用了Loader.Load方法")
}

// 测试Priest函数（已废弃但仍需测试）
func TestPriest(t *testing.T) {
	// 创建一个模拟的Loader
	loader := &mockLoader{}

	// 调用Priest函数
	err := Priest(loader)

	// 验证结果
	assert.NoError(t, err, "Priest函数不应返回错误")
	assert.True(t, loader.loadCalled, "应该调用了Loader.Load方法")
}

// 模拟Loader实现
type mockLoader struct {
	loadCalled bool
}

func (m *mockLoader) Load(goner gone.Goner, options ...gone.Option) error {
	m.loadCalled = true
	return nil
}

func (m *mockLoader) Loaded(gone.LoaderKey) bool {
	return false
}
