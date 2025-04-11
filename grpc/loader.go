package grpc

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

func ServerLoad(loader gone.Loader) error {
	return g.BuildLoadFunc(loader,
		g.L(newServer()),
	)
}

func ClientRegisterLoad(loader gone.Loader) error {
	return g.BuildLoadFunc(loader,
		g.L(NewRegister()),
	)
}

func ClientLoad(loader gone.Loader) error {
	return ClientRegisterLoad(loader)
}

func Load(loader gone.Loader) error {
	return g.BuildLoadFunc(loader,
		g.L(newServer()),
		g.L(NewRegister()),
	)
}
