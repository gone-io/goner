package prometheus

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/metric"
	"testing"
)

func TestLoad(t *testing.T) {
	gone.
		NewApp(Load).
		Run(func(reader metric.Reader) {
			assert.NotNil(t, reader)
		})
}
