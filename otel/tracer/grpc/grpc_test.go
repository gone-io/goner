package grpc

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	gone.
		NewApp(Load).
		Run(func(exporter trace.SpanExporter) {
			assert.NotNil(t, exporter)
		})

}

func TestConfig_ToOtelOptions(t *testing.T) {
	type fields struct {
		Endpoint    string
		EndpointUrl string
		Compressor  string
		Headers     map[string]string
		Duration    time.Duration
		Retry       struct {
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
	type args struct {
		options []otlptracegrpc.Option
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantLen int
	}{
		{
			name: "test",
			fields: fields{
				Endpoint:    "endpoint",
				EndpointUrl: "endpointUrl",
				Compressor:  "gzip",
				Headers:     map[string]string{"key": "value"},
				Duration:    time.Second,
				Retry: struct {
					Enabled         bool
					InitialInterval time.Duration
					MaxInterval     time.Duration
					MaxElapsedTime  time.Duration
				}{
					Enabled:         true,
					InitialInterval: time.Second,
					MaxInterval:     time.Second,
					MaxElapsedTime:  time.Second,
				},
			},
			args: args{
				options: []otlptracegrpc.Option{},
			},
			wantLen: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Endpoint:    tt.fields.Endpoint,
				EndpointUrl: tt.fields.EndpointUrl,
				Compressor:  tt.fields.Compressor,
				Headers:     tt.fields.Headers,
				Duration:    tt.fields.Duration,
				Retry:       tt.fields.Retry,
			}
			assert.Equalf(t, tt.wantLen, len(c.ToOtelOptions(tt.args.options)), "ToOtelOptions(%v)", tt.args.options)
		})
	}
}
