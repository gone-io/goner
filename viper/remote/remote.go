package remote

import (
	"github.com/gone-io/gone/v2"
	goneViper "github.com/gone-io/goner/viper"
	"github.com/spf13/viper"
	"reflect"
	"time"
)

import _ "github.com/spf13/viper/remote"

type remoteConfigure struct {
	gone.Flag

	// testFlag only for test environment
	testFlag gone.TestFlag `gone:"*" option:"allowNil"`
	logger   gone.Logger   `gone:"*" option:"lazy"`

	localConfigure gone.Configure
	viper          *viper.Viper
	keyMap         map[string][]any

	providers                 []Provider //`gone:"config,viper.remote.providers"`
	configType                string     //`gone:"config,viper.remote.type"`
	watch                     bool       //`gone:"config,viper.remote.watch"`
	useLocalConfIfKeyNotExist bool       //`gone:"config,viper.remote.useLocalConfIfKeyNotExist"`

}

type Provider struct {
	Provider  string
	Endpoint  string
	Path      string
	SecretKey string
}

func (s *remoteConfigure) Init() error {
	s.localConfigure = goneViper.New(s.testFlag)
	s.viper = viper.New()
	return s.init(s.localConfigure, s.viper)
}

func (s *remoteConfigure) doWatch(v *viper.Viper) {
	for {
		err := v.WatchRemoteConfigOnChannel()
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}

		for k, values := range s.keyMap {
			for _, ani := range values {
				err = gone.SetValue(reflect.ValueOf(ani), ani, v.GetString(k))
				if err != nil {
					s.logger.Warnf("try to set `%s` value err:%v\n", k, err)
				}
			}
		}
	}
}

func (s *remoteConfigure) init(localConfigure gone.Configure, v *viper.Viper) (err error) {
	_ = localConfigure.Get("viper.remote.providers", &s.providers, "")
	_ = localConfigure.Get("viper.remote.type", &s.configType, "")
	_ = localConfigure.Get("viper.remote.watch", &s.watch, "false")
	_ = localConfigure.Get("viper.remote.useLocalConfIfKeyNotExist", &s.useLocalConfIfKeyNotExist, "true")

	v.SetConfigType(s.configType)

	for _, p := range s.providers {
		if p.SecretKey == "" {
			err = v.AddRemoteProvider(p.Provider, p.Endpoint, p.Path)
		} else {
			err = v.AddSecureRemoteProvider(p.Provider, p.Endpoint, p.Path, p.SecretKey)
		}
		if err != nil {
			return gone.ToError(err)
		}
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		return gone.ToError(err)
	}

	if s.watch {
		s.keyMap = make(map[string][]any)
		go s.doWatch(v)
	}
	return nil
}

func (s *remoteConfigure) Get(key string, value any, defaultVal string) error {
	if s.watch {
		s.keyMap[key] = append(s.keyMap[key], value)
	}

	if s.viper == nil {
		return s.localConfigure.Get(key, value, defaultVal)
	}

	v := s.viper.Get(key)
	if value == nil || value == "" {
		return s.localConfigure.Get(key, value, defaultVal)
	}
	err := s.viper.UnmarshalKey(key, value)
	if err != nil && s.useLocalConfIfKeyNotExist {
		return s.localConfigure.Get(key, v, defaultVal)
	}
	return gone.ToError(err)
}
