package tracer

import (
	"github.com/gone-io/gone/v2"
	"github.com/google/uuid"
	"github.com/petermattis/goid"
	"sync"
)

type tracerOverGid struct {
	gone.Flag
	logger gone.Logger `gone:"*"`
}

var gidTraceIdMap sync.Map

func (s *tracerOverGid) SetTraceId(traceId string, fn func()) {
	gid := goid.Get()
	if gid == 0 {
		s.logger.Warnf("can not get goid, tracer cannot work.")
		fn()
		return
	}
	if traceId == "" {
		traceId = uuid.New().String()
	}
	gidTraceIdMap.Store(gid, traceId)
	defer func() {
		gidTraceIdMap.Delete(gid)
	}()
	fn()
}

func (s *tracerOverGid) GetTraceId() string {
	gid := goid.Get()
	if gid == 0 {
		s.logger.Warnf("can not get goid, tracer cannot work.")
		return ""
	}
	value, ok := gidTraceIdMap.Load(gid)
	if !ok {
		return ""
	}
	return value.(string)
}

func (s *tracerOverGid) Go(fn func()) {
	traceId := s.GetTraceId()
	if traceId == "" {
		go fn()
	} else {
		go func() {
			s.SetTraceId(traceId, fn)
		}()
	}
}
