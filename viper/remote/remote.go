package remote

import (
	"github.com/gone-io/gone/v2"
	goneViper "github.com/gone-io/goner/viper"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/viper"
	"time"
)

import _ "github.com/spf13/viper/remote"

//go:generate mockgen -source=remote.go -destination=mock_viper_interface_test.go -package=remote
//go:generate mockgen -destination=mock_configure_test.go -package=remote github.com/gone-io/gone/v2 Configure

type ViperInterface interface {
	ReadRemoteConfig() error
	AllSettings() map[string]any
	Get(key string) any
	MergeConfigMap(settings map[string]any) error
	UnmarshalKey(key string, rawVal any, opts ...viper.DecoderConfigOption) error
	SetConfigType(configType string)
	AddRemoteProvider(provider string, endpoint string, path string) error
	AddSecureRemoteProvider(provider string, endpoint string, path string, keyring string) error
}

type remoteConfigure struct {
	gone.Flag

	// testFlag only for test environment
	testFlag gone.TestFlag `gone:"*" option:"allowNil"`
	logger   gone.Logger   `gone:"*" option:"lazy"`

	localConfigure gone.Configure
	viper          ViperInterface
	remoteVipers   []ViperInterface
	keyMap         map[string][]any

	providers                 []Provider    //`gone:"config,viper.remote.providers"`
	watch                     bool          //`gone:"config,viper.remote.watch"`
	watchDuration             time.Duration //`gone:"config,viper.remote.watchDuration"`
	useLocalConfIfKeyNotExist bool          //`gone:"config,viper.remote.useLocalConfIfKeyNotExist"`

}

type Provider struct {
	Provider   string
	Endpoint   string
	Path       string
	ConfigType string

	//Viper uses crypt to retrieve configuration from the K/V store, which means that you can store your configuration values encrypted and have them automatically decrypted if you have the correct gpg keyring. Encryption is optional.
	Keyring string //gpg keyring
}

var newViper = func() ViperInterface {
	return viper.New()
}
var newGonerViper = goneViper.New

func (s *remoteConfigure) Init() error {
	s.localConfigure = newGonerViper(s.testFlag)
	s.viper = newViper()
	s.keyMap = make(map[string][]any)
	return s.init(s.localConfigure, s.viper)
}

func (s *remoteConfigure) doWatch(duration time.Duration) {
	for {
		time.Sleep(duration)

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

			if compare(oldValue, newValue) {
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
}

func compare(a, b any) (equal bool) {
	return cmp.Equal(a, b)
}

func (s *remoteConfigure) init(localConfigure gone.Configure, v ViperInterface) (err error) {
	_ = localConfigure.Get("viper.remote.providers", &s.providers, "")
	_ = localConfigure.Get("viper.remote.watch", &s.watch, "false")
	_ = localConfigure.Get("viper.remote.watchDuration", &s.watchDuration, "5s")
	_ = localConfigure.Get("viper.remote.useLocalConfIfKeyNotExist", &s.useLocalConfIfKeyNotExist, "true")

	for _, p := range s.providers {
		v2 := newViper()
		v2.SetConfigType(p.ConfigType)

		if p.Keyring == "" {
			err = v2.AddRemoteProvider(p.Provider, p.Endpoint, p.Path)
		} else {
			err = v2.AddSecureRemoteProvider(p.Provider, p.Endpoint, p.Path, p.Keyring)
		}
		if err != nil {
			return gone.ToError(err)
		}
		err = v2.ReadRemoteConfig()
		if err != nil {
			return gone.ToError(err)
		}
		err = v.MergeConfigMap(v2.AllSettings())
		if err != nil {
			return gone.ToError(err)
		}
		if s.watch {
			s.remoteVipers = append(s.remoteVipers, v2)
		}
	}
	if s.watch {
		go s.doWatch(s.watchDuration)
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
	if v == nil || v == "" {
		return s.localConfigure.Get(key, value, defaultVal)
	}
	err := s.viper.UnmarshalKey(key, value)
	if err != nil && s.useLocalConfIfKeyNotExist {
		return s.localConfigure.Get(key, value, defaultVal)
	}
	return gone.ToError(err)
}
