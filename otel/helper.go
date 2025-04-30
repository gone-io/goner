package otel

import (
	"github.com/gone-io/gone/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type ResourceGetter interface {
	Get() (*resource.Resource, error)
}

type helper struct {
	gone.Flag
	propagators []propagation.TextMapPropagator `gone:"*"`
	serviceName string                          `gone:"config,otel.service.name"`
}

func (s *helper) Init() {
	if len(s.propagators) == 0 {
		s.propagators = []propagation.TextMapPropagator{
			propagation.TraceContext{},
			propagation.Baggage{},
		}
	}

	propagator := propagation.NewCompositeTextMapPropagator(s.propagators...)
	otel.SetTextMapPropagator(propagator)
}

func (s *helper) Get() (*resource.Resource, error) {
	var options []attribute.KeyValue
	if s.serviceName != "" {
		options = append(options, semconv.ServiceNameKey.String(s.serviceName))
	}
	if len(options) > 0 {
		return resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				options...,
			),
		)
	}
	return resource.Default(), nil
}

// HelpSetPropagator setting propagator for openTelemetry
func HelpSetPropagator(loader gone.Loader) error {
	return loader.Load(&helper{})
}
