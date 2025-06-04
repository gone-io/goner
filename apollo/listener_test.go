package apollo

import (
	"github.com/gone-io/gone/v2"
	"testing"

	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct {
	gone.Logger
	warnings []string
}

func (m *mockLogger) Warnf(format string, args ...interface{}) {
	m.warnings = append(m.warnings, format)
}

func TestChangeListener_OnChange(t *testing.T) {
	tests := []struct {
		name           string
		namespace      string
		changes        map[string]*storage.ConfigChange
		initialValue   map[string]interface{}
		expectedValue  map[string]interface{}
		expectedErrors bool
		watch          gone.ConfWatchFunc
	}{
		{
			name:      "add new config",
			namespace: "application",
			changes: map[string]*storage.ConfigChange{
				"test.key": {
					OldValue:   "",
					NewValue:   "new value",
					ChangeType: storage.ADDED,
				},
			},
			initialValue:   map[string]interface{}{},
			expectedValue:  map[string]interface{}{"test.key": "new value"},
			expectedErrors: false,
			watch: func(oldVal, newVal any) {
				assert.Equal(t, nil, oldVal)
				assert.Equal(t, "new value", newVal)
			},
		},
		{
			name:      "modify existing config",
			namespace: "application",
			changes: map[string]*storage.ConfigChange{
				"test.key": {
					OldValue:   "old value",
					NewValue:   "updated value",
					ChangeType: storage.MODIFIED,
				},
			},
			initialValue: map[string]interface{}{
				"test": map[string]any{
					"key": "old value",
				},
			},
			expectedValue:  map[string]interface{}{"test.key": "updated value"},
			expectedErrors: false,
			watch: func(oldVal, newVal any) {
				assert.Equal(t, "old value", oldVal)
				assert.Equal(t, "updated value", newVal)
			},
		},
		{
			name:      "delete config",
			namespace: "application",
			changes: map[string]*storage.ConfigChange{
				"test.key": {
					OldValue:   "old value",
					NewValue:   "",
					ChangeType: storage.DELETED,
				},
			},
			initialValue: map[string]interface{}{
				"test": map[string]any{
					"key": "old value",
				},
			},
			expectedValue:  map[string]interface{}{"test.key": nil},
			expectedErrors: false,
			watch: func(oldVal, newVal any) {
				assert.Equal(t, "old value", oldVal)
				assert.Equal(t, nil, newVal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// 初始化监听器
			listener := &changeListener{}
			listener.Init()

			// 设置mock logger
			mockLogger := &mockLogger{}
			listener.logger = mockLogger

			// 初始化viper实例
			v := viper.New()
			_ = v.MergeConfigMap(tt.initialValue)
			listener.viper = v

			// 创建namespace对应的viper
			namespaceViper := viper.New()
			listener.vipersMap = map[string]*viper.Viper{
				tt.namespace: namespaceViper,
			}
			listener.vipers = []*viper.Viper{namespaceViper}

			event := &storage.ChangeEvent{
				Changes: tt.changes,
			}
			event.Namespace = tt.namespace

			var str = "test"
			listener.Put("test.key", &str)
			listener.Watch("test.key", tt.watch)

			// 触发配置变更
			listener.OnChange(event)

			// 验证结果
			for k, expectedVal := range tt.expectedValue {
				actualVal := listener.viper.Get(k)
				assert.Equal(t, expectedVal, actualVal)
			}

			// 验证错误日志
			if tt.expectedErrors {
				assert.NotEmpty(t, mockLogger.warnings)
			} else {
				assert.Empty(t, mockLogger.warnings)
			}
		})
	}
}

func TestChangeListener_Put(t *testing.T) {
	listener := &changeListener{}
	listener.Init()

	// 测试添加监听键值
	key := "test.key"
	value := "test value"
	listener.Put(key, value)

	// 验证键值是否正确添加到keyMap中
	values, exists := listener.keyMap[key]
	assert.True(t, exists)
	assert.Equal(t, 1, len(values))
	assert.Equal(t, value, values[0])
}
