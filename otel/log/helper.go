package log

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	otelHelper "github.com/gone-io/goner/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	otelLog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

type helper struct {
	gone.Flag

	resource       *resource.Resource        `gone:"*" option:"allowNil"`
	exporter       log.Exporter              `gone:"*" option:"allowNil"`
	logger         gone.Logger               `gone:"*" option:"lazy"`
	afterStop      gone.AfterStop            `gone:"*"`
	resourceGetter otelHelper.ResourceGetter `gone:"*"`
}

func (s *helper) Init() (err error) {
	exporter := s.exporter
	if exporter == nil {
		exporter, _ = stdoutlog.New(
			stdoutlog.WithPrettyPrint(),
			stdoutlog.WithoutTimestamps(),
		)
	}

	var options = []log.LoggerProviderOption{log.WithProcessor(log.NewBatchProcessor(exporter))}

	if s.resource != nil {
		options = append(options, log.WithResource(s.resource))
	} else {
		res, err := s.resourceGetter.Get()
		if err != nil {
			return gone.ToErrorWithMsg(err, "can not get resource")
		}
		options = append(options, log.WithResource(res))
	}

	provider := log.NewLoggerProvider(options...)

	s.afterStop(func() {
		ctx := context.Background()
		err := provider.ForceFlush(ctx)
		g.ErrorPrinter(s.logger, err, "provider.ForceFlush")
		g.ErrorPrinter(s.logger, provider.Shutdown(ctx), "otel logger provider helper shutdown")
	})

	global.SetLoggerProvider(provider)
	return nil
}

func (s *helper) Provide(_ string) (g.IsOtelLogLoaded, error) {
	return true, nil
}

// Register for openTelemetry LoggerProvider
func Register(loader gone.Loader) error {
	loader.MustLoad(gone.WrapFunctionProvider(func(tagConf string, param struct{}) (otelLog.Logger, error) {
		name, _ := gone.ParseGoneTag(tagConf)
		return global.Logger(name), nil
	}))

	loader.MustLoad(&helper{})
	return otelHelper.HelpSetPropagator(loader)
}
