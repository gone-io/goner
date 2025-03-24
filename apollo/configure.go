package apollo

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/viper"
	originViper "github.com/spf13/viper"
	"strings"
)

type apolloConfigure struct {
	gone.Flag
	client   agollo.Client
	testFlag gone.TestFlag `gone:"*" option:"allowNil"`
	logger   gone.Logger   `gone:"*" option:"lazy"`

	*viper.RemoteConfigure

	changeListener *changeListener `gone:"*"`

	appId                     string //`gone:"config,apollo.appId"`
	cluster                   string //`gone:"config,apollo.cluster"`
	ip                        string //`gone:"config,apollo.ip"`
	namespace                 string //`gone:"config,apollo.namespace"`
	secret                    string //`gone:"config,apollo.secret"`
	isBackupConfig            bool   //`gone:"config,apollo.isBackupConfig"`
	watch                     bool   //`gone:"config,apollo.watch"`
	useLocalConfIfKeyNotExist bool   //`gone:"config,apollo.useLocalConfIfKeyNotExist"`
}

func (s *apolloConfigure) init(localConfigure gone.Configure) (*config.AppConfig, error) {
	type tuple struct {
		v          any
		defaultVal string
	}

	m := map[string]*tuple{
		"apollo.appId":                     {v: &s.appId, defaultVal: ""},
		"apollo.cluster":                   {v: &s.cluster, defaultVal: "default"},
		"apollo.ip":                        {v: &s.ip, defaultVal: ""},
		"apollo.namespace":                 {v: &s.namespace, defaultVal: "application"},
		"apollo.secret":                    {v: &s.secret, defaultVal: ""},
		"apollo.isBackupConfig":            {v: &s.isBackupConfig, defaultVal: "true"},
		"apollo.watch":                     {v: &s.watch, defaultVal: "false"},
		"apollo.useLocalConfIfKeyNotExist": {v: &s.useLocalConfIfKeyNotExist, defaultVal: "true"},
	}
	for k, t := range m {
		err := localConfigure.Get(k, t.v, t.defaultVal)
		if err != nil {
			return nil, gone.ToError(err)
		}
	}

	return &config.AppConfig{
		AppID:          s.appId,
		Cluster:        s.cluster,
		IP:             s.ip,
		NamespaceName:  s.namespace,
		IsBackupConfig: s.isBackupConfig,
		Secret:         s.secret,
	}, nil
}

var startWithConfig = agollo.StartWithConfig

func (s *apolloConfigure) Init() error {
	configure := viper.New(s.testFlag)
	appConfig, err := s.init(configure)
	if err != nil {
		return gone.ToError(err)
	}

	s.client, err = startWithConfig(func() (*config.AppConfig, error) {
		return appConfig, nil
	})

	if err != nil {
		return gone.ToError(err)
	}

	namespaces := strings.Split(s.namespace, ",")

	total := originViper.New()
	for _, ns := range namespaces {
		cache := s.client.GetConfigCache(ns)
		if cache != nil {
			if s.watch {
				v := originViper.New()
				cache.Range(func(key, value interface{}) bool {
					if k, ok := key.(string); ok {
						v.Set(k, value)
					}
					return true
				})
				err = total.MergeConfigMap(v.AllSettings())
				if err != nil {
					return gone.ToError(err)
				}
				s.changeListener.AddViper(ns, v)
			} else {
				cache.Range(func(key, value interface{}) bool {
					if k, ok := key.(string); ok {
						total.Set(k, value)
					}
					return true
				})
			}
		}
	}
	s.RemoteConfigure = viper.NewRemoteConfigure(total, configure, true, s.changeListener)
	if s.watch {
		s.client.AddChangeListener(s.changeListener)
		s.changeListener.SetViper(total)
	}
	return nil
}
