package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner"
	"github.com/gone-io/goner/gin"
)

type HelloController struct {
	gone.Flag
	gin.IRouter `gone:"*"` // 注入路由器
}

// Mount 实现 gin.Controller 接口
func (h *HelloController) Mount() gin.MountError {
	h.GET("/hello", h.hello) // 注册路由
	return nil
}

func (h *HelloController) hello() (string, error) {
	return "Hello, Gone!", nil
}

func main() {
	gone.
		Load(&HelloController{}).
		Loads(goner.GinLoad).
		Serve()
}
