package consul

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

func ClientLoad(loader gone.Loader) error {
	return g.SingLoadProviderFunc(ProvideConsulClient)(loader)
}

func RegistryLoad(loader gone.Loader) error {
	loader.MustLoad(&Registry{})
	return ClientLoad(loader)
}
