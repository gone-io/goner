package log

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
	otelLog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"testing"
)

func TestRegister(t *testing.T) {
	t.Run("use stdout exporter", func(t *testing.T) {
		gone.
			NewApp(Register).
			Run(func(is g.IsOtelLogLoaded, log otelLog.Logger) {
				assert.True(t, bool(is))
				assert.NotNil(t, log)
			})
	})

	t.Run("use custom resource", func(t *testing.T) {
		gone.
			NewApp(Register).
			Loads(g.NamedThirdComponentLoadFunc("", resource.Default())).
			Run(func(is g.IsOtelLogLoaded, log otelLog.Logger) {
				assert.True(t, bool(is))
				assert.NotNil(t, log)
			})
	})
}
