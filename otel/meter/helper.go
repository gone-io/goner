package meter

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	otelHelper "github.com/gone-io/goner/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
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
			exporter, err = stdoutmetric.New(
				stdoutmetric.WithPrettyPrint(),
				stdoutmetric.WithoutTimestamps(),
			)
			if err != nil {
				return gone.ToErrorWithMsg(err, "can not create stdout trace exporter")
			}
		}
		reader = metric.NewPeriodicReader(exporter)
	}

	var options = []metric.Option{metric.WithReader(reader)}

	if s.resource != nil {
		options = append(options, metric.WithResource(s.resource))
	} else {
		res, err := s.resourceGetter.Get()
		if err != nil {
			return gone.ToErrorWithMsg(err, "can not get resource")
		}
		options = append(options, metric.WithResource(res))
	}

	meterProvider := metric.NewMeterProvider(options...)
	s.afterStop(func() {
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			s.logger.Errorf("otel meter provider helper: shutdown err: %v", err)
		}
	})
	otel.SetMeterProvider(meterProvider)
	return nil
}
func (s *providerHelper) Provide(_ string) (g.IsOtelMeterLoaded, error) {
	return true, nil
}

// Register for openTelemetry openTelemetry MeterProvider
func Register(loader gone.Loader) error {
	loader.MustLoad(&providerHelper{})
	return otelHelper.HelpSetPropagator(loader)
}
