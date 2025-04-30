package grpc

import (
	"github.com/gone-io/gone/v2"
)

// ServerLoad load server
func ServerLoad(loader gone.Loader) error {
	return loader.Load(newServer())
}

// ClientRegisterLoad load client register
func ClientRegisterLoad(loader gone.Loader) error {
	return loader.Load(NewRegister())
}

// ClientLoad @deprecated use ClientRegisterLoad instead
func ClientLoad(loader gone.Loader) error {
	return ClientRegisterLoad(loader)
}

// Load all
// @deprecated use ServerLoad and ClientRegisterLoad instead
func Load(loader gone.Loader) error {
	loader.
		MustLoad(newServer()).
		MustLoad(NewRegister())
	return nil
}
