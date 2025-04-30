package otel

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestHelpSetPropagator(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_ = os.Setenv("GONE_OTEL_SERVICE_NAME", "test")
		defer func() {
			_ = os.Unsetenv("GONE_OTEL_SERVICE_NAME")
		}()

		gone.
			NewApp(HelpSetPropagator).
			Run(func(g ResourceGetter) {
				r, err := g.Get()
				assert.Nil(t, err)
				s := r.String()
				assert.Contains(t, s, "service.name")
				assert.Contains(t, s, "service.name=test")
			})
	})
	t.Run("get default resource", func(t *testing.T) {
		gone.
			NewApp(HelpSetPropagator).
			Run(func(g ResourceGetter) {
				r, err := g.Get()
				assert.Nil(t, err)
				s := r.String()
				assert.Contains(t, s, "service.name")
				assert.NotContains(t, s, "service.name=test")
			})
	})
}
