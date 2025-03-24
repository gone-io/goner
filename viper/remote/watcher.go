package remote

import (
	"github.com/gone-io/gone/v2"
	"github.com/google/go-cmp/cmp"
	"time"
)

type watcher struct {
	gone.Flag
	viper        ViperInterface
	remoteVipers []ViperInterface
	keyMap       map[string][]any
	logger       gone.Logger `gone:"*" option:"lazy"`
}

func (s *watcher) Init() {
	s.remoteVipers = make([]ViperInterface, 0)
	s.keyMap = make(map[string][]any)
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

func (s *watcher) doWatch() {
	all := newViper()
	for _, v2 := range s.remoteVipers {
		if err := v2.ReadRemoteConfig(); err != nil {
			s.logger.Warnf("try to read remote config err:%v\n", err)
			return
		}
		if err := all.MergeConfigMap(v2.AllSettings()); err != nil {
			s.logger.Warnf("try to merge remote config err:%v\n", err)
			return
		}
	}

	needMerge := false
	for useK, values := range s.keyMap {
		oldValue := s.viper.Get(useK)
		newValue := all.Get(useK)

		if cmp.Equal(oldValue, newValue) {
			continue
		}

		for _, ani := range values {
			if err := all.UnmarshalKey(useK, ani); err != nil {
				s.logger.Warnf("try to set `%s` value err:%v\n", useK, err)
			}
		}
		needMerge = true
	}
	if needMerge {
		if err := s.viper.MergeConfigMap(all.AllSettings()); err != nil {
			s.logger.Warnf("try to merge remote config err:%v\n", err)
		}
	}
}
