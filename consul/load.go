package consul

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

func ClientLoad(loader gone.Loader) error {
	return g.SingLoadProviderFunc(ProvideConsulClient)(loader)
}

var load = g.BuildOnceLoadFunc(g.F(ClientLoad), g.L(&Registry{}))

func RegistryLoad(loader gone.Loader) error {
	return load(loader)
}
