package main

import (
	"examples/simple/internal/interface/service"
	"github.com/gone-io/gone/v2"
)

// command app example
func main() {
	gone.
		Run(func(in struct {
			service service.IService `gone:"*"`
			appName string           `gone:"config,app.name"`
		}) {
			println(in.service.SayHello(in.appName))
		})
}
