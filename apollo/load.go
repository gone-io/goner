package apollo

import "github.com/gone-io/gone/v2"

var load = gone.OnceLoad(func(loader gone.Loader) error {
	err := loader.
		Load(
			&apolloClient{},
			gone.Name(gone.ConfigureName),
			gone.IsDefault(new(gone.Configure)),
			gone.ForceReplace(),
		)
	if err != nil {
		return err
	}
	return loader.Load(&changeListener{})
})

func Load(loader gone.Loader) error {
	return load(loader)
}
