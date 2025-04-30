package grpc

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/gone-io/goner/otel/tracer"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
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

func (c *Config) ToOtelOptions(options []otlptracegrpc.Option) []otlptracegrpc.Option {
	if c.Endpoint != "" {
		options = append(options, otlptracegrpc.WithEndpoint(c.Endpoint))
	}
	if c.EndpointUrl != "" {
		options = append(options, otlptracegrpc.WithEndpointURL(c.EndpointUrl))
	}
	if c.Compressor != "" {
		options = append(options, otlptracegrpc.WithCompressor(c.Compressor))
	}

	if c.Headers != nil {
		options = append(options, otlptracegrpc.WithHeaders(c.Headers))
	}
	if c.Duration > 0 {
		options = append(options, otlptracegrpc.WithTimeout(c.Duration))
	}
	if c.Retry.Enabled {
		options = append(options, otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         c.Retry.Enabled,
			InitialInterval: c.Retry.InitialInterval,
			MaxInterval:     c.Retry.MaxInterval,
			MaxElapsedTime:  c.Retry.MaxElapsedTime,
		}))
	}
	return options
}

func Provide(_ string, i struct {
	logger gone.Logger                      `gone:"*"`
	config Config                           `gone:"config,otel.tracer.grpc"`
	creds  credentials.TransportCredentials `gone:"otel.grpc.creds" option:"allowNil"`
}) (trace.SpanExporter, error) {
	var options []otlptracegrpc.Option
	if i.creds != nil {
		options = append(options, otlptracegrpc.WithTLSCredentials(i.creds))
	} else {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(
		context.Background(),
		i.config.ToOtelOptions(options)...,
	)
	return g.ResultError(exporter, err, "can not create oltp/grpc trace exporter")
}

// Load for openTelemetry grpc trace.SpanExporter
func Load(loader gone.Loader) error {
	loader.MustLoad(gone.WrapFunctionProvider(Provide))
	return tracer.Register(loader)
}
