package http

import (
	"crypto/tls"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	gone.
		NewApp(
			Load,
			g.NamedThirdComponentLoadFunc[otlpmetrichttp.HTTPTransportProxyFunc]("", func(request *http.Request) (*url.URL, error) {
				return nil, nil
			}),
		).
		Run(func(exporter metric.Exporter) {
			assert.NotNil(t, exporter)
		})
}

func TestConfig_ToOtelOptions(t *testing.T) {
	type fields struct {
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
	type args struct {
		options []otlpmetrichttp.Option
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
				Endpoint: "http://localhost:4318",
				UrlPath:  "/v1/traces",
				Insecure: true,
				TLSCfg:   &tls.Config{},
				Headers:  map[string]string{"Content-Type": "application/json"},
				Duration: 10 * time.Second,
				Retry: struct {
					Enabled         bool
					InitialInterval time.Duration
					MaxInterval     time.Duration
					MaxElapsedTime  time.Duration
				}{
					Enabled:         true,
					InitialInterval: 0,
				},
			},
			args: args{
				options: []otlpmetrichttp.Option{},
			},
			wantLen: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Endpoint: tt.fields.Endpoint,
				UrlPath:  tt.fields.UrlPath,
				Insecure: tt.fields.Insecure,
				TLSCfg:   tt.fields.TLSCfg,
				Headers:  tt.fields.Headers,
				Duration: tt.fields.Duration,
				Retry:    tt.fields.Retry,
			}
			assert.Equalf(t, tt.wantLen, len(c.ToOtelOptions(tt.args.options)), "ToOtelOptions(%v)", tt.args.options)
		})
	}
}
