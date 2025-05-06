package gls

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/google/uuid"
	"github.com/jtolds/gls"
)

type tracer struct {
	gone.Flag
}

// var traceMap sync.Map
var (
	mgr        = gls.NewContextManager()
	traceIdKey = gls.GenSym()
)

func (t *tracer) GetTraceId() string {
	v, ok := mgr.GetValue(traceIdKey)
	if ok {
		return v.(string)
	}
	return ""
}

func (t *tracer) SetTraceId(traceId string, cb func()) {
	if traceId == "" {
		traceId = uuid.New().String()
	}
	mgr.SetValues(gls.Values{traceIdKey: traceId}, cb)
}

func (t *tracer) Go(cb func()) {
	gls.Go(cb)
}

// Load a tracer that uses `github.com/jtolds/gls` to implement invisible traceID propagation within a program.
func Load(loader gone.Loader) error {
	return loader.Load(&tracer{}, gone.IsDefault(new(g.Tracer)))
}
