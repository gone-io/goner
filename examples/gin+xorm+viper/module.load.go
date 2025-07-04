// Code generated by gonectl. DO NOT EDIT.
package template_module

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/gone-io/goner/gin"
	"github.com/gone-io/goner/tracer"
	"github.com/gone-io/goner/viper"
	"github.com/gone-io/goner/xorm/mysql"
	zap "github.com/gone-io/goner/zap"
)

// load installed gone module LoadFunc
var loaders = []gone.LoadFunc{
	gin.Load,
	tracer.Load,
	viper.Load,
	mysql.Load,
	zap.Load,
}

func GoneModuleLoad(loader gone.Loader) error {
	var ops []*g.LoadOp
	for _, f := range loaders {
		ops = append(ops, g.F(f))
	}
	return g.BuildOnceLoadFunc(ops...)(loader)
}
