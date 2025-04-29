package log

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	otelHelper "github.com/gone-io/goner/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

type helper struct {
	gone.Flag

	resource  *resource.Resource `gone:"*" option:"allowNil"`
	exporter  log.Exporter       `gone:"*" option:"allowNil"`
	logger    gone.Logger        `gone:"*"`
	afterStop gone.AfterStop     `gone:"*"`
}

func (s *helper) Init() (err error) {
	exporter := s.exporter
	if exporter == nil {
		exporter, err = stdoutlog.New(
			stdoutlog.WithPrettyPrint(),
			stdoutlog.WithoutTimestamps(),
		)
		if err != nil {
			return gone.ToErrorWithMsg(err, "can not create stdout log exporter")
		}
	}

	var options = []log.LoggerProviderOption{log.WithProcessor(log.NewBatchProcessor(exporter))}

	if s.resource != nil {
		options = append(options, log.WithResource(s.resource))
	}

	provider := log.NewLoggerProvider(options...)

	s.afterStop(func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			s.logger.Errorf("otel logger provider helper: shutdown err: %v", err)
		}
	})

	global.SetLoggerProvider(provider)
	return nil
}

func (s *helper) Provide(_ string) (g.IsOtelLogLoaded, error) {
	return true, nil
}

// Register for openTelemetry LoggerProvider
func Register(loader gone.Loader) error {
	loader.MustLoad(&helper{})
	return otelHelper.HelpSetPropagator(loader)
}
