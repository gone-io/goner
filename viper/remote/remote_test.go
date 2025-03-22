package remote

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 测试remoteConfigure的Init方法
func TestRemoteConfigure_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockLocalConfigure := NewMockConfigure(ctrl)
	mockViper := NewMockViperInterface(ctrl)

	// 设置测试数据
	providers := []Provider{
		{
			Provider:   "etcd",
			Endpoint:   "localhost:2379",
			Path:       "/config/myapp",
			ConfigType: "json",
		},
		{
			Provider:   "consul",
			Endpoint:   "localhost:8500",
			Path:       "/config/myapp",
			ConfigType: "yaml",
			Keyring:    "test-keyring",
		},
	}

	// 设置Mock行为
	mockLocalConfigure.EXPECT().Get("viper.remote.providers", gomock.Any(), "").DoAndReturn(
		func(key string, value any, defaultVal string) error {
			p := value.(*[]Provider)
			*p = providers
			return nil
		})

	mockLocalConfigure.EXPECT().Get("viper.remote.watch", gomock.Any(), "false").DoAndReturn(
		func(key string, value any, defaultVal string) error {
			p := value.(*bool)
			*p = true
			return nil
		})

	mockLocalConfigure.EXPECT().Get("viper.remote.watchDuration", gomock.Any(), "5s").DoAndReturn(
		func(key string, value any, defaultVal string) error {
			p := value.(*time.Duration)
			*p = 5 * time.Second
			return nil
		})

	mockLocalConfigure.EXPECT().Get("viper.remote.useLocalConfIfKeyNotExist", gomock.Any(), "true").DoAndReturn(
		func(key string, value any, defaultVal string) error {
			p := value.(*bool)
			*p = true
			return nil
		})

	// 设置第一个Provider的预期行为
	mockViper.EXPECT().SetConfigType(providers[0].ConfigType)
	mockViper.EXPECT().AddRemoteProvider(providers[0].Provider, providers[0].Endpoint, providers[0].Path).Return(nil)
	mockViper.EXPECT().ReadRemoteConfig().Return(nil)
	mockViper.EXPECT().AllSettings().Return(map[string]any{"key1": "value1"})
	mockViper.EXPECT().MergeConfigMap(gomock.Any()).Return(nil)

	// 设置第二个Provider的预期行为
	mockViper.EXPECT().SetConfigType(providers[1].ConfigType)
	mockViper.EXPECT().AddSecureRemoteProvider(providers[1].Provider, providers[1].Endpoint, providers[1].Path, providers[1].Keyring).Return(nil)
	mockViper.EXPECT().ReadRemoteConfig().Return(nil)
	mockViper.EXPECT().AllSettings().Return(map[string]any{"key2": "value2"})
	mockViper.EXPECT().MergeConfigMap(gomock.Any()).Return(nil)

	// 创建测试对象
	remoteCfg := &remoteConfigure{}

	// 保存原始函数并在测试结束后恢复
	origNewViper := newViper
	defer func() { newViper = origNewViper }()

	// 替换为Mock函数
	newViper = func() ViperInterface {
		return mockViper
	}

	// 执行测试
	err := remoteCfg.init(mockLocalConfigure, mockViper)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, true, remoteCfg.watch)
	assert.Equal(t, 5*time.Second, remoteCfg.watchDuration)
	assert.Equal(t, true, remoteCfg.useLocalConfIfKeyNotExist)
	assert.Equal(t, providers, remoteCfg.providers)
	assert.Len(t, remoteCfg.remoteVipers, 2)
}

// 测试remoteConfigure的Get方法 - 从远程配置获取值
func TestRemoteConfigure_Get_FromRemote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockLocalConfigure := NewMockConfigure(ctrl)
	mockViper := NewMockViperInterface(ctrl)

	// 创建测试对象
	remoteCfg := &remoteConfigure{
		localConfigure: mockLocalConfigure,
		viper:          mockViper,
		keyMap:         make(map[string][]any),
		watch:          true,
	}

	// 设置Mock行为
	testKey := "test.key"
	testValue := "test-value"
	var result string

	mockViper.EXPECT().Get(testKey).Return(testValue)
	mockViper.EXPECT().UnmarshalKey(testKey, gomock.Any()).DoAndReturn(
		func(key string, rawVal any, opts ...any) error {
			p := rawVal.(*string)
			*p = testValue
			return nil
		})

	// 执行测试
	err := remoteCfg.Get(testKey, &result, "default-value")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, testValue, result)
	assert.Contains(t, remoteCfg.keyMap, testKey)
	assert.Len(t, remoteCfg.keyMap[testKey], 1)
}

// 测试remoteConfigure的Get方法 - 当远程配置不存在时回退到本地配置
func TestRemoteConfigure_Get_FallbackToLocal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockLocalConfigure := NewMockConfigure(ctrl)
	mockViper := NewMockViperInterface(ctrl)

	// 创建测试对象
	remoteCfg := &remoteConfigure{
		localConfigure:            mockLocalConfigure,
		viper:                     mockViper,
		keyMap:                    make(map[string][]any),
		watch:                     true,
		useLocalConfIfKeyNotExist: true,
	}

	// 设置Mock行为
	testKey := "test.key"
	testValue := "local-value"
	var result string

	// 远程配置返回nil
	mockViper.EXPECT().Get(testKey).Return(nil)

	// 本地配置返回值
	mockLocalConfigure.EXPECT().Get(testKey, gomock.Any(), "default-value").DoAndReturn(
		func(key string, value any, defaultVal string) error {
			p := value.(*string)
			*p = testValue
			return nil
		})

	// 执行测试
	err := remoteCfg.Get(testKey, &result, "default-value")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, testValue, result)
	assert.Contains(t, remoteCfg.keyMap, testKey)
	assert.Len(t, remoteCfg.keyMap[testKey], 1)
}

// 测试remoteConfigure的Get方法 - 当viper为nil时直接使用本地配置
func TestRemoteConfigure_Get_NilViper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockLocalConfigure := NewMockConfigure(ctrl)

	// 创建测试对象
	remoteCfg := &remoteConfigure{
		localConfigure: mockLocalConfigure,
		viper:          nil,
		keyMap:         make(map[string][]any),
		watch:          true,
	}

	// 设置Mock行为
	testKey := "test.key"
	testValue := "local-value"
	var result string

	// 本地配置返回值
	mockLocalConfigure.EXPECT().Get(testKey, gomock.Any(), "default-value").DoAndReturn(
		func(key string, value any, defaultVal string) error {
			p := value.(*string)
			*p = testValue
			return nil
		})

	// 执行测试
	err := remoteCfg.Get(testKey, &result, "default-value")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, testValue, result)
	assert.Contains(t, remoteCfg.keyMap, testKey)
	assert.Len(t, remoteCfg.keyMap[testKey], 1)
}

// 测试compare函数
func TestCompare(t *testing.T) {
	// 测试相等的情况
	assert.True(t, compare("test", "test"))
	assert.True(t, compare(123, 123))
	assert.True(t, compare(true, true))
	assert.True(t, compare(nil, nil))
	assert.True(t, compare([]string{"a", "b"}, []string{"a", "b"}))
	assert.True(t, compare(map[string]int{"a": 1}, map[string]int{"a": 1}))

	// 测试不相等的情况
	assert.False(t, compare("test", "test2"))
	assert.False(t, compare(123, 456))
	assert.False(t, compare(true, false))
	assert.False(t, compare(nil, "not-nil"))
	assert.False(t, compare([]string{"a", "b"}, []string{"b", "a"}))
	assert.False(t, compare(map[string]int{"a": 1}, map[string]int{"a": 2}))
}
