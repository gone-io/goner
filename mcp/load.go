package goneMcp

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

// var load = g.BuildOnceLoadFunc(g.L(&serverProvider{}))
var serverLoad = g.BuildOnceLoadFunc(g.L(&serverProvider{}))

func ServerLoad(loader gone.Loader) error {
	return serverLoad(loader)
}

var clientLoad = g.BuildOnceLoadFunc(g.L(clientProvider))

func ClientLoad(loader gone.Loader) error {
	return clientLoad(loader)
}
