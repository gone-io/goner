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
	traceExporter := s.traceExporter
	if traceExporter == nil {
		traceExporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithoutTimestamps(),
		)
		if err != nil {
			return gone.ToErrorWithMsg(err, "can not create stdout trace exporter")
		}
	}

	var options = []trace.TracerProviderOption{
		trace.WithBatcher(
			s.traceExporter,
		),
	}
	if s.resource != nil {
		options = append(options, trace.WithResource(s.resource))
	} else {
		res, err := s.resourceGetter.Get()
		if err != nil {
			return gone.ToErrorWithMsg(err, "can not get resource")
		}
		options = append(options, trace.WithResource(res))
	}
	traceProvider := trace.NewTracerProvider(options...)
	otel.SetTracerProvider(traceProvider)
	s.afterStop(func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			s.logger.Errorf("otel tracer provider helper: shutdown err: %v", err)
		}
	})
	return nil
}

func (s *providerHelper) Provide(_ string) (g.IsOtelTracerLoaded, error) {
	return true, nil
}

// Register for openTelemetry TracerProvider
func Register(loader gone.Loader) error {
	loader.MustLoad(&providerHelper{})
	return otelHelper.HelpSetPropagator(loader)
}
