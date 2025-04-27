package gone_zap

import (
	"github.com/gone-io/gone/v2"
)

func Load(loader gone.Loader) error {
	loader.
		MustLoad(&atomicLevel{}).
		MustLoad(&zapLoggerProvider{}).
		MustLoad(&sugarProvider{}).
		MustLoad(&sugar{}, gone.IsDefault(new(gone.Logger)), gone.ForceReplace())
	return nil
}
