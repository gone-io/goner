package gid

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/google/uuid"
	"github.com/petermattis/goid"
	"sync"
)

type tracer struct {
	gone.Flag
}

var gidTraceIdMap sync.Map

func (s *tracer) SetTraceId(traceId string, fn func()) {
	if traceId == "" {
		traceId = uuid.New().String()
	}
	gid := goid.Get()
	gidTraceIdMap.Store(gid, traceId)
	defer func() {
		gidTraceIdMap.Delete(gid)
	}()
	fn()
}

func (s *tracer) GetTraceId() string {
	gid := goid.Get()
	value, ok := gidTraceIdMap.Load(gid)
	if !ok {
		return ""
	}
	return value.(string)
}

func (s *tracer) Go(fn func()) {
	traceId := s.GetTraceId()
	go func() {
		s.SetTraceId(traceId, fn)
	}()
}

// Load a tracer that uses `github.com/petermattis/goid` to implement invisible traceID propagation within a program.
func Load(loader gone.Loader) error {
	return loader.Load(&tracer{}, gone.IsDefault(new(g.Tracer)))
}
