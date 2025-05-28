package http

import (
	"context"
	"crypto/tls"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	logHelper "github.com/gone-io/goner/otel/log"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/sdk/log"
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

func (c *Config) ToOtelOptions(options []otlploghttp.Option) []otlploghttp.Option {
	if c.Endpoint != "" {
		options = append(options, otlploghttp.WithEndpoint(c.Endpoint))
	}
	if c.UrlPath != "" {
		options = append(options, otlploghttp.WithURLPath(c.UrlPath))
	}
	if c.Insecure {
		options = append(options, otlploghttp.WithInsecure())
	}
	if c.TLSCfg != nil {
		options = append(options, otlploghttp.WithTLSClientConfig(c.TLSCfg))
	}
	if c.Headers != nil {
		options = append(options, otlploghttp.WithHeaders(c.Headers))
	}
	if c.Duration > 0 {
		options = append(options, otlploghttp.WithTimeout(c.Duration))
	}
	if c.Retry.Enabled {
		options = append(options, otlploghttp.WithRetry(otlploghttp.RetryConfig{
			Enabled:         c.Retry.Enabled,
			InitialInterval: c.Retry.InitialInterval,
			MaxInterval:     c.Retry.MaxInterval,
			MaxElapsedTime:  c.Retry.MaxElapsedTime,
		}))
	}
	return options
}

func Provide(_ string, i struct {
	config Config                             `gone:"config,otel.log.http"`
	proxy  otlploghttp.HTTPTransportProxyFunc `gone:"otel.http.proxy" option:"allowNil"`
}) (log.Exporter, error) {
	var options []otlploghttp.Option
	if i.proxy != nil {
		options = append(options, otlploghttp.WithProxy(i.proxy))
	}

	exporter, err := otlploghttp.New(
		context.Background(),
		i.config.ToOtelOptions(options)...,
	)

	return g.ResultError(exporter, err, "can not create oltp/http log exporter")
}

// Load for openTelemetry http log.Exporter
func Load(loader gone.Loader) error {
	loader.
		MustLoad(gone.WrapFunctionProvider(Provide)).
		MustLoadX(logHelper.Register)
	return nil
}
