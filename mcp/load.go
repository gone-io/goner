package goneMcp

import (
	"github.com/gone-io/gone/v2"
)

func ServerLoad(loader gone.Loader) error {
	return loader.Load(&serverProvider{})
}

func ClientLoad(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(clientProvide))
}
