package meter

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	otelHelper "github.com/gone-io/goner/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	otelMetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

type providerHelper struct {
	gone.Flag
	resource       *resource.Resource        `gone:"*" option:"allowNil"`
	exporter       metric.Exporter           `gone:"*" option:"allowNil"`
	reader         metric.Reader             `gone:"*" option:"allowNil"`
	afterStop      gone.AfterStop            `gone:"*"`
	logger         gone.Logger               `gone:"*"`
	resourceGetter otelHelper.ResourceGetter `gone:"*"`
}

func (s *providerHelper) Init() (err error) {
	var reader metric.Reader
	if s.reader != nil {
		reader = s.reader
	} else {
		exporter := s.exporter
		if exporter == nil {
			exporter, _ = stdoutmetric.New(
				stdoutmetric.WithPrettyPrint(),
				stdoutmetric.WithoutTimestamps(),
			)
		}
		reader = metric.NewPeriodicReader(exporter)
	}

	var options = []metric.Option{metric.WithReader(reader)}

	if s.resource == nil {
		if s.resource, err = s.resourceGetter.Get(); err != nil {
			return gone.ToErrorWithMsg(err, "can not get resource")
		}
	}
	options = append(options, metric.WithResource(s.resource))

	meterProvider := metric.NewMeterProvider(options...)
	s.afterStop(func() {
		ctx := context.Background()
		g.ErrorPrinter(s.logger, meterProvider.ForceFlush(ctx), "metric provider ForceFlush")
		g.ErrorPrinter(s.logger, meterProvider.Shutdown(ctx), "metric provider Shutdown")
	})
	otel.SetMeterProvider(meterProvider)
	return nil
}

func (s *providerHelper) Provide(tagConf string) (otelMetric.Meter, error) {
	name, _ := gone.ParseGoneTag(tagConf)
	return otel.Meter(name), nil
}

var h = &providerHelper{}

// Register for openTelemetry openTelemetry MeterProvider
func Register(loader gone.Loader) error {
	if g.IsLoaded(loader, h) {
		return nil
	}

	loader.MustLoad(gone.WrapFunctionProvider(func(_ string, _ struct{}) (g.IsOtelMeterLoaded, error) {
		return true, nil
	}))
	loader.MustLoad(h)
	return otelHelper.HelpSetPropagator(loader)
}
