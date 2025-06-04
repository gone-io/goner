package remote

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/google/go-cmp/cmp"
	"time"
)

type watcher struct {
	gone.Flag
	viper        ViperInterface
	remoteVipers []ViperInterface
	keyMap       map[string][]any
	watchMap     map[string]gone.ConfWatchFunc
	logger       gone.Logger `gone:"*" option:"lazy"`
}

func (s *watcher) Init() {
	s.remoteVipers = make([]ViperInterface, 0)
	s.keyMap = make(map[string][]any)
	s.watchMap = make(map[string]gone.ConfWatchFunc)
}

func (s *watcher) SetViper(v ViperInterface) {
	s.viper = v
}

func (s *watcher) Put(key string, value any) {
	s.keyMap[key] = append(s.keyMap[key], value)
}

func (s *watcher) addViper(v2 ViperInterface) {
	s.remoteVipers = append(s.remoteVipers, v2)
}

func (s *watcher) watch(duration time.Duration) {
	for {
		time.Sleep(duration)
		s.doWatch()
	}
}

func compare(a, b any) (equal bool) {
	return cmp.Equal(a, b)
}

func (s *watcher) doWatch() {
	tmpViper := newViper()
	for _, v2 := range s.remoteVipers {
		if err := v2.ReadRemoteConfig(); err != nil {
			s.logger.Warnf("try to read remote config err:%v\n", err)
			return
		}
		if err := tmpViper.MergeConfigMap(v2.AllSettings()); err != nil {
			s.logger.Warnf("try to merge remote config err:%v\n", err)
			return
		}
	}

	needMerge := false
	for useK, values := range s.keyMap {
		oldValue := s.viper.Get(useK)
		newValue := tmpViper.Get(useK)

		if compare(oldValue, newValue) {
			continue
		}

		for _, ani := range values {
			err := tmpViper.UnmarshalKey(useK, ani)
			g.ErrorPrinter(s.logger, err, "try to unmarshal key")
		}
		needMerge = true
	}
	for key, value := range s.watchMap {
		oldVal := s.viper.Get(key)
		newVal := tmpViper.Get(key)
		if !compare(oldVal, newVal) {
			value(oldVal, newVal)
			needMerge = true
		}
	}

	if needMerge {
		s.viper = tmpViper
	}
}

func (s *watcher) Watch(key string, callback gone.ConfWatchFunc) {
	s.watchMap[key] = callback
}
