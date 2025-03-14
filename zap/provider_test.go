package gone_zap

import (
	"os"
	"testing"

	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

// 测试zapLoggerProvider的基本功能
func TestZapLoggerProvider_Basic(t *testing.T) {
	// 设置环境变量以控制日志级别
	os.Setenv("GONE_LOG_LEVEL", "debug")

	gone.NewApp(Load).Test(func(provider *zapLoggerProvider) {
		// 测试初始化
		assert.NotNil(t, provider.zapLogger, "zapLogger should be initialized")

		// 测试日志级别
		assert.Equal(t, zapcore.DebugLevel, provider.atomicLevel.Level(), "Log level should be debug")

		// 测试Provide方法
		logger, err := provider.Provide("")
		assert.Nil(t, err, "Provide should not return error")
		assert.NotNil(t, logger, "Provided logger should not be nil")

		// 测试带名称的Provide方法
		namedLogger, err := provider.Provide("tag:testLogger")
		assert.Nil(t, err, "Provide with name should not return error")
		assert.NotNil(t, namedLogger, "Provided named logger should not be nil")

		// 测试SetLevel方法
		provider.SetLevel(zapcore.InfoLevel)
		assert.Equal(t, zapcore.InfoLevel, provider.atomicLevel.Level(), "Log level should be changed to info")
	})
}

// 测试不同日志级别的配置
func TestZapLoggerProvider_LogLevels(t *testing.T) {
	// 测试不同的日志级别
	levels := []string{"debug", "info", "warn", "error", "panic", "fatal"}
	expectedLevels := []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel, zapcore.PanicLevel, zapcore.FatalLevel}

	for i, level := range levels {
		os.Setenv("GONE_LOG_LEVEL", level)

		gone.NewApp(Load).Test(func(provider *zapLoggerProvider) {
			assert.Equal(t, expectedLevels[i], provider.atomicLevel.Level(), "Log level should be "+level)
		})
	}

	// 测试无效的日志级别，应该默认为info
	os.Setenv("GONE_LOG_LEVEL", "invalid")
	gone.NewApp(Load).Test(func(provider *zapLoggerProvider) {
		assert.Equal(t, zapcore.InfoLevel, provider.atomicLevel.Level(), "Invalid log level should default to info")
	})
}

// 测试sugarProvider的功能
func TestSugarProvider(t *testing.T) {
	gone.NewApp(Load).Test(func(provider *sugarProvider) {
		// 测试Provide方法
		logger, err := provider.Provide("")
		assert.Nil(t, err, "Provide should not return error")
		assert.NotNil(t, logger, "Provided logger should not be nil")

		// 测试带名称的Provide方法
		namedLogger, err := provider.Provide("tag:testLogger")
		assert.Nil(t, err, "Provide with name should not return error")
		assert.NotNil(t, namedLogger, "Provided named logger should not be nil")

		// 验证wrapped字段已初始化
		assert.NotNil(t, provider.wrapped, "wrapped logger should be initialized")
	})
}

// 测试parseLevel函数
func TestParseLevel(t *testing.T) {
	testCases := []struct {
		level    string
		expected zapcore.Level
	}{
		{"debug", zapcore.DebugLevel},
		{"trace", zapcore.DebugLevel},
		{"info", zapcore.InfoLevel},
		{"warn", zapcore.WarnLevel},
		{"error", zapcore.ErrorLevel},
		{"panic", zapcore.PanicLevel},
		{"fatal", zapcore.FatalLevel},
		{"invalid", zapcore.InfoLevel}, // 默认为info
		{"", zapcore.InfoLevel},        // 空字符串默认为info
	}

	for _, tc := range testCases {
		result := parseLevel(tc.level)
		assert.Equal(t, tc.expected, result, "parseLevel should return correct level for "+tc.level)
	}
}
