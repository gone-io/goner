package main

import (
	"github.com/gone-io/gone/v2"
	"go.uber.org/zap"
)

type UseOriginZap struct {
	gone.Flag
	zap *zap.Logger `gone:"*"`
}

func (s *UseOriginZap) PrintLog() {
	s.zap.Info("hello", zap.String("name", "gone io"))
}
