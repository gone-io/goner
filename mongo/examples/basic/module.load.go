package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/gone-io/goner/mongo"
	"github.com/gone-io/goner/viper"
	zap "github.com/gone-io/goner/zap"
)

// load installed gone module LoadFunc
var loaders = []gone.LoadFunc{
	mongo.Load,
	viper.Load,
	zap.Load,
}

func GoneModuleLoad(loader gone.Loader) error {
	var ops []*g.LoadOp
	for _, f := range loaders {
		ops = append(ops, g.F(f))
	}
	return g.BuildOnceLoadFunc(ops...)(loader)
}

func LoadServices(loader gone.Loader) error {
	loader.
		MustLoad(&UserService{}).
		MustLoad(&DemoService{}).
		MustLoad(&AfterServerStart{})
	return nil
}
