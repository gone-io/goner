package apollo

import (
	"errors"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"testing"

	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestApolloClient_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟对象
	localConfigure := NewMockConfigure(ctrl)

	// 设置模拟对象的行为
	localConfigure.EXPECT().Get("apollo.appId", gomock.Any(), "").Return(nil).Do(
		func(key string, v any, defaultVal string) {
			*(v.(*string)) = "testApp"
		},
	)
	localConfigure.EXPECT().Get("apollo.cluster", gomock.Any(), "default").Return(nil).Do(
		func(key string, v any, defaultVal string) {
			*(v.(*string)) = "default"
		},
	)
	localConfigure.EXPECT().Get("apollo.ip", gomock.Any(), "").Return(nil).Do(
		func(key string, v any, defaultVal string) {
			*(v.(*string)) = "http://localhost:8080"
		},
	)
	localConfigure.EXPECT().Get("apollo.namespace", gomock.Any(), "application").Return(nil).Do(
		func(key string, v any, defaultVal string) {
			*(v.(*string)) = "application"
		},
	)
	localConfigure.EXPECT().Get("apollo.secret", gomock.Any(), "").Return(nil).Do(
		func(key string, v any, defaultVal string) {
			*(v.(*string)) = "secret"
		},
	)
	localConfigure.EXPECT().Get("apollo.isBackupConfig", gomock.Any(), "true").Return(nil).Do(
		func(key string, v any, defaultVal string) {
			*(v.(*bool)) = true
		},
	)
	localConfigure.EXPECT().Get("apollo.watch", gomock.Any(), "false").Return(nil).Do(
		func(key string, v any, defaultVal string) {
			*(v.(*bool)) = false
		},
	)
	localConfigure.EXPECT().Get("apollo.useLocalConfIfKeyNotExist", gomock.Any(), "true").Return(nil).Do(
		func(key string, v any, defaultVal string) {
			*(v.(*bool)) = true
		},
	)

	mockClient := NewMockClient(ctrl)

	// 创建apolloClient实例
	client := &apolloClient{
		changeListener: &changeListener{},
	}
	client.localConfigure = localConfigure

	// 执行初始化
	client.init(localConfigure, func(loadAppConfig func() (*config.AppConfig, error)) (agollo.Client, error) {
		return mockClient, nil
	})

	// 验证配置是否正确读取
	assert.Equal(t, "testApp", client.appId)
	assert.Equal(t, "default", client.cluster)
	assert.Equal(t, "http://localhost:8080", client.ip)
	assert.Equal(t, "application", client.namespace)
	assert.Equal(t, "secret", client.secret)
	assert.Equal(t, true, client.isBackupConfig)
	assert.Equal(t, false, client.watch)
	assert.Equal(t, true, client.useLocalConfIfKeyNotExist)
	assert.NotNil(t, client.apolloClient)
}

func TestApolloClient_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟对象
	localConfigure := NewMockConfigure(ctrl)
	mockClient := NewMockClient(ctrl)
	mockCache := NewMockCacheInterface(ctrl)

	// 设置模拟对象的行为
	mockClient.EXPECT().GetConfigCache("application").Return(mockCache).AnyTimes()
	mockCache.EXPECT().Get("test.key").Return("test-value", nil).AnyTimes()

	// 创建apolloClient实例
	client := &apolloClient{
		localConfigure:            localConfigure,
		apolloClient:              mockClient,
		namespace:                 "application",
		changeListener:            &changeListener{},
		watch:                     false,
		useLocalConfIfKeyNotExist: true,
	}

	// 测试从Apollo获取配置
	var value string
	err := client.Get("test.key", &value, "default-value")
	assert.Nil(t, err)
	assert.Equal(t, "test-value", value)

	// 测试获取配置失败时使用本地配置
	mockCache.EXPECT().Get("test.not-exist").Return(nil, errors.New("key not found")).AnyTimes()
	localConfigure.EXPECT().Get("test.not-exist", gomock.Any(), "default-value").Return(nil).Do(
		func(key string, v any, defaultVal string) {
			*(v.(*string)) = "local-value"
		},
	)

	var localValue string
	err = client.Get("test.not-exist", &localValue, "default-value")
	assert.Nil(t, err)
	assert.Equal(t, "local-value", localValue)

	// 测试不使用本地配置
	client.useLocalConfIfKeyNotExist = false
	var noLocalValue string
	err = client.Get("test.not-exist", &noLocalValue, "default-value")
	assert.Nil(t, err)
	assert.Equal(t, "", noLocalValue) // 应该是空值，因为Apollo没有找到且不使用本地配置
}

func TestApolloClient_Get_WithWatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟对象
	localConfigure := NewMockConfigure(ctrl)
	mockClient := NewMockClient(ctrl)
	mockCache := NewMockCacheInterface(ctrl)

	// 设置模拟对象的行为
	mockClient.EXPECT().GetConfigCache("application").Return(mockCache).AnyTimes()
	mockCache.EXPECT().Get("test.key").Return("test-value", nil).AnyTimes()

	// 创建changeListener
	listener := &changeListener{}
	listener.Init()

	// 创建apolloClient实例
	client := &apolloClient{
		localConfigure:            localConfigure,
		apolloClient:              mockClient,
		namespace:                 "application",
		changeListener:            listener,
		watch:                     true,
		useLocalConfIfKeyNotExist: true,
	}

	// 测试带监听的配置获取
	var value string
	err := client.Get("test.key", &value, "default-value")
	assert.Nil(t, err)
	assert.Equal(t, "test-value", value)

	// 验证监听器是否正确注册了key
	_, exists := listener.keyMap["test.key"]
	assert.True(t, exists)

	// 测试配置变更通知
	changes := make(map[string]*storage.ConfigChange)
	changes["test.key"] = &storage.ConfigChange{
		OldValue:   "test-value",
		NewValue:   "new-value",
		ChangeType: storage.MODIFIED,
	}

	changeEvent := &storage.ChangeEvent{
		Changes: changes,
	}

	// 触发配置变更
	listener.OnChange(changeEvent)

	// 验证值是否被更新
	assert.Equal(t, "new-value", value)
}

func TestSetValue(t *testing.T) {
	// 测试字符串值设置
	var strValue string
	err := setValue(&strValue, "test-string")
	assert.Nil(t, err)
	assert.Equal(t, "test-string", strValue)

	// 测试整数值设置
	var intValue int
	err = setValue(&intValue, "123")
	assert.Nil(t, err)
	assert.Equal(t, 123, intValue)

	// 测试布尔值设置
	var boolValue bool
	err = setValue(&boolValue, "true")
	assert.Nil(t, err)
	assert.Equal(t, true, boolValue)

	// 测试JSON对象设置
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var structValue TestStruct
	err = setValue(&structValue, map[string]any{"name": "test", "age": 30})
	assert.Nil(t, err)
	assert.Equal(t, "test", structValue.Name)
	assert.Equal(t, 30, structValue.Age)
}

func TestApolloClient_MultiNamespace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟对象
	localConfigure := NewMockConfigure(ctrl)
	mockClient := NewMockClient(ctrl)
	mockCache1 := NewMockCacheInterface(ctrl)
	mockCache2 := NewMockCacheInterface(ctrl)

	// 设置模拟对象的行为
	mockClient.EXPECT().GetConfigCache("application").Return(mockCache1).AnyTimes()
	mockClient.EXPECT().GetConfigCache("test.yml").Return(mockCache2).AnyTimes()

	// 第一个namespace没有找到key
	mockCache1.EXPECT().Get("test.key").Return(nil, errors.New("key not found")).AnyTimes()
	// 第二个namespace找到了key
	mockCache2.EXPECT().Get("test.key").Return("test-value-ns2", nil).AnyTimes()

	// 创建apolloClient实例
	client := &apolloClient{
		localConfigure:            localConfigure,
		apolloClient:              mockClient,
		namespace:                 "application,test.yml",
		changeListener:            &changeListener{},
		watch:                     false,
		useLocalConfIfKeyNotExist: true,
	}

	// 测试从第二个namespace获取配置
	var value string
	err := client.Get("test.key", &value, "default-value")
	assert.Nil(t, err)
	assert.Equal(t, "test-value-ns2", value)
}

func Test_Integration(t *testing.T) {
	// 这是一个集成测试，使用gone框架的测试机制
	// 仅在有Apollo服务器时运行
	t.Skip("Skip integration test that requires Apollo server")

	gone.
		NewApp(Load). // 加载Apollo模块
		Test(func(in struct {
			a string  `gone:"config,test.a"`
			b *string `gone:"config,test.b"`
			c string  `gone:"config,test.c"`
			d string  `gone:"config,test.d"`
			e string  `gone:"config,test.e"`
			x *string `gone:"config,test.x"`
		}) {
			// 验证配置值
			assert.Equal(t, "value-a", in.a)
			assert.Equal(t, "value-b", *in.b)
			assert.Equal(t, "value-c", in.c)
			assert.Equal(t, "value-d", in.d)
			assert.Equal(t, "value-e", in.e)
			assert.Equal(t, "value-x", *in.x)
		})
}
