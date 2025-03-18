package apollo

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/internal/json"
	viper "github.com/gone-io/goner/viper"
	"reflect"
	"strings"
)

type apolloClient struct {
	gone.Flag
	localConfigure gone.Configure
	apolloClient   agollo.Client

	changeListener *changeListener `gone:"*"`

	// testFlag only for test environment
	testFlag gone.TestFlag `gone:"*" option:"allowNil"`

	//lazy fill，resolve circular dependency
	logger gone.Logger `gone:"*" option:"lazy"`

	appId                     string //`gone:"config,apollo.appId"`
	cluster                   string //`gone:"config,apollo.cluster"`
	ip                        string //`gone:"config,apollo.ip"`
	namespace                 string //`gone:"config,apollo.namespace"`
	secret                    string //`gone:"config,apollo.secret"`
	isBackupConfig            bool   //`gone:"config,apollo.isBackupConfig"`
	watch                     bool   //`gone:"config,apollo.watch"`
	useLocalConfIfKeyNotExist bool   //`gone:"config,apollo.useLocalConfIfKeyNotExist"`
}

func (s *apolloClient) init(localConfigure gone.Configure, startWithConfig func(loadAppConfig func() (*config.AppConfig, error)) (agollo.Client, error)) {
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
			panic(err)
		}
	}

	c := &config.AppConfig{
		AppID:          s.appId,
		Cluster:        s.cluster,
		IP:             s.ip,
		NamespaceName:  s.namespace,
		IsBackupConfig: s.isBackupConfig,
		Secret:         s.secret,
	}
	client, err := startWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		panic(err)
	}
	s.apolloClient = client
	if s.watch {
		client.AddChangeListener(s.changeListener)
	}
}

func (s *apolloClient) Init() {
	s.localConfigure = viper.New(s.testFlag)
	s.init(s.localConfigure, agollo.StartWithConfig)
}

func (s *apolloClient) Get(key string, v any, defaultVal string) error {
	if s.watch {
		s.changeListener.Put(key, v)
	}

	if s.apolloClient == nil {
		return s.localConfigure.Get(key, v, defaultVal)
	}

	namespaces := strings.Split(s.namespace, ",")
	for _, ns := range namespaces {
		cache := s.apolloClient.GetConfigCache(ns)
		if cache != nil {
			if value, err := cache.Get(key); err == nil {
				err = setValue(v, value)
				if err != nil {
					s.warnf("try to set `%s` value err:%v\n", key, err)
				} else {
					return nil
				}
			} else {
				s.warnf("get `%s` value from apollo ns(%s) err:%v\n", key, ns, err)
			}
		}
	}
	if s.useLocalConfIfKeyNotExist {
		return s.localConfigure.Get(key, v, defaultVal)
	}
	return nil
}

func setValue(v any, value any) error {
	if str, ok := value.(string); ok {
		return gone.ToError(gone.SetValue(reflect.ValueOf(v), v, str))
	} else {
		marshal, err := json.Marshal(value)
		if err != nil {
			return gone.ToError(err)
		}
		return gone.ToError(gone.SetValue(reflect.ValueOf(v), v, string(marshal)))
	}
}

func (s *apolloClient) warnf(format string, args ...any) {
	if s.logger != nil {
		s.logger.Warnf(format, args...)
	}
}
