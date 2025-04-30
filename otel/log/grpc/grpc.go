package grpc

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	logHelper "github.com/gone-io/goner/otel/log"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc/credentials"
	"time"
)

type Config struct {
	Endpoint    string
	EndpointUrl string
	Compressor  string //support gzip
	Headers     map[string]string
	Duration    time.Duration

	Retry struct {
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

func (c *Config) ToOtelOptions(options []otlploggrpc.Option) []otlploggrpc.Option {
	if c.Endpoint != "" {
		options = append(options, otlploggrpc.WithEndpoint(c.Endpoint))
	}
	if c.EndpointUrl != "" {
		options = append(options, otlploggrpc.WithEndpointURL(c.EndpointUrl))
	}
	if c.Compressor != "" {
		options = append(options, otlploggrpc.WithCompressor(c.Compressor))
	}

	if c.Headers != nil {
		options = append(options, otlploggrpc.WithHeaders(c.Headers))
	}
	if c.Duration > 0 {
		options = append(options, otlploggrpc.WithTimeout(c.Duration))
	}
	if c.Retry.Enabled {
		options = append(options, otlploggrpc.WithRetry(otlploggrpc.RetryConfig{
			Enabled:         c.Retry.Enabled,
			InitialInterval: c.Retry.InitialInterval,
			MaxInterval:     c.Retry.MaxInterval,
			MaxElapsedTime:  c.Retry.MaxElapsedTime,
		}))
	}
	return options
}

func Provide(_ string, i struct {
	config Config                           `gone:"config,otel.meter.grpc"`
	creds  credentials.TransportCredentials `gone:"otel.grpc.creds" option:"allowNil"`
}) (log.Exporter, error) {
	var options []otlploggrpc.Option
	if i.creds != nil {
		options = append(options, otlploggrpc.WithTLSCredentials(i.creds))
	} else {
		options = append(options, otlploggrpc.WithInsecure())
	}

	exporter, err := otlploggrpc.New(
		context.Background(),
		i.config.ToOtelOptions(options)...,
	)

	return g.ResultError(exporter, err, "can not create oltp/grpc log exporter")
}

// Load for openTelemetry grpc log.Exporter
func Load(loader gone.Loader) error {
	loader.MustLoad(gone.WrapFunctionProvider(Provide))
	return logHelper.Register(loader)
}
