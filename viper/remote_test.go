package viper

import (
	"github.com/gone-io/gone/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 测试RemoteConfigure的构造函数
func TestNewRemoteConfigure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockViper := NewMockKeyGetter(ctrl)
	mockLocalConfigure := NewMockConfigure(ctrl)
	mockWatcher := NewMockWatcherKeeper(ctrl)

	// 创建RemoteConfigure实例
	remoteCfg := NewRemoteConfigure(mockViper, mockLocalConfigure, true, mockWatcher)

	// 验证结果
	assert.NotNil(t, remoteCfg)
	assert.Equal(t, mockViper, remoteCfg.Viper)
	assert.Equal(t, mockLocalConfigure, remoteCfg.local)
	assert.Equal(t, mockWatcher, remoteCfg.watcher)
	assert.True(t, remoteCfg.useLocalConfIfKeyNotExist)
}

// 测试Get方法 - 从远程配置获取值
func TestRemoteConfigure_Get_FromRemote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockViper := NewMockKeyGetter(ctrl)
	mockLocalConfigure := NewMockConfigure(ctrl)
	mockWatcher := NewMockWatcherKeeper(ctrl)

	// 创建RemoteConfigure实例
	remoteCfg := NewRemoteConfigure(mockViper, mockLocalConfigure, true, mockWatcher)

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
	mockWatcher.EXPECT().Put(testKey, &result)

	// 执行测试
	err := remoteCfg.Get(testKey, &result, "default-value")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, testValue, result)
}

// 测试Get方法 - 当远程配置不存在时回退到本地配置
func TestRemoteConfigure_Get_FallbackToLocal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockViper := NewMockKeyGetter(ctrl)
	mockLocalConfigure := NewMockConfigure(ctrl)
	mockWatcher := NewMockWatcherKeeper(ctrl)

	// 创建RemoteConfigure实例
	remoteCfg := NewRemoteConfigure(mockViper, mockLocalConfigure, true, mockWatcher)

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
	mockWatcher.EXPECT().Put(testKey, &result)

	// 执行测试
	err := remoteCfg.Get(testKey, &result, "default-value")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, testValue, result)
}

// 测试Get方法 - 当远程配置解析失败时回退到本地配置
func TestRemoteConfigure_Get_FallbackToLocalWhenUnmarshalFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockViper := NewMockKeyGetter(ctrl)
	mockLocalConfigure := NewMockConfigure(ctrl)
	mockWatcher := NewMockWatcherKeeper(ctrl)

	// 创建RemoteConfigure实例
	remoteCfg := NewRemoteConfigure(mockViper, mockLocalConfigure, true, mockWatcher)

	// 设置Mock行为
	testKey := "test.key"
	remoteValue := "invalid-value"
	localValue := "local-value"
	var result string

	// 远程配置返回值但解析失败
	mockViper.EXPECT().Get(testKey).Return(remoteValue)
	mockViper.EXPECT().UnmarshalKey(testKey, gomock.Any()).Return(gone.ToError("unmarshal error"))

	// 本地配置返回值
	mockLocalConfigure.EXPECT().Get(testKey, gomock.Any(), "default-value").DoAndReturn(
		func(key string, value any, defaultVal string) error {
			p := value.(*string)
			*p = localValue
			return nil
		})
	mockWatcher.EXPECT().Put(testKey, &result)

	// 执行测试
	err := remoteCfg.Get(testKey, &result, "default-value")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, localValue, result)
}

// 测试Get方法 - 当viper为nil时直接使用本地配置
func TestRemoteConfigure_Get_NilViper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockLocalConfigure := NewMockConfigure(ctrl)
	mockWatcher := NewMockWatcherKeeper(ctrl)

	// 创建RemoteConfigure实例，viper为nil
	remoteCfg := NewRemoteConfigure(nil, mockLocalConfigure, true, mockWatcher)

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
	mockWatcher.EXPECT().Put(testKey, &result)

	// 执行测试
	err := remoteCfg.Get(testKey, &result, "default-value")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, testValue, result)
}
