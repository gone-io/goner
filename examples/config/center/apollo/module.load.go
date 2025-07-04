// Code generated by gonectl. DO NOT EDIT.
package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/apollo"
	"github.com/gone-io/goner/g"
)

// load installed gone module LoadFunc
var loaders = []gone.LoadFunc{
	apollo.Load,
}

func GoneModuleLoad(loader gone.Loader) error {
	var ops []*g.LoadOp
	for _, f := range loaders {
		ops = append(ops, g.F(f))
	}
	return g.BuildOnceLoadFunc(ops...)(loader)
}
