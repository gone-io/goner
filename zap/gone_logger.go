package gone_zap

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newGoneLogger(provider *zapLoggerProvider) *goneLogger {
	s := goneLogger{
		provider: provider,
	}
	err := s.Init()
	g.PanicIfErr(err)
	return &s
}

type goneLogger struct {
	gone.Flag
	*zap.SugaredLogger
	provider *zapLoggerProvider `gone:"*"`
}

func (l *goneLogger) GonerName() string {
	return "gone-logger"
}

func (l *goneLogger) Init() error {
	logger, err := l.provider.Provide("")
	if err != nil {
		return gone.ToError(err)
	}
	l.SugaredLogger = logger.Sugar()
	return nil
}

func (l *goneLogger) GetLevel() gone.LoggerLevel {
	switch l.SugaredLogger.Level() {
	case zap.DebugLevel:
		return gone.DebugLevel
	case zap.InfoLevel:
		return gone.InfoLevel
	case zap.WarnLevel:
		return gone.WarnLevel
	case zap.ErrorLevel:
		return gone.ErrorLevel
	default:
		return gone.ErrorLevel
	}
}

func (l *goneLogger) SetLevel(level gone.LoggerLevel) {
	var zapLevel zapcore.Level
	switch level {
	case gone.DebugLevel:
		zapLevel = zap.DebugLevel
	case gone.InfoLevel:
		zapLevel = zap.InfoLevel
	case gone.WarnLevel:
		zapLevel = zap.WarnLevel
	case gone.ErrorLevel:
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}
	l.provider.SetLevel(zapLevel)
}
