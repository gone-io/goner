package schedule

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/robfig/cron/v3"
	"time"
)

type schedule struct {
	gone.Flag

	logger      gone.Logger   `gone:"*"`
	schedulers  []Scheduler   `gone:"*"`
	locker      DoLocker      `gone:"*" option:"allowNil"`
	tracer      g.Tracer      `gone:"*" option:"allowNil"`
	isCluster   bool          `gone:"config,schedule.in-cluster=true"`
	lockTime    time.Duration `gone:"config,schedule.lockTime,default=10s"`
	checkPeriod time.Duration `gone:"config,schedule.checkPeriod,default=2s"`

	cronTab *cron.Cron
}

func (s *schedule) Start() error {
	if len(s.schedulers) == 0 {
		s.logger.Warnf("no scheduler found")
		return nil
	}

	if s.isCluster {
		if s.locker == nil {
			return gone.ToError("in cluster mod, must load a `DoLocker`. you can use `goner/redis`.")
		}
	} else {
		s.logger.Warnf("`schedule` is running in single instance mod.")
	}

	s.cronTab = cron.New(cron.WithSeconds())

	for _, o := range s.schedulers {
		o.Cron(func(spec string, jobName JobName, fn func()) {
			_, err := s.cronTab.AddFunc(spec, s.wrapFn(fn, jobName))
			if err != nil {
				panic("cron.AddFunc for " + string(jobName) + " err:" + err.Error())
			}
			s.logger.Infof("Add cron item: %s => %s : %s", spec, jobName, gone.GetFuncName(fn))
		})
	}
	s.cronTab.Start()
	return nil
}

func (s *schedule) wrapFn(fn func(), jobName JobName) func() {
	return func() {
		f := func() {
			defer func() {
				if err := recover(); err != nil {
					e := gone.NewInnerErrorSkip(fmt.Sprintf("panic: %v", err), gone.PanicError, 3)
					s.logger.Errorf("%v", e)
				}
			}()
			if s.locker != nil {
				lockKey := fmt.Sprintf("lock-job:%s", jobName)
				err := s.locker.LockAndDo(lockKey, fn, s.lockTime, s.checkPeriod)
				if err != nil {
					s.logger.Warnf("cron get lock err:%v", err)
				}
			} else {
				fn()
			}
		}
		if s.tracer != nil {
			s.tracer.SetTraceId("", f)
		} else {
			f()
		}
	}
}

func (s *schedule) Stop() error {
	if s.cronTab == nil {
		return nil
	}
	s.cronTab.Stop()
	return nil
}
