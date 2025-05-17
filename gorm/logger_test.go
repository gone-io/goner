package gorm

import (
	"context"
	"errors"
	mock "github.com/gone-io/gone"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm/logger"
)

// TestLogMode 测试LogMode方法
func TestLogMode(t *testing.T) {
	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLogger
	mockLogger := mock.NewMockLogger(ctrl)

	// 创建iLogger实例
	l := &iLogger{
		log:           mockLogger,
		LogLevel:      logger.Info,
		SlowThreshold: 200 * time.Millisecond,
	}

	// 测试LogMode方法
	newLoggerInterface := l.LogMode(logger.Error)
	newLogger, ok := newLoggerInterface.(*iLogger)

	// 验证结果
	assert.True(t, ok, "LogMode应该返回*iLogger类型")
	assert.Equal(t, logger.Error, newLogger.LogLevel, "LogLevel应该被更新为Error")
	assert.Equal(t, l.log, newLogger.log, "log字段应该保持不变")
	assert.Equal(t, l.SlowThreshold, newLogger.SlowThreshold, "SlowThreshold应该保持不变")
	// 确保原始logger没有被修改
	assert.Equal(t, logger.Info, l.LogLevel, "原始logger的LogLevel不应被修改")
}

// TestInfo 测试Info方法
func TestInfo(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		logLevel  logger.LogLevel
		expectLog bool
	}{
		{"LogLevel高于Info", logger.Info, true},
		{"LogLevel等于Info", logger.Info, true},
		{"LogLevel低于Info", logger.Silent, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建gomock控制器
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 创建MockLogger
			mockLogger := mock.NewMockLogger(ctrl)

			// 创建iLogger实例
			l := &iLogger{
				log:      mockLogger,
				LogLevel: tt.logLevel,
			}

			// 设置期望
			if tt.expectLog {
				mockLogger.EXPECT().
					Infof("test message", gomock.Any(), "param1", "param2").
					Times(1)
			}

			// 执行测试
			l.Info(ctx, "test message", "param1", "param2")
		})
	}
}

// TestWarn 测试Warn方法
func TestWarn(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		logLevel  logger.LogLevel
		expectLog bool
	}{
		{"LogLevel高于Warn", logger.Info, true},
		{"LogLevel等于Warn", logger.Warn, true},
		{"LogLevel低于Warn", logger.Error, false},
		{"LogLevel为Silent", logger.Silent, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建gomock控制器
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 创建MockLogger
			mockLogger := mock.NewMockLogger(ctrl)

			// 创建iLogger实例
			l := &iLogger{
				log:      mockLogger,
				LogLevel: tt.logLevel,
			}

			// 设置期望
			if tt.expectLog {
				mockLogger.EXPECT().
					Warnf("test message", gomock.Any(), "param1", "param2").
					Times(1)
			}

			// 执行测试
			l.Warn(ctx, "test message", "param1", "param2")
		})
	}
}

// TestError 测试Error方法
func TestError(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		logLevel  logger.LogLevel
		expectLog bool
	}{
		{"LogLevel等于Error", logger.Error, true},
		{"LogLevel低于Error", logger.Silent, false},
		{"LogLevel为Silent", logger.Silent, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建gomock控制器
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 创建MockLogger
			mockLogger := mock.NewMockLogger(ctrl)

			// 创建iLogger实例
			l := &iLogger{
				log:      mockLogger,
				LogLevel: tt.logLevel,
			}

			// 设置期望
			if tt.expectLog {
				mockLogger.EXPECT().
					Errorf("test message", gomock.Any(), "param1", "param2").
					Times(1)
			}

			// 执行测试
			l.Error(ctx, "test message", "param1", "param2")
		})
	}
}

// TestTrace_Silent 测试Trace方法在Silent日志级别下的行为
func TestTrace_Silent(t *testing.T) {
	ctx := context.Background()

	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLogger
	mockLogger := mock.NewMockLogger(ctrl)

	// 创建iLogger实例，设置为Silent级别
	l := &iLogger{
		log:      mockLogger,
		LogLevel: logger.Silent,
	}

	// 不期望任何日志调用
	// 执行测试
	l.Trace(ctx, time.Now(), func() (string, int64) {
		return "SELECT * FROM users", 10
	}, nil)
}

// TestTrace_Error 测试Trace方法在有错误时的行为
func TestTrace_Error(t *testing.T) {
	ctx := context.Background()
	testErr := errors.New("test error")

	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLogger
	mockLogger := mock.NewMockLogger(ctrl)

	// 创建iLogger实例
	l := &iLogger{
		log:      mockLogger,
		LogLevel: logger.Error,
	}

	// 设置期望 - 应该记录错误
	mockLogger.EXPECT().
		Debugf(gomock.Any(), testErr, gomock.Any(), gomock.Any(), "SELECT * FROM users").
		Times(1)

	// 执行测试
	l.Trace(ctx, time.Now(), func() (string, int64) {
		return "SELECT * FROM users", 10
	}, testErr)

	// 测试行数为-1的情况
	mockLogger.EXPECT().
		Debugf(gomock.Any(), testErr, gomock.Any(), "-", "SELECT * FROM users").
		Times(1)

	// 执行测试
	l.Trace(ctx, time.Now(), func() (string, int64) {
		return "SELECT * FROM users", -1
	}, testErr)
}

// TestTrace_SlowSQL 测试Trace方法对慢查询的处理
func TestTrace_SlowSQL(t *testing.T) {
	ctx := context.Background()

	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLogger
	mockLogger := mock.NewMockLogger(ctrl)

	// 创建iLogger实例，设置慢查询阈值为100ms
	l := &iLogger{
		log:           mockLogger,
		LogLevel:      logger.Warn,
		SlowThreshold: 100 * time.Millisecond,
	}

	// 设置期望 - 应该记录慢查询
	mockLogger.EXPECT().
		Debugf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), "SELECT * FROM users").
		Times(1)

	// 执行测试 - 使用sleep模拟慢查询
	beginTime := time.Now().Add(-200 * time.Millisecond) // 模拟查询已经运行了200ms
	l.Trace(ctx, beginTime, func() (string, int64) {
		return "SELECT * FROM users", 10
	}, nil)

	// 测试行数为-1的情况
	mockLogger.EXPECT().
		Debugf(gomock.Any(), gomock.Any(), gomock.Any(), "-", "SELECT * FROM users").
		Times(1)

	// 执行测试
	l.Trace(ctx, beginTime, func() (string, int64) {
		return "SELECT * FROM users", -1
	}, nil)
}

// TestTrace_Info 测试Trace方法在Info级别的行为
func TestTrace_Info(t *testing.T) {
	ctx := context.Background()

	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLogger
	mockLogger := mock.NewMockLogger(ctrl)

	// 创建iLogger实例
	l := &iLogger{
		log:      mockLogger,
		LogLevel: logger.Info,
	}

	// 设置期望 - 应该记录SQL
	mockLogger.EXPECT().
		Debugf(gomock.Any(), gomock.Any(), gomock.Any(), "SELECT * FROM users").
		Times(1)

	// 执行测试
	l.Trace(ctx, time.Now(), func() (string, int64) {
		return "SELECT * FROM users", 10
	}, nil)

	// 测试行数为-1的情况
	mockLogger.EXPECT().
		Debugf(gomock.Any(), gomock.Any(), "-", "SELECT * FROM users").
		Times(1)

	// 执行测试
	l.Trace(ctx, time.Now(), func() (string, int64) {
		return "SELECT * FROM users", -1
	}, nil)
}

// TestTrace_RecordNotFoundError 测试Trace方法对RecordNotFound错误的处理
func TestTrace_RecordNotFoundError(t *testing.T) {
	ctx := context.Background()

	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLogger
	mockLogger := mock.NewMockLogger(ctrl)

	// 创建iLogger实例
	l := &iLogger{
		log:      mockLogger,
		LogLevel: logger.Error,
	}

	// 不期望任何日志调用，因为RecordNotFound错误应该被忽略
	// 执行测试
	l.Trace(ctx, time.Now(), func() (string, int64) {
		return "SELECT * FROM users WHERE id = 1", 0
	}, logger.ErrRecordNotFound)
}
