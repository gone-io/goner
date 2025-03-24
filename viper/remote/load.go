package remote

import "github.com/gone-io/gone/v2"

var load = gone.OnceLoad(func(loader gone.Loader) error {
	err := loader.Load(&watcher{})
	if err != nil {
		return gone.ToError(err)
	}

	return loader.Load(
		&remoteConfigure{},
		gone.Name(gone.ConfigureName),
		gone.IsDefault(new(gone.Configure)),
		gone.ForceReplace(),
	)
})

func Load(loader gone.Loader) error {
	return load(loader)
}
