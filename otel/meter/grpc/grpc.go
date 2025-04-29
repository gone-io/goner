package grpc

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/otel/meter"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
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

func (c *Config) ToOtelOptions(options []otlpmetricgrpc.Option) []otlpmetricgrpc.Option {
	if c.Endpoint != "" {
		options = append(options, otlpmetricgrpc.WithEndpoint(c.Endpoint))
	}
	if c.EndpointUrl != "" {
		options = append(options, otlpmetricgrpc.WithEndpointURL(c.EndpointUrl))
	}
	if c.Compressor != "" {
		options = append(options, otlpmetricgrpc.WithCompressor(c.Compressor))
	}

	if c.Headers != nil {
		options = append(options, otlpmetricgrpc.WithHeaders(c.Headers))
	}
	if c.Duration > 0 {
		options = append(options, otlpmetricgrpc.WithTimeout(c.Duration))
	}
	if c.Retry.Enabled {
		options = append(options, otlpmetricgrpc.WithRetry(otlpmetricgrpc.RetryConfig{
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
}) (metric.Exporter, error) {
	var options []otlpmetricgrpc.Option
	if i.creds != nil {
		options = append(options, otlpmetricgrpc.WithTLSCredentials(i.creds))
	} else {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		i.config.ToOtelOptions(options)...,
	)

	if err != nil {
		return nil, gone.ToErrorWithMsg(err, "can not create stdout trace exporter")
	}
	return exporter, nil
}

// Load for openTelemetry grpc metric.Exporter
func Load(loader gone.Loader) error {
	loader.MustLoad(gone.WrapFunctionProvider(Provide))
	return meter.Register(loader)
}
