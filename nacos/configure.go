package nacos

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	goneViper "github.com/gone-io/goner/viper"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

type confGroup struct {
	Group  string
	Format string
}

type configure struct {
	gone.Flag
	// testFlag only for test environment
	testFlag gone.TestFlag `gone:"*" option:"allowNil"`
	logger   gone.Logger   `gone:"*" option:"lazy"`

	localConfigure gone.Configure
	client         config_client.IConfigClient

	dataId                    string      //nacos.dataId
	groups                    []confGroup // nacos.groups
	watch                     bool        //`gone:"config,nacos.watch"`
	useLocalConfIfKeyNotExist bool        //`gone:"config,nacos.useLocalConfIfKeyNotExist"`

	groupConfMap map[string]*viper.Viper
	viper        *viper.Viper
	keyMap       map[string][]any
}

func (s *configure) OnChange(namespace, group, dataId, content string) {
	var format = "yaml"
	for _, v := range s.groups {
		if v.Group == group {
			format = v.Format
			break
		}
	}

	v2 := viper.New()
	v2.SetConfigType(format)
	err := v2.ReadConfig(strings.NewReader(content))
	if err != nil {
		s.logger.Errorf("failed to read config file, err: %v", err)
		return
	}
	s.groupConfMap[group] = v2
	tmpViper := viper.New()
	for _, g := range s.groups {
		err = tmpViper.MergeConfigMap(s.groupConfMap[g.Group].AllSettings())
		if err != nil {
			s.logger.Errorf("failed to merge config, err: %v", err)
			return
		}
	}

	for k, values := range s.keyMap {
		oldValue := s.viper.Get(k)
		newValue := tmpViper.Get(k)
		if oldValue != newValue {
			for _, ani := range values {
				err = gone.SetValue(reflect.ValueOf(ani), ani, s.viper.GetString(k))
				if err != nil {
					s.logger.Warnf("try to set `%s` value err:%v\n", k, err)
				}
			}
			s.viper.Set(k, newValue)
		}
	}
}

func (s *configure) getConfigContent(localConfigure gone.Configure, client config_client.IConfigClient) (err error) {
	err = localConfigure.Get("nacos.dataId", &s.dataId, "")
	if err != nil {
		return
	}
	if s.dataId == "" {
		return gone.NewInnerError("nacos config dataId is empty", gone.InjectError)
	}
	err = localConfigure.Get("nacos.groups", &s.groups, "")
	if err != nil {
		return gone.NewInnerError("nacos config groups is empty", gone.InjectError)
	}
	_ = localConfigure.Get("nacos.watch", &s.watch, "false")
	_ = localConfigure.Get("nacos.useLocalConfIfKeyNotExist", &s.useLocalConfIfKeyNotExist, "true")

	s.groupConfMap = make(map[string]*viper.Viper)
	s.viper = viper.New()

	for _, g := range s.groups {
		param := vo.ConfigParam{
			DataId: s.dataId,
			Group:  g.Group,
		}
		if s.watch {
			param.OnChange = s.OnChange
		}
		content, err := client.GetConfig(param)
		if err != nil {
			return gone.ToError(err)
		}
		v2 := viper.New()
		v2.SetConfigType(g.Format)
		err = v2.ReadConfig(strings.NewReader(content))
		if err != nil {
			return gone.ToError(err)
		}
		s.groupConfMap[g.Group] = v2
		err = s.viper.MergeConfigMap(v2.AllSettings())
		if err != nil {
			return gone.ToError(err)
		}
	}

	keys := s.viper.AllKeys()
	fmt.Printf("%#+v", keys)

	s.keyMap = make(map[string][]any)
	return nil
}

func (s *configure) init(
	config gone.Configure,
	createClient func(param vo.NacosClientParam) (iClient config_client.IConfigClient, err error),
) (config_client.IConfigClient, error) {
	var clientConfig constant.ClientConfig
	err := config.Get("nacos.client", &clientConfig, "")
	if err != nil {
		return nil, gone.NewInnerError("failed to get nacos client config", gone.InjectError)
	}

	var serverConfigs []constant.ServerConfig
	err = config.Get("nacos.server", &serverConfigs, "")
	if err != nil {
		return nil, gone.NewInnerError("failed to get nacos server config", gone.InjectError)
	}
	configClient, err := createClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	return configClient, gone.ToError(err)
}

func (s *configure) Init() (err error) {
	s.localConfigure = goneViper.New(s.testFlag)
	s.client, err = s.init(s.localConfigure, clients.NewConfigClient)
	if err != nil {
		return
	}
	return s.getConfigContent(s.localConfigure, s.client)
}

func (s *configure) Get(key string, v any, defaultVal string) error {
	if s.watch {
		s.keyMap[key] = append(s.keyMap[key], v)
	}

	if s.client == nil || s.viper == nil {
		return s.localConfigure.Get(key, v, defaultVal)
	}

	err := gone.SetValue(reflect.ValueOf(v), v, s.viper.GetString(key))
	if err != nil {
		s.logger.Warnf("try to set `%s` value err:%v\n", key, err)
	}
	if s.useLocalConfIfKeyNotExist {
		return s.localConfigure.Get(key, v, defaultVal)
	}
	return nil
}
