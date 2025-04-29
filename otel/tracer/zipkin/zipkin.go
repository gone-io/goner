package zipkin

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/otel/tracer"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/trace"
)

func Provide(_ string, i struct {
	url     string            `gone:"config,otel.tracer.zipkin.url"`
	headers map[string]string `gone:"config,otel.tracer.zipkin.headers"`
}) (trace.SpanExporter, error) {
	exporter, err := zipkin.New(i.url, zipkin.WithHeaders(i.headers))
	if err != nil {
		return nil, gone.ToErrorWithMsg(err, "can not create zipkin exporter")
	}
	return exporter, nil
}

// Load for openTelemetry zipkin trace.SpanExporter
func Load(loader gone.Loader) error {
	loader.MustLoad(gone.WrapFunctionProvider(Provide))
	return tracer.Register(loader)
}
