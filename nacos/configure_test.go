package nacos

import (
	"strings"
	"testing"

	"github.com/gone-io/gone/v2"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	originViper "github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -destination=mock_configure_test.go -package=nacos github.com/gone-io/gone/v2 Configure,Logger

//go:generate mockgen -destination=mock_nacos_client_test.go -package=nacos github.com/nacos-group/nacos-sdk-go/v2/clients/config_client IConfigClient

func Test_configure_init(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	mockClient := NewMockIConfigClient(ctr)
	mockConfigure := NewMockConfigure(ctr)

	t.Run("Success Case", func(t *testing.T) {
		clientConfig := constant.ClientConfig{}
		serverConfigs := []constant.ServerConfig{{
			IpAddr: "localhost",
			Port:   8848,
		}}

		mockConfigure.EXPECT().
			Get("nacos.client", gomock.Any(), "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*constant.ClientConfig)) = clientConfig
				return nil
			})

		mockConfigure.EXPECT().
			Get("nacos.server", gomock.Any(), "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*[]constant.ServerConfig)) = serverConfigs
				return nil
			})

		c := configure{}
		client, err := c.init(mockConfigure, func(param vo.NacosClientParam) (iClient config_client.IConfigClient, err error) {
			return mockClient, nil
		})

		assert.NoError(t, err)
		assert.Equal(t, mockClient, client)
	})

	t.Run("Client Config Error", func(t *testing.T) {
		mockConfigure.EXPECT().
			Get("nacos.client", gomock.Any(), "").
			Return(gone.NewInnerError("client config error", gone.InjectError))

		c := configure{}
		client, err := c.init(mockConfigure, func(param vo.NacosClientParam) (iClient config_client.IConfigClient, err error) {
			return mockClient, nil
		})

		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "failed to get nacos client config")
	})

	t.Run("Server Config Error", func(t *testing.T) {
		clientConfig := constant.ClientConfig{}
		mockConfigure.EXPECT().
			Get("nacos.client", gomock.Any(), "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*constant.ClientConfig)) = clientConfig
				return nil
			})

		mockConfigure.EXPECT().
			Get("nacos.server", gomock.Any(), "").
			Return(gone.NewInnerError("server config error", gone.InjectError))

		c := configure{}
		client, err := c.init(mockConfigure, func(param vo.NacosClientParam) (iClient config_client.IConfigClient, err error) {
			return mockClient, nil
		})

		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "failed to get nacos server config")
	})

	t.Run("Client Creation Error", func(t *testing.T) {
		clientConfig := constant.ClientConfig{}
		serverConfigs := []constant.ServerConfig{{
			IpAddr: "localhost",
			Port:   8848,
		}}

		mockConfigure.EXPECT().
			Get("nacos.client", gomock.Any(), "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*constant.ClientConfig)) = clientConfig
				return nil
			})

		mockConfigure.EXPECT().
			Get("nacos.server", gomock.Any(), "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*[]constant.ServerConfig)) = serverConfigs
				return nil
			})

		c := configure{}
		client, err := c.init(mockConfigure, func(param vo.NacosClientParam) (iClient config_client.IConfigClient, err error) {
			return nil, gone.NewInnerError("client creation error", gone.InjectError)
		})

		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "client creation error")
	})
}

func Test_configure_getConfigContent(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	mockClient := NewMockIConfigClient(ctr)
	mockConfigure := NewMockConfigure(ctr)

	t.Run("Success Case", func(t *testing.T) {
		dataId := "test-data-id"
		groups := []confGroup{{
			Group:  "test-group",
			Format: "yaml",
		}}
		configContent := "key: value"

		mockConfigure.EXPECT().
			Get("nacos.dataId", gomock.Any(), "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*string)) = dataId
				return nil
			})

		mockConfigure.EXPECT().
			Get("nacos.groups", gomock.Any(), "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*[]confGroup)) = groups
				return nil
			})

		mockConfigure.EXPECT().
			Get("nacos.watch", gomock.Any(), "false").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*bool)) = true
				return nil
			})

		mockConfigure.EXPECT().
			Get("nacos.useLocalConfIfKeyNotExist", gomock.Any(), "true").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*bool)) = true
				return nil
			})

		mockClient.EXPECT().
			GetConfig(gomock.Any()).
			Return(configContent, nil)

		mockClient.EXPECT().ListenConfig(gomock.Any()).Return(nil)

		c := configure{}
		err := c.getConfigContent(mockConfigure, mockClient)

		assert.NoError(t, err)
		assert.Equal(t, dataId, c.dataId)
		assert.Equal(t, groups, c.groups)
		assert.True(t, c.watch)
		assert.True(t, c.useLocalConfIfKeyNotExist)
		assert.NotNil(t, c.viper)
		assert.NotNil(t, c.groupConfMap)
		assert.NotNil(t, c.keyMap)
	})

	t.Run("Empty DataId", func(t *testing.T) {
		mockConfigure.EXPECT().
			Get("nacos.dataId", gomock.Any(), "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*string)) = ""
				return nil
			})

		c := configure{}
		err := c.getConfigContent(mockConfigure, mockClient)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nacos config dataId is empty")
	})

	t.Run("Empty Groups", func(t *testing.T) {
		mockConfigure.EXPECT().
			Get("nacos.dataId", gomock.Any(), "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*string)) = "test-data-id"
				return nil
			})

		mockConfigure.EXPECT().
			Get("nacos.groups", gomock.Any(), "").
			Return(gone.NewInnerError("nacos config groups is empty", gone.InjectError))

		c := configure{}
		err := c.getConfigContent(mockConfigure, mockClient)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nacos config groups is empty")
	})

	t.Run("GetConfig Error", func(t *testing.T) {
		dataId := "test-data-id"
		groups := []confGroup{{
			Group:  "test-group",
			Format: "yaml",
		}}

		mockConfigure.EXPECT().
			Get("nacos.dataId", gomock.Any(), "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*string)) = dataId
				return nil
			})

		mockConfigure.EXPECT().
			Get("nacos.groups", gomock.Any(), "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*[]confGroup)) = groups
				return nil
			})

		mockConfigure.EXPECT().
			Get("nacos.watch", gomock.Any(), "false").
			Return(nil)

		mockConfigure.EXPECT().
			Get("nacos.useLocalConfIfKeyNotExist", gomock.Any(), "true").
			Return(nil)

		mockClient.EXPECT().
			GetConfig(gomock.Any()).
			Return("", gone.NewInnerError("failed to get config", gone.InjectError))

		c := configure{}
		err := c.getConfigContent(mockConfigure, mockClient)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get config")
	})

	t.Run("Invalid Config Format", func(t *testing.T) {
		dataId := "test-data-id"
		groups := []confGroup{{
			Group:  "test-group",
			Format: "yaml",
		}}
		configContent := "invalid:yaml:content"

		mockConfigure.EXPECT().
			Get("nacos.dataId", gomock.Any(), "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*string)) = dataId
				return nil
			})

		mockConfigure.EXPECT().
			Get("nacos.groups", gomock.Any(), "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*[]confGroup)) = groups
				return nil
			})

		mockConfigure.EXPECT().
			Get("nacos.watch", gomock.Any(), "false").
			Return(nil)

		mockConfigure.EXPECT().
			Get("nacos.useLocalConfIfKeyNotExist", gomock.Any(), "true").
			Return(nil)

		mockClient.EXPECT().
			GetConfig(gomock.Any()).
			Return(configContent, nil)

		c := configure{}
		err := c.getConfigContent(mockConfigure, mockClient)

		assert.Error(t, err)
	})
}

func Test_configure_OnChange(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	mockLogger := NewMockLogger(ctr)

	t.Run("Success Case", func(t *testing.T) {
		c := configure{
			logger: mockLogger,
			groups: []confGroup{
				{Group: "test-group", Format: "yaml"},
				{Group: "test-group-2", Format: "yaml"},
			},
			groupConfMap: make(map[string]*originViper.Viper),
			viper:        originViper.New(),
			keyMap:       make(map[string][]any),
			watch:        true,
		}

		// 初始化配置
		initialContent := "key1: value1\nkey2: value2"
		v := originViper.New()
		v.SetConfigType("yaml")
		err := v.ReadConfig(strings.NewReader(initialContent))
		assert.NoError(t, err)
		c.groupConfMap["test-group"] = v
		c.groupConfMap["test-group-2"] = originViper.New()
		c.viper.MergeConfigMap(v.AllSettings())

		var key1 string
		var key2 string
		_ = c.Get("key1", &key1, "")
		_ = c.Get("key2", &key2, "")

		assert.Equal(t, "value1", key1)
		assert.Equal(t, "value2", key2)

		// 模拟配置变更
		newContent := "key1: new-value1\nkey2: value2"
		c.OnChange("", "test-group", "", newContent)

		// 验证配置是否正确更新
		assert.Equal(t, "new-value1", key1)
		assert.Equal(t, "value2", key2)
	})

	t.Run("Invalid YAML Format", func(t *testing.T) {
		c := configure{
			logger: mockLogger,
			groups: []confGroup{
				{Group: "test-group", Format: "yaml"},
			},
			groupConfMap: make(map[string]*originViper.Viper),
			viper:        originViper.New(),
			keyMap:       make(map[string][]any),
		}

		mockLogger.EXPECT().
			Errorf(gomock.Any(), gomock.Any()).
			Times(1)

		// 传入无效的YAML内容
		invalidContent := "invalid:yaml:content"
		c.OnChange("", "test-group", "", invalidContent)

		// 验证配置未发生变化
		assert.Empty(t, c.groupConfMap["test-group"])
	})

	t.Run("Config Merge Error", func(t *testing.T) {
		c := configure{
			logger: mockLogger,
			groups: []confGroup{
				{Group: "test-group", Format: "yaml"},
				{Group: "test-group-2", Format: "yaml"},
			},
			groupConfMap: make(map[string]*originViper.Viper),
			viper:        originViper.New(),
			keyMap:       make(map[string][]any),
		}

		mockLogger.EXPECT().
			Errorf(gomock.Any(), gomock.Any()).
			Times(1)

		// 初始化配置
		initialContent := "key1: value1\nkey2: value2"
		v := originViper.New()
		v.SetConfigType("yaml")
		err := v.ReadConfig(strings.NewReader(initialContent))
		assert.NoError(t, err)
		c.groupConfMap["test-group-2"] = v

		// 传入无效的YAML内容到另一个组
		invalidContent := "invalid:yaml:content"
		c.OnChange("", "test-group", "", invalidContent)

		// 验证配置未发生变化
		assert.Empty(t, c.groupConfMap["test-group"])
	})

	t.Run("Value Update", func(t *testing.T) {
		var testValue string
		c := configure{
			logger: mockLogger,
			groups: []confGroup{
				{Group: "test-group", Format: "yaml"},
			},
			groupConfMap: make(map[string]*originViper.Viper),
			viper:        originViper.New(),
			keyMap:       make(map[string][]any),
		}

		// 初始化配置和keyMap
		initialContent := "key1: value1"
		v := originViper.New()
		v.SetConfigType("yaml")
		err := v.ReadConfig(strings.NewReader(initialContent))
		assert.NoError(t, err)
		c.groupConfMap["test-group"] = v
		c.viper.MergeConfigMap(v.AllSettings())
		c.keyMap["key1"] = []any{&testValue}

		// 模拟配置变更
		newContent := "key1: new-value1"
		c.OnChange("", "test-group", "", newContent)

		// 验证值是否被正确更新
		assert.Equal(t, "new-value1", testValue)
		assert.Equal(t, "new-value1", c.viper.GetString("key1"))
	})
}
