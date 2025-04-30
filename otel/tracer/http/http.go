package http

import (
	"context"
	"crypto/tls"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/gone-io/goner/otel/tracer"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

type Config struct {
	Endpoint string
	UrlPath  string
	Insecure bool
	TLSCfg   *tls.Config
	Headers  map[string]string
	Duration time.Duration
	Retry    struct {
		// Enabled indicates whether to not retry sending batches in case of
		// export failure.
		Enabled bool
		// InitialInterval the time to wait after the first failure before
		// retrying.
		InitialInterval time.Duration
		// MaxInterval is the upper bound on backoff interval. Once this value is
		// reached the delay between consecutive retries will always be
		// `MaxInterval`.
		MaxInterval time.Duration
		// MaxElapsedTime is the maximum amount of time (including retries) spent
		// trying to send a request/batch.  Once this value is reached, the data
		// is discarded.
		MaxElapsedTime time.Duration
	}
}

func (c *Config) ToOtelOptions(options []otlptracehttp.Option) []otlptracehttp.Option {
	if c.Endpoint != "" {
		options = append(options, otlptracehttp.WithEndpoint(c.Endpoint))
	}
	if c.UrlPath != "" {
		options = append(options, otlptracehttp.WithURLPath(c.UrlPath))
	}
	if c.Insecure {
		options = append(options, otlptracehttp.WithInsecure())
	}
	if c.TLSCfg != nil {
		options = append(options, otlptracehttp.WithTLSClientConfig(c.TLSCfg))
	}
	if c.Headers != nil {
		options = append(options, otlptracehttp.WithHeaders(c.Headers))
	}
	if c.Duration > 0 {
		options = append(options, otlptracehttp.WithTimeout(c.Duration))
	}
	if c.Retry.Enabled {
		options = append(options, otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         c.Retry.Enabled,
			InitialInterval: c.Retry.InitialInterval,
			MaxInterval:     c.Retry.MaxInterval,
			MaxElapsedTime:  c.Retry.MaxElapsedTime,
		}))
	}
	return options
}

func Provide(_ string, i struct {
	logger gone.Logger                          `gone:"*"`
	config Config                               `gone:"config,otel.tracer.http"`
	proxy  otlptracehttp.HTTPTransportProxyFunc `gone:"*" option:"allowNil"`
}) (trace.SpanExporter, error) {
	var options []otlptracehttp.Option
	if i.proxy != nil {
		options = append(options, otlptracehttp.WithProxy(i.proxy))
	}

	exporter, err := otlptracehttp.New(
		context.Background(),
		i.config.ToOtelOptions(options)...,
	)
	return g.ResultError(exporter, err, "can not create oltp/http trace exporter")
}

// Load for openTelemetry http trace.SpanExporter
func Load(loader gone.Loader) error {
	loader.MustLoad(gone.WrapFunctionProvider(Provide))
	return tracer.Register(loader)
}
