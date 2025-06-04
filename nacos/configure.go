package nacos

import (
	"github.com/go-viper/encoding/javaproperties"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	goneViper "github.com/gone-io/goner/viper"
	"github.com/google/go-cmp/cmp"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
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
	watchMap     map[string]gone.ConfWatchFunc
	viper        *viper.Viper
	keyMap       map[string][]any
}

func (s *configure) newConf(group, content string) *viper.Viper {
	var format = "yaml"
	for _, v := range s.groups {
		if v.Group == group {
			format = v.Format
			break
		}
	}
	v2, _ := newViperByFormat(format)
	v2.SetConfigType(format)
	err := v2.ReadConfig(strings.NewReader(content))
	g.ErrorPrinter(s.logger, err, "failed to read config file")
	return v2
}

func (s *configure) OnChange(_, group, _, content string) {
	s.groupConfMap[group] = s.newConf(group, content)

	tmpViper := viper.New()
	for _, gro := range s.groups {
		err := tmpViper.MergeConfigMap(s.groupConfMap[gro.Group].AllSettings())
		g.ErrorPrinter(s.logger, err, "failed to merge config")
	}

	for k, values := range s.keyMap {
		oldValue := s.viper.Get(k)
		newValue := tmpViper.Get(k)
		if compare(oldValue, newValue) {
			continue
		}
		for _, ani := range values {
			err := tmpViper.UnmarshalKey(k, ani)
			g.ErrorPrinter(s.logger, err, "try to unmarshal key")
		}
	}

	for key, value := range s.watchMap {
		oldVal := s.viper.Get(key)
		newVal := tmpViper.Get(key)
		if !compare(oldVal, newVal) {
			value(oldVal, newVal)
		}
	}
	s.viper = tmpViper
}
func compare(a, b any) (equal bool) {
	return cmp.Equal(a, b)
}

func newViperByFormat(format string) (*viper.Viper, error) {
	if format != "properties" {
		return viper.New(), nil
	}

	codecRegistry := viper.NewCodecRegistry()
	codec := &javaproperties.Codec{}
	err := codecRegistry.RegisterCodec("properties", codec)
	if err != nil {
		return nil, gone.ToError(err)
	}
	return viper.NewWithOptions(viper.WithCodecRegistry(codecRegistry)), nil
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

	for _, group := range s.groups {
		param := vo.ConfigParam{
			DataId: s.dataId,
			Group:  group.Group,
		}

		content, err := client.GetConfig(param)
		if err != nil {
			return gone.ToError(err)
		}
		v2, err := newViperByFormat(group.Format)
		if err != nil {
			return gone.ToError(err)
		}
		v2.SetConfigType(group.Format)
		err = v2.ReadConfig(strings.NewReader(content))
		if err != nil {
			return gone.ToError(err)
		}
		s.groupConfMap[group.Group] = v2
		err = s.viper.MergeConfigMap(v2.AllSettings())
		if err != nil {
			return gone.ToError(err)
		}

		if s.watch {
			param.OnChange = s.OnChange
			err = client.ListenConfig(param)
			if err != nil {
				return gone.ToError(err)
			}
		}
	}
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

var newConfigClient = clients.NewConfigClient

func (s *configure) Init() {
	s.watchMap = make(map[string]gone.ConfWatchFunc)
	s.localConfigure = goneViper.New(s.testFlag)
	var err error

	s.client, err = s.init(s.localConfigure, newConfigClient)
	g.PanicIfErr(err)
	err = s.getConfigContent(s.localConfigure, s.client)
	g.PanicIfErr(err)
}

func (s *configure) Get(key string, value any, defaultVal string) error {
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

func (s *configure) Notify(key string, callback gone.ConfWatchFunc) {
	s.watchMap[key] = callback
}
