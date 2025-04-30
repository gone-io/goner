package gin

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/gone-io/goner/gin"
	"github.com/gone-io/goner/otel/meter/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type controller struct {
	gone.Flag

	router      gin.IRouter `gone:"*"` // 注入路由器
	metricsPath string      `gone:"config,otel.meter.prometheus.path=/metrics"`
}

func (c *controller) Mount() (err g.MountError) {
	var handler = promhttp.Handler()
	c.router.Any(c.metricsPath, func(ctx gin.Context) {
		handler.ServeHTTP(ctx.Writer, ctx.Request)
	})
	return nil
}

// Load prometheus http handler for metrics point
func Load(loader gone.Loader) error {
	loader.MustLoad(&controller{})
	g.PanicIfErr(gin.Load(loader))
	return prometheus.Load(loader)
}
