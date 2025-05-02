package gone_zap

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

//type sugarProvider struct {
//	gone.Flag
//
//	zapLogger *zap.Logger `gone:"*"`
//	wrapped   *wrappedLogger
//}
//
//func (s *sugarProvider) Provide(tagConf string) (Logger, error) {
//	if s.wrapped == nil {
//		s.wrapped = &wrappedLogger{Logger: s.zapLogger}
//	}
//
//	_, keys := gone.TagStringParse(tagConf)
//	if len(keys) > 0 {
//		if keys[0] != "" {
//			return s.wrapped.Named(keys[0]), nil
//		}
//	}
//	return s.wrapped, nil
//}

// GetTraceIdFromCtx which get traceId from context and convert to zap.Field
// Examples:
//
//	logger := zap.NewNop()
//	logger.With(CtxTraceId(context.Background())).Info("info")
//	logger.Info("info", CtxTraceId(context.Background()))
func GetTraceIdFromCtx(ctx context.Context) zap.Field {
	span := trace.SpanFromContext(ctx)
	spanContext := span.SpanContext()
	var traceId string
	if spanContext.IsValid() {
		traceId = spanContext.TraceID().String()
	}
	return zap.String(traceIdKey, traceId)
}

var _ g.CtxLogger = (*ctxLogger)(nil)

type ctxLogger struct {
	gone.Flag
	provider *zapLoggerProvider `gone:"*"`
	tracer   g.Tracer           `gone:"*" option:"allowNil"`
}

func (c ctxLogger) Ctx(ctx context.Context) gone.Logger {
	logger := newGoneLogger(c.provider)
	span := trace.SpanFromContext(ctx)
	spanContext := span.SpanContext()
	var traceId string
	if spanContext.IsValid() {
		traceId = spanContext.TraceID().String()
	} else if c.tracer != nil {
		traceId = c.tracer.GetTraceId()
	}
	if traceId != "" {
		logger.SugaredLogger = logger.With(zap.String(traceIdKey, traceId))
	} else {
		logger.Warnf("traceId is empty; maybe, you should execute `gonectl install goner/trace` 、 `gonectl install goner/otel/tracer/http`")
	}
	return logger
}
