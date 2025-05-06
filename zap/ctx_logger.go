package gone_zap

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var _ g.CtxLogger = (*ctxLogger)(nil)

type ctxLogger struct {
	gone.Flag
	provider *zapLoggerProvider `gone:"*"`
	tracer   g.Tracer           `gone:"*" option:"allowNil"`
}

func (c ctxLogger) Ctx(ctx context.Context) gone.Logger {
	logger := newGoneLogger(c.provider)
	spanContext := trace.SpanContextFromContext(ctx)
	var traceId string
	if spanContext.HasTraceID() {
		logger.SugaredLogger = logger.With(zap.Any(contextKey, ctx))
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
