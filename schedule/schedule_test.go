package schedule

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/redis"
	"github.com/gone-io/goner/tracer"
	gone_viper "github.com/gone-io/goner/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"sync"
	"testing"
	"time"
)

type locker struct {
	gone.Flag
	redis.Locker
}

func (l *locker) LockAndDo(key string, fn func(), lockTime, checkPeriod time.Duration) (err error) {
	fn()
	return nil
}

func Test_schedule_Start(t *testing.T) {
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
		NewApp(tracer.Priest, Load, redis.Load, gone_viper.Load, func(loader gone.Loader) error {
			return loader.Load(scheduler)
		}).
		Test(func(s *schedule) {
			assert.True(t, s.isCluster)
			time.Sleep(2 * time.Second)

		})

	mu.Lock()
	assert.Equal(t, 2, i)
	mu.Unlock()
}
