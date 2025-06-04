package remote

import (
	"testing"

	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 测试remoteConfigure的Init方法 - AddRemoteProvider失败
func TestRemoteConfigure_Init_AddRemoteProviderFail(t *testing.T) {
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
			*p = false
			return nil
		})

	mockLocalConfigure.EXPECT().Get("viper.remote.watchDuration", gomock.Any(), "5s").Return(nil)
	mockLocalConfigure.EXPECT().Get("viper.remote.useLocalConfIfKeyNotExist", gomock.Any(), "true").Return(nil)

	// 设置Provider的预期行为 - AddRemoteProvider失败
	mockViper.EXPECT().SetConfigType(providers[0].ConfigType)
	mockViper.EXPECT().AddRemoteProvider(providers[0].Provider, providers[0].Endpoint, providers[0].Path).Return(gone.ToError("add remote provider failed"))

	// 创建测试对象
	remoteCfg := &remoteConfigure{}

	// 保存原始函数并在测试结束后恢复
	origNewViper := newViper
	defer func() { newViper = origNewViper }()

	// 替换为Mock函数
	newViper = func() ViperInterface {
		return mockViper
	}
	newRemoteViper = newViper

	// 执行测试
	err := remoteCfg.init(mockLocalConfigure, mockViper)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "add remote provider failed")
}

// 测试remoteConfigure的Init方法 - AddSecureRemoteProvider失败
func TestRemoteConfigure_Init_AddSecureRemoteProviderFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockLocalConfigure := NewMockConfigure(ctrl)
	mockViper := NewMockViperInterface(ctrl)

	// 设置测试数据
	providers := []Provider{
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
			*p = false
			return nil
		})

	mockLocalConfigure.EXPECT().Get("viper.remote.watchDuration", gomock.Any(), "5s").Return(nil)
	mockLocalConfigure.EXPECT().Get("viper.remote.useLocalConfIfKeyNotExist", gomock.Any(), "true").Return(nil)

	// 设置Provider的预期行为 - AddSecureRemoteProvider失败
	mockViper.EXPECT().SetConfigType(providers[0].ConfigType)
	mockViper.EXPECT().AddSecureRemoteProvider(providers[0].Provider, providers[0].Endpoint, providers[0].Path, providers[0].Keyring).Return(gone.ToError("add secure remote provider failed"))

	// 创建测试对象
	remoteCfg := &remoteConfigure{}

	// 保存原始函数并在测试结束后恢复
	origNewViper := newViper
	defer func() { newViper = origNewViper }()

	// 替换为Mock函数
	newViper = func() ViperInterface {
		return mockViper
	}
	newRemoteViper = newViper

	// 执行测试
	err := remoteCfg.init(mockLocalConfigure, mockViper)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "add secure remote provider failed")
}

// 测试remoteConfigure的Init方法 - ReadRemoteConfig失败
func TestRemoteConfigure_Init_ReadRemoteConfigFail(t *testing.T) {
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
			*p = false
			return nil
		})

	mockLocalConfigure.EXPECT().Get("viper.remote.watchDuration", gomock.Any(), "5s").Return(nil)
	mockLocalConfigure.EXPECT().Get("viper.remote.useLocalConfIfKeyNotExist", gomock.Any(), "true").Return(nil)

	// 设置Provider的预期行为 - ReadRemoteConfig失败
	mockViper.EXPECT().SetConfigType(providers[0].ConfigType)
	mockViper.EXPECT().AddRemoteProvider(providers[0].Provider, providers[0].Endpoint, providers[0].Path).Return(nil)
	mockViper.EXPECT().ReadRemoteConfig().Return(gone.ToError("read remote config failed"))

	// 创建测试对象
	remoteCfg := &remoteConfigure{}

	// 保存原始函数并在测试结束后恢复
	origNewViper := newViper
	defer func() { newViper = origNewViper }()

	// 替换为Mock函数
	newViper = func() ViperInterface {
		return mockViper
	}
	newRemoteViper = newViper

	// 执行测试
	err := remoteCfg.init(mockLocalConfigure, mockViper)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read remote config failed")
}

// 测试remoteConfigure的Init方法 - MergeConfigMap失败
func TestRemoteConfigure_Init_MergeConfigMapFail(t *testing.T) {
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
			*p = false
			return nil
		})

	mockLocalConfigure.EXPECT().Get("viper.remote.watchDuration", gomock.Any(), "5s").Return(nil)
	mockLocalConfigure.EXPECT().Get("viper.remote.useLocalConfIfKeyNotExist", gomock.Any(), "true").Return(nil)

	// 设置Provider的预期行为 - MergeConfigMap失败
	mockViper.EXPECT().SetConfigType(providers[0].ConfigType)
	mockViper.EXPECT().AddRemoteProvider(providers[0].Provider, providers[0].Endpoint, providers[0].Path).Return(nil)
	mockViper.EXPECT().ReadRemoteConfig().Return(nil)
	mockViper.EXPECT().AllSettings().Return(map[string]any{"key": "value"})
	mockViper.EXPECT().MergeConfigMap(gomock.Any()).Return(gone.ToError("merge config map failed"))

	// 创建测试对象
	remoteCfg := &remoteConfigure{}

	// 保存原始函数并在测试结束后恢复
	origNewViper := newViper
	defer func() { newViper = origNewViper }()

	// 替换为Mock函数
	newViper = func() ViperInterface {
		return mockViper
	}
	newRemoteViper = newViper

	// 执行测试
	err := remoteCfg.init(mockLocalConfigure, mockViper)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "merge config map failed")
}
