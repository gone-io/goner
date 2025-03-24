package remote

import (
	"testing"

	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 测试watcher的Init方法
func TestWatcher_Init(t *testing.T) {
	w := &watcher{}
	w.Init()

	// 验证初始化后的状态
	assert.NotNil(t, w.remoteVipers)
	assert.NotNil(t, w.keyMap)
	assert.Len(t, w.remoteVipers, 0)
	assert.Len(t, w.keyMap, 0)
}

// 测试watcher的SetViper方法
func TestWatcher_SetViper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockViper := NewMockViperInterface(ctrl)
	w := &watcher{}
	w.Init()
	w.SetViper(mockViper)

	// 验证viper已被设置
	assert.Equal(t, mockViper, w.viper)
}

// 测试watcher的Put方法
func TestWatcher_Put(t *testing.T) {
	w := &watcher{}
	w.Init()

	// 测试添加单个值
	testKey := "test.key"
	testValue := "test-value"
	w.Put(testKey, testValue)

	// 验证值已被添加到keyMap
	assert.Contains(t, w.keyMap, testKey)
	assert.Len(t, w.keyMap[testKey], 1)
	assert.Equal(t, testValue, w.keyMap[testKey][0])

	// 测试添加多个值到同一个key
	testValue2 := "test-value-2"
	w.Put(testKey, testValue2)

	// 验证值已被添加到keyMap
	assert.Contains(t, w.keyMap, testKey)
	assert.Len(t, w.keyMap[testKey], 2)
	assert.Equal(t, testValue, w.keyMap[testKey][0])
	assert.Equal(t, testValue2, w.keyMap[testKey][1])
}

// 测试watcher的addViper方法
func TestWatcher_AddViper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockViper1 := NewMockViperInterface(ctrl)
	mockViper2 := NewMockViperInterface(ctrl)

	w := &watcher{}
	w.Init()

	// 添加第一个viper
	w.addViper(mockViper1)
	assert.Len(t, w.remoteVipers, 1)
	assert.Equal(t, mockViper1, w.remoteVipers[0])

	// 添加第二个viper
	w.addViper(mockViper2)
	assert.Len(t, w.remoteVipers, 2)
	assert.Equal(t, mockViper1, w.remoteVipers[0])
	assert.Equal(t, mockViper2, w.remoteVipers[1])
}

// 测试watcher的doWatch方法 - 成功更新配置
func TestWatcher_DoWatch_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockLogger := NewMockLogger(ctrl)
	mockViper := NewMockViperInterface(ctrl)
	mockRemoteViper := NewMockViperInterface(ctrl)
	mockAllViper := NewMockViperInterface(ctrl)

	// 保存原始函数并在测试结束后恢复
	origNewViper := newViper
	defer func() { newViper = origNewViper }()

	// 替换为Mock函数
	newViper = func() ViperInterface {
		return mockAllViper
	}

	// 创建测试对象
	w := &watcher{
		logger:       mockLogger,
		viper:        mockViper,
		remoteVipers: []ViperInterface{mockRemoteViper},
		keyMap:       make(map[string][]any),
	}

	// 设置测试数据
	testKey := "test.key"
	oldValue := "old-value"
	newValue := "new-value"
	testStruct := &struct{ Value string }{}
	w.Put(testKey, testStruct)

	// 设置Mock行为
	mockRemoteViper.EXPECT().ReadRemoteConfig().Return(nil)
	mockAllViper.EXPECT().MergeConfigMap(gomock.Any()).Return(nil)
	mockAllViper.EXPECT().AllSettings().Return(map[string]any{testKey: newValue})
	mockRemoteViper.EXPECT().AllSettings().Return(map[string]any{testKey: newValue})

	mockViper.EXPECT().Get(testKey).Return(oldValue)
	mockAllViper.EXPECT().Get(testKey).Return(newValue)
	mockAllViper.EXPECT().UnmarshalKey(testKey, gomock.Any()).DoAndReturn(
		func(key string, rawVal any, opts ...any) error {
			p := rawVal.(*struct{ Value string })
			p.Value = newValue
			return nil
		})
	mockViper.EXPECT().MergeConfigMap(gomock.Any()).Return(nil)

	// 执行测试
	w.doWatch()

	// 验证结果
	assert.Equal(t, newValue, testStruct.Value)
}

// 测试watcher的doWatch方法 - 读取远程配置失败
func TestWatcher_DoWatch_ReadRemoteConfigFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockLogger := NewMockLogger(ctrl)
	mockViper := NewMockViperInterface(ctrl)
	mockRemoteViper := NewMockViperInterface(ctrl)

	// 创建测试对象
	w := &watcher{
		logger:       mockLogger,
		viper:        mockViper,
		remoteVipers: []ViperInterface{mockRemoteViper},
		keyMap:       make(map[string][]any),
	}

	// 设置Mock行为 - 读取远程配置失败
	mockRemoteViper.EXPECT().ReadRemoteConfig().Return(gone.ToError("read remote config failed"))
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any())

	// 执行测试
	w.doWatch()

	// 验证结果 - 由于读取失败，不会有进一步的操作
}

// 测试watcher的doWatch方法 - 合并配置失败
func TestWatcher_DoWatch_MergeConfigFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockLogger := NewMockLogger(ctrl)
	mockViper := NewMockViperInterface(ctrl)
	mockRemoteViper := NewMockViperInterface(ctrl)
	mockAllViper := NewMockViperInterface(ctrl)

	// 保存原始函数并在测试结束后恢复
	origNewViper := newViper
	defer func() { newViper = origNewViper }()

	// 替换为Mock函数
	newViper = func() ViperInterface {
		return mockAllViper
	}

	// 创建测试对象
	w := &watcher{
		logger:       mockLogger,
		viper:        mockViper,
		remoteVipers: []ViperInterface{mockRemoteViper},
		keyMap:       make(map[string][]any),
	}

	// 设置Mock行为 - 读取成功但合并失败
	mockRemoteViper.EXPECT().ReadRemoteConfig().Return(nil)
	mockRemoteViper.EXPECT().AllSettings().Return(map[string]any{"key": "value"})
	mockAllViper.EXPECT().MergeConfigMap(gomock.Any()).Return(gone.ToError("merge config failed"))
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any())

	// 执行测试
	w.doWatch()

	// 验证结果 - 由于合并失败，不会有进一步的操作
}

// 测试watcher的doWatch方法 - 值未变化
func TestWatcher_DoWatch_NoValueChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockViper := NewMockViperInterface(ctrl)
	mockRemoteViper := NewMockViperInterface(ctrl)
	mockAllViper := NewMockViperInterface(ctrl)

	// 保存原始函数并在测试结束后恢复
	origNewViper := newViper
	defer func() { newViper = origNewViper }()

	// 替换为Mock函数
	newViper = func() ViperInterface {
		return mockAllViper
	}

	// 创建测试对象
	w := &watcher{
		viper:        mockViper,
		remoteVipers: []ViperInterface{mockRemoteViper},
		keyMap:       make(map[string][]any),
	}

	// 设置测试数据
	testKey := "test.key"
	testValue := "test-value"
	testStruct := &struct{ Value string }{}
	w.Put(testKey, testStruct)

	// 设置Mock行为 - 值未变化
	mockRemoteViper.EXPECT().ReadRemoteConfig().Return(nil)
	mockAllViper.EXPECT().MergeConfigMap(gomock.Any()).Return(nil)
	mockRemoteViper.EXPECT().AllSettings().Return(map[string]any{testKey: testValue})

	mockViper.EXPECT().Get(testKey).Return(testValue)
	mockAllViper.EXPECT().Get(testKey).Return(testValue)

	// 执行测试
	w.doWatch()

	// 验证结果 - 由于值未变化，不会调用UnmarshalKey和MergeConfigMap
}

// 测试watcher的doWatch方法 - UnmarshalKey失败
func TestWatcher_DoWatch_UnmarshalKeyFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockLogger := NewMockLogger(ctrl)
	mockViper := NewMockViperInterface(ctrl)
	mockRemoteViper := NewMockViperInterface(ctrl)
	mockAllViper := NewMockViperInterface(ctrl)

	// 保存原始函数并在测试结束后恢复
	origNewViper := newViper
	defer func() { newViper = origNewViper }()

	// 替换为Mock函数
	newViper = func() ViperInterface {
		return mockAllViper
	}

	// 创建测试对象
	w := &watcher{
		logger:       mockLogger,
		viper:        mockViper,
		remoteVipers: []ViperInterface{mockRemoteViper},
		keyMap:       make(map[string][]any),
	}

	// 设置测试数据
	testKey := "test.key"
	oldValue := "old-value"
	newValue := "new-value"
	testStruct := &struct{ Value string }{}
	w.Put(testKey, testStruct)

	// 设置Mock行为 - UnmarshalKey失败
	mockRemoteViper.EXPECT().ReadRemoteConfig().Return(nil)
	mockAllViper.EXPECT().MergeConfigMap(gomock.Any()).Return(nil)
	mockAllViper.EXPECT().AllSettings().Return(map[string]any{testKey: newValue})
	mockRemoteViper.EXPECT().AllSettings().Return(map[string]any{testKey: newValue})

	mockViper.EXPECT().Get(testKey).Return(oldValue)
	mockAllViper.EXPECT().Get(testKey).Return(newValue)
	mockAllViper.EXPECT().UnmarshalKey(testKey, gomock.Any()).Return(gone.ToError("unmarshal key failed"))
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any())
	mockViper.EXPECT().MergeConfigMap(gomock.Any()).Return(nil)

	// 执行测试
	w.doWatch()

	// 验证结果 - 尽管UnmarshalKey失败，但仍会调用MergeConfigMap
}

// 测试watcher的doWatch方法 - 最终MergeConfigMap失败
func TestWatcher_DoWatch_FinalMergeConfigFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建Mock对象
	mockLogger := NewMockLogger(ctrl)
	mockViper := NewMockViperInterface(ctrl)
	mockRemoteViper := NewMockViperInterface(ctrl)
	mockAllViper := NewMockViperInterface(ctrl)

	// 保存原始函数并在测试结束后恢复
	origNewViper := newViper
	defer func() { newViper = origNewViper }()

	// 替换为Mock函数
	newViper = func() ViperInterface {
		return mockAllViper
	}

	// 创建测试对象
	w := &watcher{
		logger:       mockLogger,
		viper:        mockViper,
		remoteVipers: []ViperInterface{mockRemoteViper},
		keyMap:       make(map[string][]any),
	}

	// 设置测试数据
	testKey := "test.key"
	oldValue := "old-value"
	newValue := "new-value"
	testStruct := &struct{ Value string }{}
	w.Put(testKey, testStruct)

	// 设置Mock行为 - 最终MergeConfigMap失败
	mockRemoteViper.EXPECT().ReadRemoteConfig().Return(nil)
	mockAllViper.EXPECT().MergeConfigMap(gomock.Any()).Return(nil)
	mockAllViper.EXPECT().AllSettings().Return(map[string]any{testKey: newValue})
	mockRemoteViper.EXPECT().AllSettings().Return(map[string]any{testKey: newValue})

	mockViper.EXPECT().Get(testKey).Return(oldValue)
	mockAllViper.EXPECT().Get(testKey).Return(newValue)
	mockAllViper.EXPECT().UnmarshalKey(testKey, gomock.Any()).DoAndReturn(
		func(key string, rawVal any, opts ...any) error {
			p := rawVal.(*struct{ Value string })
			p.Value = newValue
			return nil
		})
	mockViper.EXPECT().MergeConfigMap(gomock.Any()).Return(gone.ToError("final merge config failed"))
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any())

	// 执行测试
	w.doWatch()

	// 验证结果 - 尽管最终MergeConfigMap失败，但值已经被更新
	assert.Equal(t, newValue, testStruct.Value)
}
