package main

import (
	"examples/simple/service"
	"github.com/gone-io/gone/v2"
)

func main() {
	gone.Run(func(in struct {
		service service.Service `gone:"*"`
		appName string          `gone:"config,app.name"`
	}) {
		println(in.service.SayHello(in.appName))
	})
}
