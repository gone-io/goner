package viper

import (
	"github.com/gone-io/gone/v2"
	"github.com/spf13/viper"
)

type WatcherKeeper interface {
	Put(key string, value any)
}

type KeyGetter interface {
	Get(key string) any
	UnmarshalKey(key string, value any, opts ...viper.DecoderConfigOption) error
}

type RemoteConfigure struct {
	gone.Flag
	Viper   KeyGetter
	watcher WatcherKeeper
	local   gone.Configure

	useLocalConfIfKeyNotExist bool
}

func (s *RemoteConfigure) Get(key string, value any, defaultVal string) error {
	if s.watcher != nil {
		s.watcher.Put(key, value)
	}

	if s.Viper == nil {
		return s.local.Get(key, value, defaultVal)
	}

	v := s.Viper.Get(key)
	if v == nil || v == "" {
		return s.local.Get(key, value, defaultVal)
	}
	err := s.Viper.UnmarshalKey(key, value)
	if err != nil && s.useLocalConfIfKeyNotExist {
		return s.local.Get(key, value, defaultVal)
	}
	return gone.ToError(err)
}

func NewRemoteConfigure(viper KeyGetter, localConfigure gone.Configure, useLocalConfIfKeyNotExist bool, watcher WatcherKeeper) *RemoteConfigure {
	return &RemoteConfigure{
		Viper:                     viper,
		local:                     localConfigure,
		useLocalConfIfKeyNotExist: useLocalConfIfKeyNotExist,
		watcher:                   watcher,
	}
}
