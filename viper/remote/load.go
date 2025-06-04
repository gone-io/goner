package remote

import "github.com/gone-io/gone/v2"

func Load(loader gone.Loader) error {
	return loader.MustLoad(&watcher{}).Load(
		&remoteConfigure{},
		gone.Name(gone.ConfigureName),
		gone.IsDefault(new(gone.Configure)),
		gone.ForceReplace(),
	)
}
