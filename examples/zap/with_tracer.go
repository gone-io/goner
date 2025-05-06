package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"go.uber.org/zap"
)

type UseTracer struct {
	gone.Flag
	logger gone.Logger `gone:"*"`
	zap    *zap.Logger `gone:"*"`
	tracer g.Tracer    `gone:"*"`
}

func (s *UseTracer) PrintLog() {
	s.tracer.SetTraceId("", func() {
		s.logger.Infof("hello with traceId")
		s.zap.Info("hello with traceId")
	})
}
