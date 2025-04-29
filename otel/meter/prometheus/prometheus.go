package prometheus

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/otel/meter"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

func Provide(_ string, i struct{}) (metric.Reader, error) {
	reader, err := prometheus.New()
	if err != nil {
		return nil, gone.ToErrorWithMsg(err, "can not create stdout trace reader")
	}
	return reader, nil
}

// Load for openTelemetry prometheus metric.Reader
func Load(loader gone.Loader) error {
	loader.MustLoad(gone.WrapFunctionProvider(Provide))
	return meter.Register(loader)
}
