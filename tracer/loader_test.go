package tracer

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad(t *testing.T) {
	gone.
		NewApp(Load).
		Run(func(tracer g.Tracer) {
			assert.NotNil(t, tracer)
		})
}

func TestLoadGidTracer(t *testing.T) {
	gone.
		NewApp(LoadGidTracer).
		Run(func(tracer g.Tracer) {
			assert.NotNil(t, tracer)
		})
}
