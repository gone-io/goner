package otel

import (
	"github.com/gone-io/gone/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type helper struct {
	gone.Flag
	propagators []propagation.TextMapPropagator `gone:"*"`
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

// HelpSetPropagator setting propagator for openTelemetry
func HelpSetPropagator(loader gone.Loader) error {
	return loader.Load(&helper{})
}
