package apollo

import (
	"github.com/gone-io/gone/v2"
)

func Load(loader gone.Loader) error {
	loader.
		MustLoad(&apolloConfigure{},
			gone.Name(gone.ConfigureName),
			gone.IsDefault(new(gone.Configure)),
			gone.ForceReplace(),
		).
		MustLoad(&changeListener{})
	return nil
}
