package tracer

import (
	"github.com/gone-io/gone/v2"
	"github.com/google/uuid"
	"github.com/petermattis/goid"
	"sync"
)

type tracerOverGid struct {
	gone.Flag
}

var gidTraceIdMap sync.Map

func (s *tracerOverGid) SetTraceId(traceId string, fn func()) {
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

func (s *tracerOverGid) GetTraceId() string {
	gid := goid.Get()
	value, ok := gidTraceIdMap.Load(gid)
	if !ok {
		return ""
	}
	return value.(string)
}

func (s *tracerOverGid) Go(fn func()) {
	traceId := s.GetTraceId()
	go func() {
		s.SetTraceId(traceId, fn)
	}()
}
