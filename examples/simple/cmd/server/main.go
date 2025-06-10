package main

import "github.com/gone-io/gone/v2"

// server app example
func main() {
	gone.Serve()
}

// AfterServerStart  example for using hook to do something after server start
type AfterServerStart struct {
	gone.Flag
	afterStart gone.AfterStart `gone:"*"`
	logger     gone.Logger     `gone:"*"`
}

func (s *AfterServerStart) Init() {
	s.afterStart(func() {
		s.logger.Infof("after server start")
		s.logger.Infof("press `ctr + c` to stop!")
	})
}
