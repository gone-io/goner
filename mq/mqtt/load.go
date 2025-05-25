package mqtt

import "github.com/gone-io/gone/v2"

// Load mqtt client
func Load(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(ProvideClient))
}
