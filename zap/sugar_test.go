package gone_zap

import (
	"fmt"
	"os"
	"testing"

	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type tracer struct {
	gone.Flag
	g.Tracer
}

func (t *tracer) SetTraceId(traceId string, fn func()) {
	fn()
}

// GetTraceId Get the traceId of the current goroutine
func (t *tracer) GetTraceId() string {
	return "trace-id"
}

// Go Start a new goroutine instead of `go func`, which can pass the traceId to the new goroutine.
func (t *tracer) Go(fn func()) {
	go fn()
}

func TestSugar_GetLevel(t *testing.T) {
	tests := []struct {
		name      string
		zapLevel  zapcore.Level
		wantLevel gone.LoggerLevel
	}{
		{"debug level", zap.DebugLevel, gone.DebugLevel},
		{"info level", zap.InfoLevel, gone.InfoLevel},
		{"warn level", zap.WarnLevel, gone.WarnLevel},
		{"error level", zap.ErrorLevel, gone.ErrorLevel},
		{"fatal level", zap.FatalLevel, gone.ErrorLevel},
		{"panic level", zap.PanicLevel, gone.ErrorLevel},
		{"debug-1 level", zapcore.Level(-1), gone.DebugLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建一个测试用的provider
			provider := &zapLoggerProvider{
				atomicLevel: newAtomicLevel(tt.zapLevel),
				output:      "stdout",
				beforeStop: func(process gone.Process) {

				},
			}
			provider.SetLevel(tt.zapLevel)
			_ = provider.Init()

			// 创建logger
			logger, _ := provider.Provide("")
			s := &sugar{
				SugaredLogger: logger.Sugar(),
				provider:      provider,
			}

			// 测试GetLevel
			gotLevel := s.GetLevel()
			assert.Equal(t, tt.wantLevel, gotLevel)
		})
	}
}

func TestSugar_SetLevel(t *testing.T) {
	tests := []struct {
		name       string
		inputLevel gone.LoggerLevel
		wantLevel  zapcore.Level
	}{
		{"debug level", gone.DebugLevel, zap.DebugLevel},
		{"info level", gone.InfoLevel, zap.InfoLevel},
		{"warn level", gone.WarnLevel, zap.WarnLevel},
		{"error level", gone.ErrorLevel, zap.ErrorLevel},
		{"invalid level", gone.LoggerLevel(99), zap.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建一个测试用的provider
			provider := &zapLoggerProvider{
				atomicLevel: newAtomicLevel(zap.ErrorLevel),
				output:      "stdout",
				beforeStop: func(process gone.Process) {

				},
			}

			_ = provider.Init()

			// 创建logger
			logger, _ := provider.Provide("")
			s := &sugar{
				SugaredLogger: logger.Sugar(),
				provider:      provider,
			}

			// 设置日志级别
			s.SetLevel(tt.inputLevel)

			// 验证provider中的日志级别是否正确设置
			assert.Equal(t, tt.wantLevel, provider.atomicLevel.Level())
		})
	}
}

func TestNewSugar(t *testing.T) {
	_ = os.Setenv("GONE_LOG_LEVEL", "debug")
	gone.
		NewApp(Load).
		Load(&tracer{}).
		Test(func(log gone.Logger, tracer g.Tracer, in struct {
			level string `gone:"config,log.level"`
		}) {
			log.Infof("level:%s", in.level)
			if in.level != "debug" {
				t.Fatal("log level error")
			}
			tracer.SetTraceId("", func() {
				log.Debugf("debug log")
				log.Infof("info log")
				log.Warnf("warn log")
				log.Errorf("error log")
			})
		})
}

func Test_sugar_Init(t *testing.T) {
	type fields struct {
		Flag          gone.Flag
		SugaredLogger *zap.SugaredLogger
		provider      *zapLoggerProvider
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "test error",
			fields: fields{
				Flag:          gone.Flag{},
				SugaredLogger: zap.NewNop().Sugar(),
				provider: &zapLoggerProvider{
					atomicLevel: newAtomicLevel(zap.ErrorLevel),
					output:      "stdout",
					beforeStop: func(process gone.Process) {
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &sugar{
				Flag:          tt.fields.Flag,
				SugaredLogger: tt.fields.SugaredLogger,
				provider:      tt.fields.provider,
			}
			tt.wantErr(t, l.Init(), fmt.Sprintf("Init()"))
		})
	}
}
