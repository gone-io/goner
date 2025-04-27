package grpc

import (
	"github.com/gone-io/gone/v2"
)

func ServerLoad(loader gone.Loader) error {
	return loader.Load(newServer())
}

func ClientRegisterLoad(loader gone.Loader) error {
	return loader.Load(NewRegister())
}

func ClientLoad(loader gone.Loader) error {
	return ClientRegisterLoad(loader)
}

func Load(loader gone.Loader) error {
	loader.
		MustLoad(newServer()).
		MustLoad(NewRegister())
	return nil
}
