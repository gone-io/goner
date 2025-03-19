package gone_zap

import (
	"github.com/gone-io/gone/v2"
)

// Load load zap logger
var load = gone.OnceLoad(func(loader gone.Loader) error {
	err := loader.Load(&atomicLevel{})
	if err != nil {
		return err
	}

	err = loader.Load(&zapLoggerProvider{})
	if err != nil {
		return err
	}
	err = loader.Load(&sugarProvider{})
	if err != nil {
		return err
	}
	return loader.Load(&sugar{}, gone.IsDefault(new(gone.Logger)), gone.ForceReplace())
})

func Load(loader gone.Loader) error {
	return load(loader)
}

// Priest Deprecated, use Load instead
func Priest(loader gone.Loader) error {
	return Load(loader)
}
