package tracer

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/tracer/gid"
	"github.com/gone-io/goner/tracer/gls"
)

// Load a tracer that uses `github.com/jtolds/gls` to implement invisible traceID propagation within a program.
func Load(loader gone.Loader) error {
	return gls.Load(loader)
}

// LoadGidTracer a tracer that uses `github.com/petermattis/goid` to implement invisible traceID propagation within a program.
func LoadGidTracer(loader gone.Loader) error {
	return gid.Load(loader)
}
