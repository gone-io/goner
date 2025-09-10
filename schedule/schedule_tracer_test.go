package schedule

import (
	"github.com/gone-io/goner/g"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 自定义tracer实现，用于测试
type mockTracer struct {
	gone.Flag
	g.Tracer
	called bool
	once   sync.Once
}

func (m *mockTracer) SetTraceId(traceId string, fn func()) {
	m.once.Do(func() {
		m.called = true
	})
	fn()
}

// 测试tracer在没有设置的情况下的行为
func Test_schedule_NoTracer(t *testing.T) {
	_ = os.Setenv("GONE_SCHEDULE_IN-CLUSTER", "false")
	defer func() {
		_ = os.Setenv("GONE_SCHEDULE_IN-CLUSTER", "true")
	}()

	controller := gomock.NewController(t)
	defer controller.Finish()

	var mu sync.Mutex
	i := 0
	scheduler := NewMockScheduler(controller)
	scheduler.EXPECT().
		Cron(gomock.Any()).
		Do(func(run RunFuncOnceAt) {
			run("0/1 * * * * *", "test-NoTracer", func() {
				mu.Lock()
				i++
				mu.Unlock()
			})
		})

	gone.
		NewApp(Load).
		Load(scheduler).
		Test(func(s *schedule) {
			assert.Nil(t, s.tracer)
			time.Sleep(2 * time.Second)
		})

	mu.Lock()
	assert.Equal(t, 2, i)
	mu.Unlock()
}

// 测试自定义tracer的行为
func Test_schedule_CustomTracer(t *testing.T) {
	_ = os.Setenv("GONE_SCHEDULE_IN-CLUSTER", "false")
	defer func() {
		_ = os.Setenv("GONE_SCHEDULE_IN-CLUSTER", "true")
	}()

	controller := gomock.NewController(t)
	defer controller.Finish()

	var mu sync.Mutex
	i := 0
	scheduler := NewMockScheduler(controller)
	scheduler.EXPECT().
		Cron(gomock.Any()).
		Do(func(run RunFuncOnceAt) {
			run("0/1 * * * * *", "test-CustomTracer", func() {
				mu.Lock()
				i++
				mu.Unlock()
			})
		})

	// 创建自定义tracer
	customTracer := &mockTracer{}

	gone.
		NewApp(Load, func(loader gone.Loader) error {
			return loader.Load(customTracer)
		}).
		Load(scheduler).
		Test(func(s *schedule) {
			// 验证tracer已被注入
			assert.NotNil(t, s.tracer)

			// 直接测试wrapFn函数
			wrappedFn := s.wrapFn(func() {
				mu.Lock()
				i++
				mu.Unlock()
			}, "custom-tracer-test")

			// 执行包装函数
			wrappedFn()

			// 验证自定义tracer被调用
			assert.True(t, customTracer.called)

			// 验证函数已执行
			mu.Lock()
			assert.True(t, i >= 1)
			mu.Unlock()
		})
}

// 测试在集群模式下tracer和locker的交互
func Test_schedule_TracerWithLocker(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var mu sync.Mutex
	i := 0
	scheduler := NewMockScheduler(controller)
	scheduler.EXPECT().
		Cron(gomock.Any()).
		Do(func(run RunFuncOnceAt) {
			run("0/1 * * * * *", "test-TracerWithLocker", func() {
				mu.Lock()
				i++
				mu.Unlock()
			})
		}).
		AnyTimes()

	// 创建自定义tracer
	customTracer := &mockTracer{}

	gone.
		NewApp(Load, func(loader gone.Loader) error {
			_ = loader.Load(customTracer)
			return loader.Load(&locker{})
		}).
		Load(scheduler).
		Test(func(s *schedule) {

			// 验证tracer和locker都已被注入
			assert.NotNil(t, s.tracer)
			assert.NotNil(t, s.locker)

			time.Sleep(2 * time.Second)

			// 验证函数已执行
			mu.Lock()
			assert.Equal(t, 2, i)
			mu.Unlock()

			// 验证自定义tracer被调用
			assert.True(t, customTracer.called)
		})
}
