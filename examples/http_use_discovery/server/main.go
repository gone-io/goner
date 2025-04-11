package main

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner"
	"github.com/gone-io/goner/gin"
	"github.com/gone-io/goner/nacos"
	"github.com/gone-io/goner/viper"
)

func main() {
	gone.
		NewApp(goner.GinLoad, nacos.RegistryLoad, viper.Load).
		Load(&HelloController{}).
		Serve()
}

type HelloController struct {
	gone.Flag
	gin.RouteGroup `gone:"*"`
}

func (c *HelloController) Mount() gin.MountError {
	c.GET("/hello", func(in struct {
		name string `gone:"http,query"`
	}) string {
		return fmt.Sprintf("hello, %s", in.name)
	})
	return nil
}
