package tracer

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	otelHelper "github.com/gone-io/goner/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	otelTrace "go.opentelemetry.io/otel/trace"
)

type providerHelper struct {
	gone.Flag
	resource       *resource.Resource        `gone:"*" option:"allowNil"`
	traceExporter  trace.SpanExporter        `gone:"*" option:"allowNil"`
	afterStop      gone.AfterStop            `gone:"*"`
	logger         gone.Logger               `gone:"*"`
	resourceGetter otelHelper.ResourceGetter `gone:"*"`
}

func (s *providerHelper) Init() (err error) {
	if s.traceExporter == nil {
		s.traceExporter, _ = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithoutTimestamps(),
		)
	}

	var options = []trace.TracerProviderOption{
		trace.WithBatcher(
			s.traceExporter,
		),
	}
	if s.resource == nil {
		if s.resource, err = s.resourceGetter.Get(); err != nil {
			return gone.ToErrorWithMsg(err, "can not get resource")
		}
	}
	options = append(options, trace.WithResource(s.resource))
	traceProvider := trace.NewTracerProvider(options...)
	otel.SetTracerProvider(traceProvider)
	s.afterStop(func() {
		ctx := context.Background()
		g.ErrorPrinter(s.logger, traceProvider.ForceFlush(ctx), "tracer provider ForceFlush")
		g.ErrorPrinter(s.logger, traceProvider.Shutdown(ctx), "tracer provider Shutdown")
	})
	return nil
}

func (s *providerHelper) Provide(tagConf string) (otelTrace.Tracer, error) {
	name, _ := gone.ParseGoneTag(tagConf)
	return otel.Tracer(name), nil
}

// Register for openTelemetry TracerProvider
func Register(loader gone.Loader) error {
	loader.MustLoad(gone.WrapFunctionProvider(func(tagConf string, param struct{}) (g.IsOtelTracerLoaded, error) {
		return true, nil
	}))

	loader.MustLoad(&providerHelper{})
	return otelHelper.HelpSetPropagator(loader)
}
