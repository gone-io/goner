package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"os"
)

func main() {
	//设置服务名称
	_ = os.Setenv("GONE_OTEL_SERVICE_NAME", "simple demo")

	gone.
		Loads(GoneModuleLoad).
		Load(&YourComponent{}).
		Run(func(c *YourComponent) {
			// 调用组件中的方法
			c.HandleRequest(context.Background())
		})
}
