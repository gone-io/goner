package schedule

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/tracer"
	"github.com/robfig/cron/v3"
	"time"
)

var load = gone.OnceLoad(func(loader gone.Loader) error {
	err := tracer.Load(loader)
	if err != nil {
		return gone.ToError(err)
	}
	return loader.Load(&schedule{})
})

func Load(loader gone.Loader) error {
	return load(loader)
}

// Priest Deprecated, use Load instead
func Priest(loader gone.Loader) error {
	return Load(loader)
}

type schedule struct {
	gone.Flag

	log        gone.Logger      `gone:"gone-logger"`
	tracer     tracer.Tracer    `gone:"gone-tracer"`
	schedulers []Scheduler      `gone:"*"`
	gKeeper    gone.GonerKeeper `gone:"*"`

	isCluster   bool          `gone:"config,schedule.in-cluster=true"`
	lockTime    time.Duration `gone:"config,schedule.lockTime,default=10s"`
	checkPeriod time.Duration `gone:"config,schedule.checkPeriod,default=2s"`

	cronTab *cron.Cron
	locker  DoLocker
}

func (s *schedule) Start() error {
	if len(s.schedulers) == 0 {
		s.log.Warnf("no scheduler found")
		return nil
	}

	if s.isCluster {
		locker := s.gKeeper.GetGonerByType(gone.GetInterfaceType(new(DoLocker)))
		if locker == nil {
			return gone.ToError("in cluster mod, must load a `DoLocker`. you can use `goner/redis`.")
		}
		s.locker = locker.(DoLocker)
	} else {
		s.log.Warnf("`schedule` is running in single instance mod.")
	}

	s.cronTab = cron.New(cron.WithSeconds())

	for _, o := range s.schedulers {
		o.Cron(func(spec string, jobName JobName, fn func()) {
			fnWrap := func() {
				s.tracer.SetTraceId("", func() {
					defer func() {
						if err := recover(); err != nil {
							e := gone.NewInnerErrorSkip(fmt.Sprintf("panic: %v", err), gone.PanicError, 3)
							s.log.Errorf("%v", e)
						}
					}()
					if s.locker != nil {
						lockKey := fmt.Sprintf("lock-job:%s", jobName)
						err := s.locker.LockAndDo(lockKey, fn, s.lockTime, s.checkPeriod)
						if err != nil {
							s.log.Warnf("cron get lock err:%v", err)
						}
					} else {
						fn()
					}
				})
			}
			_, err := s.cronTab.AddFunc(spec, fnWrap)

			if err != nil {
				panic("cron.AddFunc for " + string(jobName) + " err:" + err.Error())
			}
			s.log.Infof("Add cron item: %s => %s : %s", spec, jobName, gone.GetFuncName(fn))
		})
	}
	s.cronTab.Start()
	return nil
}

func (s *schedule) Stop() error {
	s.cronTab.Stop()
	return nil
}
