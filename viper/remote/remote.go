package remote

import (
	"github.com/gone-io/gone/v2"
	goneViper "github.com/gone-io/goner/viper"
	"github.com/spf13/viper"
	"time"
)

import _ "github.com/spf13/viper/remote"

type remoteConfigure struct {
	gone.Flag
	*goneViper.RemoteConfigure

	testFlag gone.TestFlag `gone:"*" option:"allowNil"` // testFlag only for test environment
	logger   gone.Logger   `gone:"*" option:"lazy"`
	w        *watcher      `gone:"*"`

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

	// Viper uses crypt to retrieve configuration from the K/V store, which means that you can store your configuration
	// values encrypted and have them automatically decrypted if you have the correct gpg keyring.
	// Encryption is optional.
	Keyring string //gpg keyring
}

var newViper = func() ViperInterface {
	return viper.New()
}
var newRemoteViper = newViper

var newGonerViper = goneViper.New

func (s *remoteConfigure) Init() error {
	return s.init(newGonerViper(s.testFlag), newViper())
}

func (s *remoteConfigure) init(localConfigure gone.Configure, v ViperInterface) (err error) {
	_ = localConfigure.Get("viper.remote.providers", &s.providers, "")
	_ = localConfigure.Get("viper.remote.watch", &s.watch, "false")
	_ = localConfigure.Get("viper.remote.watchDuration", &s.watchDuration, "5s")
	_ = localConfigure.Get("viper.remote.useLocalConfIfKeyNotExist", &s.useLocalConfIfKeyNotExist, "true")

	if s.watch {
		s.w.SetViper(v)
		go s.w.watch(s.watchDuration)
	}

	for _, p := range s.providers {
		v2 := newRemoteViper()
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
			s.w.addViper(v2)
		}
	}

	s.RemoteConfigure = goneViper.NewRemoteConfigure(
		v,
		localConfigure,
		s.useLocalConfIfKeyNotExist,
		s.w,
	)
	return nil
}

func (s *remoteConfigure) Notify(key string, callback gone.ConfWatchFunc) {
	s.w.Watch(key, callback)
}
