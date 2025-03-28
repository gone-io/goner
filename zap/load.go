package gone_zap

import (
	"github.com/gone-io/gone/v2"
)

// Load load zap logger
var load = gone.OnceLoad(func(loader gone.Loader) error {
	if err := loader.Load(&atomicLevel{}); err != nil {
		return err
	}
	if err := loader.Load(&zapLoggerProvider{}); err != nil {
		return err
	}
	if err := loader.Load(&sugarProvider{}); err != nil {
		return err
	}
	return loader.Load(&sugar{}, gone.IsDefault(new(gone.Logger)), gone.ForceReplace())
})

func Load(loader gone.Loader) error {
	return load(loader)
}
