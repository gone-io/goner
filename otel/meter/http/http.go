package http

import (
	"context"
	"crypto/tls"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/otel/meter"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
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

func (c *Config) ToOtelOptions(options []otlpmetrichttp.Option) []otlpmetrichttp.Option {
	if c.Endpoint != "" {
		options = append(options, otlpmetrichttp.WithEndpoint(c.Endpoint))
	}
	if c.UrlPath != "" {
		options = append(options, otlpmetrichttp.WithURLPath(c.UrlPath))
	}
	if c.Insecure {
		options = append(options, otlpmetrichttp.WithInsecure())
	}
	if c.TLSCfg != nil {
		options = append(options, otlpmetrichttp.WithTLSClientConfig(c.TLSCfg))
	}
	if c.Headers != nil {
		options = append(options, otlpmetrichttp.WithHeaders(c.Headers))
	}
	if c.Duration > 0 {
		options = append(options, otlpmetrichttp.WithTimeout(c.Duration))
	}
	if c.Retry.Enabled {
		options = append(options, otlpmetrichttp.WithRetry(otlpmetrichttp.RetryConfig{
			Enabled:         c.Retry.Enabled,
			InitialInterval: c.Retry.InitialInterval,
			MaxInterval:     c.Retry.MaxInterval,
			MaxElapsedTime:  c.Retry.MaxElapsedTime,
		}))
	}
	return options
}

func Provide(_ string, i struct {
	config Config                                `gone:"config,otel.meter.http"`
	proxy  otlpmetrichttp.HTTPTransportProxyFunc `gone:"otel.http.proxy" option:"allowNil"`
}) (metric.Exporter, error) {
	var options []otlpmetrichttp.Option
	if i.proxy != nil {
		options = append(options, otlpmetrichttp.WithProxy(i.proxy))
	}

	traceExporter, err := otlpmetrichttp.New(
		context.Background(),
		i.config.ToOtelOptions(options)...,
	)

	if err != nil {
		return nil, gone.ToErrorWithMsg(err, "can not create stdout trace exporter")
	}
	return traceExporter, nil
}

// Load for openTelemetry http metric.Exporter
func Load(loader gone.Loader) error {
	loader.MustLoad(gone.WrapFunctionProvider(Provide))
	return meter.Register(loader)
}
