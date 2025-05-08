package tracer

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
	"testing"
)

func TestRegister(t *testing.T) {
	gone.
		NewApp(Register).
		Loads(func(loader gone.Loader) error {
			return Register(loader)
		}).
		Run(func(is g.IsOtelTracerLoaded, tracer trace.Tracer) {
			assert.True(t, bool(is))
			assert.NotNil(t, tracer)
		})
}
