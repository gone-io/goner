package prometheus

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/gone-io/goner/otel/meter"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

func Provide(_ string, _ struct{}) (metric.Reader, error) {
	reader, err := prometheus.New()
	return g.ResultError(reader, err, "can not create prometheus metric reader")
}

// Load for openTelemetry prometheus metric.Reader
func Load(loader gone.Loader) error {
	loader.
		MustLoad(gone.WrapFunctionProvider(Provide)).
		MustLoadX(meter.Register)
	return nil
}
