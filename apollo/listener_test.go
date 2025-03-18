package apollo

import (
	"testing"

	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestChangeListener_Init(t *testing.T) {
	// 创建changeListener实例
	listener := &changeListener{}

	// 执行初始化
	listener.Init()

	// 验证keyMap已初始化
	assert.NotNil(t, listener.keyMap, "keyMap应该被初始化")
	assert.Empty(t, listener.keyMap, "keyMap应该为空")
}

func TestChangeListener_Put(t *testing.T) {
	// 创建changeListener实例
	listener := &changeListener{}
	listener.Init()

	// 测试数据
	key := "test.key"
	var value string

	// 执行Put操作
	listener.Put(key, &value)

	// 验证key-value是否正确存储
	assert.Contains(t, listener.keyMap, key, "keyMap应该包含指定的key")
	assert.Equal(t, &value, listener.keyMap[key], "keyMap中存储的值应该是正确的引用")
}

func TestChangeListener_OnChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟logger
	mockLogger := NewMockLogger(ctrl)

	// 创建changeListener实例
	listener := &changeListener{
		logger: mockLogger,
	}
	listener.Init()

	// 测试数据
	key := "test.key"
	var value string
	listener.Put(key, &value)

	// 创建配置变更事件
	changeEvent := &storage.ChangeEvent{
		Changes: map[string]*storage.ConfigChange{
			key: {
				OldValue:   "",
				NewValue:   "new-value",
				ChangeType: storage.MODIFIED,
			},
		},
	}

	// 执行OnChange
	listener.OnChange(changeEvent)

	// 验证值是否被更新
	assert.Equal(t, "new-value", value, "值应该被更新为new-value")
}

func TestChangeListener_OnChange_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟logger
	mockLogger := NewMockLogger(ctrl)
	// 期望调用Warnf方法
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

	// 创建changeListener实例
	listener := &changeListener{
		logger: mockLogger,
	}
	listener.Init()

	// 测试数据 - 使用int类型，但新值为字符串，应该导致错误
	key := "test.key"
	var value int
	listener.Put(key, &value)

	// 创建配置变更事件 - 新值为非数字字符串，无法转换为int
	changeEvent := &storage.ChangeEvent{
		Changes: map[string]*storage.ConfigChange{
			key: {
				OldValue:   "0",
				NewValue:   "not-a-number",
				ChangeType: storage.MODIFIED,
			},
		},
	}

	// 执行OnChange - 应该记录警告日志
	listener.OnChange(changeEvent)

	// 值不应该被更新
	assert.Equal(t, 0, value, "值不应该被更新")
}

func TestChangeListener_OnChange_NotModified(t *testing.T) {
	// 创建changeListener实例
	listener := &changeListener{}
	listener.Init()

	// 测试数据
	key := "test.key"
	var value string = "original"
	listener.Put(key, &value)

	// 创建配置变更事件 - 使用非MODIFIED类型
	changeEvent := &storage.ChangeEvent{
		Changes: map[string]*storage.ConfigChange{
			key: {
				OldValue:   "original",
				NewValue:   "new-value",
				ChangeType: storage.ADDED, // 非MODIFIED类型
			},
		},
	}

	// 执行OnChange
	listener.OnChange(changeEvent)

	// 验证值没有被更新
	assert.Equal(t, "original", value, "值不应该被更新")
}

func TestChangeListener_OnNewestChange(t *testing.T) {
	// 创建changeListener实例
	listener := &changeListener{}

	// 执行OnNewestChange - 这是一个空实现，不应该有任何影响
	listener.OnNewestChange(nil)

	// 没有断言，因为这个方法是空实现
	// 这个测试只是为了覆盖率
}
