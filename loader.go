package goner

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/gone-io/goner/gin"
	"github.com/gone-io/goner/tracer"
	"github.com/gone-io/goner/viper"
	zap "github.com/gone-io/goner/zap"
)

var baseLoad = g.BuildOnceLoadFunc(
	g.F(tracer.Load),
	g.F(viper.Load),
	g.F(zap.Load),
)

func BaseLoad(loader gone.Loader) error {
	return baseLoad(loader)
}

var ginLoad = g.BuildOnceLoadFunc(
	g.F(BaseLoad),
	g.F(gin.Load),
)

func GinLoad(loader gone.Loader) error {
	return ginLoad(loader)
}
