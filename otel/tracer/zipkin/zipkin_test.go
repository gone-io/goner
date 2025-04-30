package zipkin

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"
	"testing"
)

func TestLoad(t *testing.T) {
	gone.NewApp(Load).Run(func(tracer trace.SpanExporter) {
		assert.NotNil(t, tracer)
	})
}
