package schedule

import (
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/tracer"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 模拟获取锁失败的DoLocker
type failLocker struct {
	gone.Flag
}

func (l *failLocker) LockAndDo(key string, fn func(), lockTime, checkPeriod time.Duration) (err error) {
	return errors.New("failed to acquire lock")
}

// 测试多个scheduler同时注册的情况
func Test_schedule_MultipleSchedulers(t *testing.T) {
	_ = os.Setenv("GONE_SCHEDULE_IN-CLUSTER", "false")
	defer func() {
		_ = os.Setenv("GONE_SCHEDULE_IN-CLUSTER", "true")
	}()

	controller := gomock.NewController(t)
	defer controller.Finish()

	var mu sync.Mutex
	i := 0
	j := 0

	// 创建第一个scheduler
	scheduler1 := NewMockScheduler(controller)
	scheduler1.EXPECT().
		Cron(gomock.Any()).
		Do(func(run RunFuncOnceAt) {
			run("0/1 * * * * *", "test-MultipleSchedulers", func() {
				mu.Lock()
				i++
				mu.Unlock()
			})
		})

	// 创建第二个scheduler
	scheduler2 := NewMockScheduler(controller)
	scheduler2.EXPECT().
		Cron(gomock.Any()).
		Do(func(run RunFuncOnceAt) {
			run("0/1 * * * * *", "test2", func() {
				mu.Lock()
				j++
				mu.Unlock()
			})
		})

	gone.
		NewApp(Load).
		Load(scheduler1).
		Load(scheduler2).
		Test(func(s *schedule) {
			assert.Equal(t, 2, len(s.schedulers))
			time.Sleep(2 * time.Second)
		})

	mu.Lock()
	assert.Equal(t, 2, i)
	assert.Equal(t, 2, j)
	mu.Unlock()
}

// 测试不同cron表达式的解析和执行
func Test_schedule_DifferentCronExpressions(t *testing.T) {
	_ = os.Setenv("GONE_SCHEDULE_IN-CLUSTER", "false")
	defer func() {
		_ = os.Setenv("GONE_SCHEDULE_IN-CLUSTER", "true")
	}()

	controller := gomock.NewController(t)
	defer controller.Finish()

	var mu sync.Mutex
	counts := map[string]int{
		"every_second":    0,
		"every_2_seconds": 0,
	}

	scheduler := NewMockScheduler(controller)
	scheduler.EXPECT().
		Cron(gomock.Any()).
		Do(func(run RunFuncOnceAt) {
			// 每秒执行一次
			run("* * * * * *", "every_second", func() {
				mu.Lock()
				counts["every_second"]++
				mu.Unlock()
			})

			// 每2秒执行一次
			run("*/2 * * * * *", "every_2_seconds", func() {
				mu.Lock()
				counts["every_2_seconds"]++
				mu.Unlock()
			})
		})

	gone.
		NewApp(Load).
		Load(scheduler).
		Test(func(s *schedule) {
			time.Sleep(3 * time.Second)
		})

	mu.Lock()
	// 3秒内应该执行约3次
	assert.GreaterOrEqual(t, counts["every_second"], 3)
	// 3秒内应该执行约1-2次
	assert.GreaterOrEqual(t, counts["every_2_seconds"], 1)
	assert.LessOrEqual(t, counts["every_2_seconds"], 2)
	mu.Unlock()
}

// 测试锁定时间和检查周期配置
func Test_schedule_LockTimeAndCheckPeriod(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	scheduler := NewMockScheduler(controller)
	scheduler.EXPECT().
		Cron(gomock.Any()).
		Do(func(run RunFuncOnceAt) {
			run("0/1 * * * * *", "test-LockTimeAndCheckPeriod", func() {})
		}).
		AnyTimes()

	gone.
		NewApp(Load, func(loader gone.Loader) error {
			return loader.Load(&locker{})
		}).
		Load(scheduler).
		Test(func(s *schedule) {
			// 修改锁定时间和检查周期
			s.lockTime = 5 * time.Second
			s.checkPeriod = 1 * time.Second
			s.isCluster = true

			err := s.Start()
			assert.Nil(t, err)
			assert.Equal(t, 5*time.Second, s.lockTime)
			assert.Equal(t, 1*time.Second, s.checkPeriod)
		})
}

// 测试在集群模式下锁获取失败的情况
func Test_schedule_LockFailure(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var mu sync.Mutex
	i := 0
	scheduler := NewMockScheduler(controller)
	scheduler.EXPECT().
		Cron(gomock.Any()).
		Do(func(run RunFuncOnceAt) {
			run("0/1 * * * * *", "test-LockFailure", func() {
				mu.Lock()
				i++
				mu.Unlock()
			})
		}).
		AnyTimes()

	gone.
		NewApp(Load, func(loader gone.Loader) error {
			return loader.Load(&failLocker{})
		}).
		Load(scheduler).
		Test(func(s *schedule) {
			s.isCluster = true
			err := s.Start()
			assert.Nil(t, err)
			assert.NotNil(t, s.locker)
			time.Sleep(2 * time.Second)
		})

	// 由于锁获取失败，任务不会执行，计数应该为0
	mu.Lock()
	assert.Equal(t, 0, i)
	mu.Unlock()
}

// 测试tracer在不同情况下的行为
func Test_schedule_TracerBehavior(t *testing.T) {
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
			run("0/1 * * * * *", "test", func() {
				mu.Lock()
				i++
				mu.Unlock()
			})
		})

	// 创建一个自定义的tracer加载器，用于验证tracer的行为
	gone.
		NewApp(tracer.Load, Load).
		Load(scheduler).
		Test(func(s *schedule) {
			assert.NotNil(t, s.tracer)

			// 直接测试wrapFn函数，验证tracer的行为
			wrappedFn := s.wrapFn(func() {
				mu.Lock()
				i++
				mu.Unlock()
			}, "tracer-test")

			// 执行包装函数
			wrappedFn()

			// 验证函数已执行
			mu.Lock()
			assert.Equal(t, 1, i)
			mu.Unlock()
		})
}
