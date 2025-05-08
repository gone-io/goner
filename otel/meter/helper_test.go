package meter

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	otelMetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"testing"
)

type mockReader struct {
	gone.Flag
	metric.Reader
}

func TestRegister(t *testing.T) {
	t.Run("use stdout exporter", func(t *testing.T) {
		gone.
			NewApp(Register).
			Loads(func(loader gone.Loader) error {
				return Register(loader)
			}).
			Run(func(is g.IsOtelMeterLoaded, meter otelMetric.Meter) {
				assert.True(t, bool(is))
				assert.NotNil(t, meter)
			})
	})

	t.Run("use custom reader", func(t *testing.T) {
		exporter, _ := stdoutmetric.New(
			stdoutmetric.WithPrettyPrint(),
			stdoutmetric.WithoutTimestamps(),
		)

		gone.
			NewApp(Register).
			Load(&mockReader{
				Reader: metric.NewPeriodicReader(exporter),
			}).
			Run(func(is g.IsOtelMeterLoaded, meter otelMetric.Meter) {
				assert.True(t, bool(is))
				assert.NotNil(t, meter)
			})
	})

	t.Run("use custom resource", func(t *testing.T) {
		gone.
			NewApp(Register).
			Loads(g.NamedThirdComponentLoadFunc("", resource.Default())).
			Run(func(is g.IsOtelMeterLoaded, meter otelMetric.Meter) {
				assert.True(t, bool(is))
				assert.NotNil(t, meter)
			})
	})
}
