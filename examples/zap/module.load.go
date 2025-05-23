// Code generated by gonectl. DO NOT EDIT.
package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/gone-io/goner/tracer/gid"
	zap "github.com/gone-io/goner/zap"
)

// load installed gone module LoadFunc
var loaders = []gone.LoadFunc{
	gid.Load,
	zap.Load,
}

func GoneModuleLoad(loader gone.Loader) error {
	var ops []*g.LoadOp
	for _, f := range loaders {
		ops = append(ops, g.F(f))
	}
	return g.BuildOnceLoadFunc(ops...)(loader)
}
