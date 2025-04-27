package etcd

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

func ClientLoad(loader gone.Loader) error {
	return g.SingLoadProviderFunc(ProvideEtecd3Client)(loader)
}

func RegistryLoad(loader gone.Loader) error {
	loader.MustLoad(&Registry{})
	return ClientLoad(loader)
}
