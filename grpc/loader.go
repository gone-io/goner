package grpc

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

var serverLoad = g.BuildOnceLoadFunc(
	g.L(newServer()),
)

func ServerLoad(loader gone.Loader) error {
	return serverLoad(loader)
}

var clientLoad = g.BuildOnceLoadFunc(
	g.L(NewRegister()),
)

func ClientRegisterLoad(loader gone.Loader) error {
	return clientLoad(loader)
}

func ClientLoad(loader gone.Loader) error {
	return ClientRegisterLoad(loader)
}

var load = g.BuildOnceLoadFunc(
	g.L(newServer()),
	g.L(NewRegister()),
)

func Load(loader gone.Loader) error {
	return load(loader)
}
