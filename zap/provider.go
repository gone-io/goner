package gone_zap

import (
	"errors"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/log/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/fs"
	"strings"
)

type zapLoggerProvider struct {
	gone.Flag

	disableStacktrace   bool   `gone:"config,log.disable-stacktrace,default=false"`
	stackTraceLevel     string `gone:"config,log.stacktrace-level,default=error"`
	reportCaller        bool   `gone:"config,log.report-caller,default=true"`
	encoder             string `gone:"config,log.encoder,default=console"`
	output              string `gone:"config,log.output,default=stdout"`
	errOutput           string `gone:"config,log.error-output"`
	rotationOutput      string `gone:"config,log.rotation.output"`
	rotationErrorOutput string `gone:"config,log.rotation.error-output"`
	rotationMaxSize     int    `gone:"config,log.rotation.max-size,default=100"`
	rotationMaxFiles    int    `gone:"config,log.rotation.max-files,default=10"`
	rotationMaxAge      int    `gone:"config,log.rotation.max-age,default=30"`
	rotationLocalTime   bool   `gone:"config,log.rotation.local-time,default=true"`
	rotationCompress    bool   `gone:"config,log.rotation.compress,default=false"`
	otelEnable          bool   `gone:"config,log.otel.enable=false"`
	otelOnly            bool   `gone:"config,log.otel.only=true"`
	otelLogName         string `gone:"config,log.otel.log-name=zap"`

	useEncoder      zapcore.Encoder   `gone:"*" option:"allowNil"`
	tracer          g.Tracer          `gone:"*" option:"allowNil"`
	isOtelLogLoaded g.IsOtelLogLoaded `gone:"*" option:"allowNil"`
	atomicLevel     *atomicLevel      `gone:"*"`
	beforeStop      gone.BeforeStop   `gone:"*"`

	zapLogger *zap.Logger
}

func (s *zapLoggerProvider) Provide(tagConf string) (*zap.Logger, error) {
	_, keys := gone.TagStringParse(tagConf)
	if len(keys) > 0 && keys[0] != "" {
		return s.zapLogger.Named(keys[0]), nil
	}
	if s.zapLogger == nil {
		return nil, gone.ToError("zapLogger is nil")
	}
	return s.zapLogger, nil
}

func (s *zapLoggerProvider) processError(err error) {
	var e *fs.PathError
	if errors.As(err, &e) && strings.HasPrefix(e.Path, "/dev/") {
		return
	}
	g.ErrorPrinter(gone.GetDefaultLogger(), err, "zapLogger.Sync")
}

func (s *zapLoggerProvider) Init() error {
	if s.zapLogger == nil {
		s.zapLogger = s.create()
		s.beforeStop(func() {
			s.processError(s.zapLogger.Sync())
		})
	}
	return nil
}

func (s *zapLoggerProvider) SetLevel(level zapcore.Level) {
	s.atomicLevel.SetLevel(level)
}

func (s *zapLoggerProvider) createFileCore() (core zapcore.Core) {
	outputs := strings.Split(s.output, ",")
	sink, closeOut, err := zap.Open(outputs...)
	g.PanicIfErr(gone.ToErrorWithMsg(err, "create zap output sink"))

	if s.rotationOutput != "" {
		rotationWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   s.rotationOutput,
			MaxSize:    s.rotationMaxSize, // megabytes
			MaxBackups: s.rotationMaxFiles,
			MaxAge:     s.rotationMaxAge, // days
			Compress:   s.rotationCompress,
		})

		sink = zap.CombineWriteSyncers(sink, rotationWriter)
	}

	errOutputs := strings.Split(s.errOutput, ",")
	var errSink zapcore.WriteSyncer
	if s.errOutput != "" && len(errOutputs) > 0 {
		errSink, _, err = zap.Open(errOutputs...)
		if err != nil {
			closeOut()
			g.PanicIfErr(gone.ToErrorWithMsg(err, "create zap error sink"))
		}
	}

	if s.rotationErrorOutput != "" {
		rotationWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   s.rotationErrorOutput,
			MaxSize:    s.rotationMaxSize, // megabytes
			MaxBackups: s.rotationMaxFiles,
			MaxAge:     s.rotationMaxAge, // days
			Compress:   s.rotationCompress,
		})
		if errSink == nil {
			errSink = rotationWriter
		} else {
			errSink = zap.CombineWriteSyncers(rotationWriter, errSink)
		}
	}

	if s.useEncoder == nil {
		if s.encoder == "console" {
			config := zap.NewDevelopmentEncoderConfig()
			config.ConsoleSeparator = "|"
			s.useEncoder = zapcore.NewConsoleEncoder(config)
		} else {
			s.useEncoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		}
		s.useEncoder = NewTraceEncoder(s.useEncoder, s.tracer)
	}

	core = zapcore.NewCore(
		s.useEncoder,
		sink,
		s.atomicLevel,
	)

	if errSink != nil {
		core = zapcore.NewTee(
			core,
			zapcore.NewCore(
				s.useEncoder,
				errSink,
				s.atomicLevel,
			),
		)
	}
	return core
}

func (s *zapLoggerProvider) createCore() (core zapcore.Core) {
	if bool(s.isOtelLogLoaded) && s.otelEnable {
		provider := global.GetLoggerProvider()
		core = otelzap.NewCore("otelLogName", otelzap.WithLoggerProvider(provider))

		if !s.otelOnly {
			core = zapcore.NewTee(
				core,
				s.createFileCore(),
			)
		}
	} else {
		core = s.createFileCore()
	}
	return core
}

func (s *zapLoggerProvider) create() *zap.Logger {
	var core = s.createCore()
	var opts []zap.Option
	if !s.disableStacktrace {
		opts = append(opts, zap.AddStacktrace(parseLevel(s.stackTraceLevel)))
	}
	if s.reportCaller {
		opts = append(opts, zap.AddCaller())
	}
	logger := zap.New(core, opts...)
	return logger
}
