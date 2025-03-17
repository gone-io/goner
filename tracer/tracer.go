package tracer

import (
	"github.com/gone-io/gone/v2"
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
