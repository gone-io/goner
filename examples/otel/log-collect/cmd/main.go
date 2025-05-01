package main

import (
	"github.com/gone-io/gone/v2"
)

func main() {
	gone.Run(func(logger gone.Logger, i struct {
		name string `gone:"config,otel.service.name"`
	}) {
		logger.Infof("service name: %s", i.name)
		logger.Infof("hello world")
		logger.Debugf("debug info")
		logger.Warnf("warn info")
		logger.Errorf("error info")
	})
}
