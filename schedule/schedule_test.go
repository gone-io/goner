package schedule

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/redis"
	"github.com/gone-io/goner/tracer"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type locker struct {
	gone.Flag
	redis.Locker
}

func (l *locker) LockAndDo(key string, fn func(), lockTime, checkPeriod time.Duration) (err error) {
	fn()
	return nil
}

func Test_schedule_Start_SingleInstance(t *testing.T) {
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
				println("test")
				mu.Lock()
				i++
				mu.Unlock()
			})
		})

	gone.
		NewApp(Load).
		Load(scheduler).
		Test(func(s *schedule) {
			assert.False(t, s.isCluster)
			time.Sleep(2 * time.Second)
		})

	mu.Lock()
	assert.Equal(t, 2, i)
	mu.Unlock()
}

func Test_schedule_Start_Cluster(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	var mu sync.Mutex
	i := 0
	scheduler := NewMockScheduler(controller)
	scheduler.EXPECT().
		Cron(gomock.Any()).
		Do(func(run RunFuncOnceAt) {
			run("0/1 * * * * *", "test", func() {
				println("test")
				mu.Lock()
				i++
				mu.Unlock()
			})
		})

	gone.
		NewApp(tracer.Load, Load, func(loader gone.Loader) error {
			return loader.Load(&locker{})
		}).
		Load(scheduler).
		Test(func(s *schedule) {
			assert.True(t, s.isCluster)
			assert.NotNil(t, s.locker)
			time.Sleep(2 * time.Second)
		})

	mu.Lock()
	assert.Equal(t, 2, i)
	mu.Unlock()
}

func Test_schedule_Start_NoSchedulers(t *testing.T) {
	gone.
		NewApp(Load).
		Test(func(s *schedule) {
			assert.NotNil(t, s)
		})
}

func Test_schedule_Start_ClusterNoLocker(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	scheduler := NewMockScheduler(controller)

	gone.
		NewApp().
		Load(scheduler).
		Test(func(s *schedule) {
			s.isCluster = true
			err := s.Start()
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "must load a `DoLocker`")
		})
}

func Test_schedule_Stop(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	scheduler := NewMockScheduler(controller)
	scheduler.EXPECT().
		Cron(gomock.Any()).
		Do(func(run RunFuncOnceAt) {
			run("0/1 * * * * *", "test", func() {})
		})

	gone.
		Load(scheduler).
		Test(func(s *schedule) {
			s.isCluster = false
			_ = s.Start()
			err := s.Stop()
			assert.Nil(t, err)
		})
}

func Test_schedule_Start_WithTracer(t *testing.T) {
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
				println("test with tracer")
				mu.Lock()
				i++
				mu.Unlock()
			})
		})

	gone.
		NewApp(tracer.Load, Load).
		Load(scheduler).
		Test(func(s *schedule) {
			assert.NotNil(t, s.tracer)
			time.Sleep(2 * time.Second)
		})

	mu.Lock()
	assert.Equal(t, 2, i)
	mu.Unlock()
}

func Test_schedule_wrapFn_Panic(t *testing.T) {
	_ = os.Setenv("GONE_SCHEDULE_IN-CLUSTER", "false")
	defer func() {
		_ = os.Setenv("GONE_SCHEDULE_IN-CLUSTER", "true")
	}()

	controller := gomock.NewController(t)
	defer controller.Finish()
	scheduler := NewMockScheduler(controller)
	scheduler.EXPECT().
		Cron(gomock.Any()).
		Do(func(run RunFuncOnceAt) {
			run("0/1 * * * * *", "panic-test", func() {
				panic("test panic")
			})
		}).
		AnyTimes()

	gone.
		NewApp(Load).
		Load(scheduler).
		Test(func(s *schedule) {
			_ = s.Start()
			// 直接调用包装函数，测试panic恢复
			wrappedFn := s.wrapFn(func() {
				panic("test panic recovery")
			}, "panic-recovery")

			// 不应该导致测试崩溃
			wrappedFn()
		})
}
