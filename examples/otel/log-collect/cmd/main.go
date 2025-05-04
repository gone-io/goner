package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"go.opentelemetry.io/otel"
)

func main() {
	gone.Run(func(logger gone.Logger, ctxLogger g.CtxLogger, gTracer g.Tracer, i struct {
		name string `gone:"config,otel.service.name"`
	}) {
		//logger.Infof("service name: %s", i.name)
		//logger.Infof("hello world")
		//logger.Debugf("debug info")
		//logger.Warnf("warn info")
		//logger.Errorf("error info")

		tracer := otel.Tracer("test-tracer")
		ctx, span := tracer.Start(context.Background(), "test-run")
		defer span.End()

		log := ctxLogger.Ctx(ctx)
		log.Infof("hello world with traceId")
		log.Warnf("debug info with traceId")

		//set traceId
		gTracer.SetTraceId(span.SpanContext().TraceID().String(), func() {
			doSomething(logger, log)
		})
	})
}

func doSomething(logger gone.Logger, log gone.Logger) {

	logger.Infof("get traceId by using trace.Trace")

	log.Infof("traceId setted by ctx logger")
}
