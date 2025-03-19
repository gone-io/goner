package gone_zap

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func parseLevel(level string) zapcore.Level {
	switch level {
	default:
		return zap.InfoLevel
	case "debug", "trace":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	}
}

type atomicLevel struct {
	gone.Flag
	level *string `gone:"config,log.level,default=info"`
}

func (a atomicLevel) Enabled(level zapcore.Level) bool {
	zapLevel := parseLevel(*a.level)
	return zapLevel <= level
}

func (a atomicLevel) Level() zapcore.Level {
	fmt.Printf("a.level:%v\n", a.level)
	return parseLevel(*a.level)
}

func (a atomicLevel) SetLevel(level zapcore.Level) {
	if a.level == nil {
		a.level = new(string)
	}

	switch level {
	case zap.DebugLevel:
		*a.level = "debug"
	case zap.InfoLevel:
		*a.level = "info"
	case zap.WarnLevel:
		*a.level = "warn"
	case zap.ErrorLevel:
		*a.level = "error"
	case zap.PanicLevel:
		*a.level = "panic"
	case zap.FatalLevel:
		*a.level = "fatal"
	default:
		*a.level = "info"
	}
}
